# Secure Always-Online Client (Go DRM System)

一个使用Go语言实现的健壮的"永久在线"客户端应用程序包装器，具有DRM（数字版权管理）和网络安全功能。

## 功能特性

### 核心功能
- ✅ **跨平台硬件ID生成**：支持 Windows、Linux、macOS
- ✅ **许可证激活系统**：与远程服务器进行密钥验证
- ✅ **持久心跳监控**：后台Goroutine维持连接验证
- ✅ **智能重试机制**：3次重试，2秒延迟
- ✅ **强制终止开关**：验证失败时立即终止应用
- ✅ **模块化架构**：清晰的代码组织结构

### 安全特性
- 🔒 基于JWT令牌的认证
- 🔒 硬件绑定防止许可证共享
- 🔒 预留SSL Pinning接口（生产环境推荐）
- 🔒 请求签名和时间戳验证（TODO注释标注）

## 项目结构

```
网络验证/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块定义
├── config.json.example     # 配置文件示例
├── auth/
│   └── auth.go            # 许可证激活和验证逻辑
├── hwid/
│   └── hwid.go            # 跨平台硬件ID生成
└── heartbeat/
    └── heartbeat.go       # 心跳监控和强制退出机制
```

## 快速开始

### 1. 环境要求

- Go 1.21 或更高版本
- 网络连接到许可证服务器

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置

复制配置文件示例：

```bash
cp config.json.example config.json
```

编辑 `config.json`：

```json
{
  "server_url": "http://your-license-server.com",
  "license_key": "YOUR-LICENSE-KEY-HERE",
  "heartbeat_interval_seconds": 30,
  "max_retries": 3,
  "retry_delay_seconds": 2
}
```

**配置说明：**
- `server_url`: 许可证服务器地址
- `license_key`: 许可证密钥（可留空，运行时输入）
- `heartbeat_interval_seconds`: 心跳间隔（秒）
- `max_retries`: 心跳失败最大重试次数
- `retry_delay_seconds`: 重试延迟（秒）

### 4. 编译

```bash
# 编译当前平台
go build -o secure-client

# 跨平台编译示例
# Windows
GOOS=windows GOARCH=amd64 go build -o secure-client.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o secure-client-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o secure-client-mac
```

### 5. 运行

```bash
./secure-client
```

## 架构与逻辑流程

### 启动流程

```
1. 加载配置 → 2. 生成HWID → 3. 激活许可证 → 4. 启动心跳 → 5. 运行业务逻辑
```

#### 1️⃣ 启动阶段（Startup）
- 生成稳定的硬件ID（HWID）基于 CPU/磁盘/主板
- 从配置文件或用户输入获取许可证密钥
- 发送 `POST /api/activate` 请求：`{key, hwid}`
- 接收并存储JWT令牌

#### 2️⃣ 运行阶段（Runtime - Heartbeat）
- 后台Goroutine每30秒发送 `POST /api/heartbeat`
- **重试逻辑**：失败时重试3次，间隔2秒
- **Kill Switch**：所有重试失败或服务器返回"Banned/Expired"时，调用 `ForceExit()` 立即终止进程

#### 3️⃣ 业务逻辑（Business Logic）
- 只有激活成功后才执行 `RunMainApp()` 函数
- 这是实际软件功能的占位符

### API 接口契约（Mock）

#### 激活接口
```http
POST http://localhost:8080/api/activate
Content-Type: application/json

{
  "key": "LICENSE-KEY-HERE",
  "hwid": "abc123..."
}
```

**成功响应 (200):**
```json
{
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**失败响应 (403):**
```json
{
  "error": "Invalid license key"
}
```

#### 心跳接口
```http
POST http://localhost:8080/api/heartbeat
Authorization: Bearer <token>
```

**正常响应 (200):**
```json
{
  "status": "alive"
}
```

**无效响应 (401/403):**
```json
{
  "status": "dead"
}
```

## 代码模块说明

### 1. hwid/hwid.go - 硬件ID生成

**功能：**
- 跨平台硬件指纹生成
- Windows: 使用WMIC获取CPU/主板/磁盘序列号
- Linux: 读取 `/proc/cpuinfo` 和 `/etc/machine-id`
- macOS: 使用 `ioreg` 获取硬件UUID和序列号
- 返回SHA256哈希值作为稳定标识

**关键函数：**
```go
func GetHardwareID() (string, error)
```

### 2. auth/auth.go - 认证模块

**功能：**
- 许可证激活逻辑
- JWT令牌管理
- 心跳请求发送
- SSL Pinning预留接口（TODO注释）

**关键类型：**
```go
type Client struct {
    ServerURL  string
    HTTPClient *http.Client
    Token      string
}
```

**关键函数：**
```go
func (c *Client) Activate(licenseKey, hwid string) error
func (c *Client) Heartbeat() error
```

### 3. heartbeat/heartbeat.go - 心跳监控

**功能：**
- 后台Goroutine心跳循环
- 重试逻辑（3次，2秒延迟）
- 强制退出机制
- 错误回调支持

**关键类型：**
```go
type Monitor struct {
    client        AuthClient
    interval      time.Duration
    maxRetries    int
    retryDelay    time.Duration
}
```

**关键函数：**
```go
func (m *Monitor) Start()                    // 启动监控
func ForceExit(reason string)                // 强制终止
func GracefulShutdown(reason string)         // 优雅关闭
```

### 4. main.go - 主程序

**功能：**
- 应用程序入口点
- 配置加载
- 模块编排
- 业务逻辑占位符

**关键函数：**
```go
func main()                      // 主入口
func loadConfig() (*Config, error)
func RunMainApp()                // 业务逻辑占位符（替换为实际代码）
```

## 生产环境部署建议

### 安全加固

1. **启用SSL Pinning**（见 `auth/auth.go:35` 注释）
   ```go
   Transport: &http.Transport{
       TLSClientConfig: &tls.Config{
           RootCAs:      certPool,          // 固定证书池
           MinVersion:   tls.VersionTLS12,
       },
   }
   ```

2. **添加请求签名**（见 `auth/auth.go:67` 和 `auth/auth.go:121` 注释）
   ```go
   req.Header.Set("X-Request-Signature", generateHMAC(jsonData, secret))
   ```

3. **实现防篡改机制**（见 `heartbeat/heartbeat.go:99` 注释）
   - 清理敏感数据
   - 删除临时文件
   - 写入审计日志

### 编译优化

```bash
# 去除调试信息，减小体积
go build -ldflags="-s -w" -o secure-client

# 使用upx压缩（可选）
upx --best secure-client
```

### 混淆保护（可选）

考虑使用Go混淆工具保护二进制文件：
- [garble](https://github.com/burrowers/garble)
- [gobfuscate](https://github.com/unixpickle/gobfuscate)

## 测试

### 模拟许可证服务器

创建一个简单的测试服务器（test-server.go）：

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/api/activate", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "status": "success",
            "token":  "mock-jwt-token-abc123",
        })
    })

    http.HandleFunc("/api/heartbeat", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "status": "alive",
        })
    })

    log.Println("Mock server running on :8080")
    http.ListenAndServe(":8080", nil)
}
```

运行测试服务器：
```bash
go run test-server.go
```

### 单元测试

运行所有测试：
```bash
go test ./...
```

## 常见问题

### Q: 如何更换业务逻辑？
**A:** 修改 `main.go` 中的 `RunMainApp()` 函数，替换为您的实际应用代码。

### Q: 心跳间隔太频繁怎么办？
**A:** 在 `config.json` 中调整 `heartbeat_interval_seconds` 参数。

### Q: 如何处理网络不稳定？
**A:** 增加 `max_retries` 和 `retry_delay_seconds` 参数。

### Q: 可以在无网络环境使用吗？
**A:** 不可以。这是"Always-Online"系统，必须保持网络连接。如需离线模式，需要修改架构。

### Q: 如何禁用心跳监控？
**A:** 不建议禁用。如果确实需要，注释掉 `main.go` 中的心跳启动代码（第64-72行）。

## 许可证

本项目代码仅供学习和研究使用。

## 贡献

欢迎提交Issue和Pull Request。

## 联系方式

如有问题，请创建Issue或联系项目维护者。

---

**⚠️ 重要提示：**
- 此系统设计用于合法的软件保护目的
- 请确保遵守当地法律法规
- 不要用于恶意软件或非法用途
- 生产环境部署前请进行充分的安全审计
