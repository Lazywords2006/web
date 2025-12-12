package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 认证客户端
type Client struct {
	ServerURL  string
	HTTPClient *http.Client
	Token      string
}

// ActivateRequest 激活请求结构
type ActivateRequest struct {
	Key  string `json:"key"`
	HWID string `json:"hwid"`
}

// ActivateResponse 激活响应结构
type ActivateResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	Error  string `json:"error,omitempty"`
}

// HeartbeatResponse 心跳响应结构
type HeartbeatResponse struct {
	Status string `json:"status"`
}

// NewClient 创建新的认证客户端
func NewClient(serverURL string) *Client {
	return &Client{
		ServerURL: serverURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			// TODO: 生产环境中应启用SSL Pinning
			// 使用自定义Transport配置TLS证书验证
			// Transport: &http.Transport{
			//     TLSClientConfig: &tls.Config{
			//         RootCAs:      certPool,          // 固定证书池
			//         MinVersion:   tls.VersionTLS12,  // 最低TLS版本
			//         CipherSuites: secureCiphers,     // 安全的加密套件
			//     },
			// }
		},
	}
}

// Activate 激活许可证
// 发送许可证密钥和硬件ID到服务器进行验证
// 成功时返回JWT令牌并存储在客户端实例中
func (c *Client) Activate(licenseKey, hwid string) error {
	// 构建请求体
	reqBody := ActivateRequest{
		Key:  licenseKey,
		HWID: hwid,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/api/activate", c.ServerURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SecureClient/1.0")

	// TODO: 添加请求签名以防止中间人攻击
	// req.Header.Set("X-Request-Signature", generateHMAC(jsonData, secret))

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send activation request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// 解析响应
	var activateResp ActivateResponse
	if err := json.Unmarshal(body, &activateResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		if activateResp.Error != "" {
			return fmt.Errorf("activation failed: %s", activateResp.Error)
		}
		return fmt.Errorf("activation failed with status code: %d", resp.StatusCode)
	}

	// 验证响应
	if activateResp.Status != "success" {
		return fmt.Errorf("activation unsuccessful: %s", activateResp.Status)
	}

	if activateResp.Token == "" {
		return fmt.Errorf("no token received from server")
	}

	// 存储令牌
	c.Token = activateResp.Token
	return nil
}

// Heartbeat 发送心跳请求
// 向服务器发送心跳以验证许可证仍然有效
// 返回错误表示许可证已失效或连接失败
func (c *Client) Heartbeat() error {
	if c.Token == "" {
		return fmt.Errorf("no token available, please activate first")
	}

	// 创建心跳请求
	url := fmt.Sprintf("%s/api/heartbeat", c.ServerURL)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create heartbeat request: %w", err)
	}

	// 设置Authorization头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("User-Agent", "SecureClient/1.0")

	// TODO: 添加请求时间戳和签名防止重放攻击
	// timestamp := time.Now().Unix()
	// req.Header.Set("X-Timestamp", strconv.FormatInt(timestamp, 10))
	// req.Header.Set("X-Signature", generateSignature(c.Token, timestamp))

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read heartbeat response: %w", err)
	}

	// 检查状态码
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("license invalidated by server (status: %d)", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat failed with status code: %d", resp.StatusCode)
	}

	// 解析响应
	var heartbeatResp HeartbeatResponse
	if err := json.Unmarshal(body, &heartbeatResp); err != nil {
		return fmt.Errorf("failed to parse heartbeat response: %w", err)
	}

	// 验证心跳状态
	if heartbeatResp.Status != "alive" {
		return fmt.Errorf("license is no longer active: %s", heartbeatResp.Status)
	}

	return nil
}

// GetToken 获取当前存储的令牌
func (c *Client) GetToken() string {
	return c.Token
}

// IsAuthenticated 检查是否已认证
func (c *Client) IsAuthenticated() bool {
	return c.Token != ""
}
