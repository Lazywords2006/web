package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// MockActivateRequest 激活请求
type MockActivateRequest struct {
	Key  string `json:"key"`
	HWID string `json:"hwid"`
}

// MockActivateResponse 激活响应
type MockActivateResponse struct {
	Status string `json:"status"`
	Token  string `json:"token,omitempty"`
	Error  string `json:"error,omitempty"`
}

// MockHeartbeatResponse 心跳响应
type MockHeartbeatResponse struct {
	Status string `json:"status"`
}

// 有效的测试许可证密钥
var validLicenseKeys = map[string]bool{
	"TEST-LICENSE-KEY-12345":     true,
	"DEMO-KEY-ABCDEF":            true,
	"VALID-KEY-XYZ789":           true,
}

// 存储激活的令牌（简单示例，生产环境应使用数据库）
var activeTokens = make(map[string]bool)

func main() {
	log.Println("=== Mock License Server ===")
	log.Println("Starting on http://localhost:8080")
	log.Println("\nValid test license keys:")
	for key := range validLicenseKeys {
		log.Printf("  - %s", key)
	}
	log.Println("\nAPI Endpoints:")
	log.Println("  POST /api/activate   - License activation")
	log.Println("  POST /api/heartbeat  - Heartbeat validation")
	log.Println("========================================\n")

	http.HandleFunc("/api/activate", handleActivate)
	http.HandleFunc("/api/heartbeat", handleHeartbeat)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// handleActivate 处理许可证激活请求
func handleActivate(w http.ResponseWriter, r *http.Request) {
	// 只接受POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求
	var req MockActivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[Activate] Bad request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MockActivateResponse{
			Status: "error",
			Error:  "Invalid request format",
		})
		return
	}

	log.Printf("[Activate] Request: key=%s, hwid=%s...", req.Key, req.HWID[:16])

	// 验证许可证密钥
	if !validLicenseKeys[req.Key] {
		log.Printf("[Activate] REJECTED: Invalid key")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(MockActivateResponse{
			Status: "error",
			Error:  "Invalid license key",
		})
		return
	}

	// 生成Mock JWT令牌
	token := generateMockToken(req.Key, req.HWID)
	activeTokens[token] = true

	log.Printf("[Activate] SUCCESS: Token issued: %s...", token[:20])

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MockActivateResponse{
		Status: "success",
		Token:  token,
	})
}

// handleHeartbeat 处理心跳请求
func handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	// 只接受POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 检查Authorization头
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Printf("[Heartbeat] REJECTED: No Authorization header")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(MockHeartbeatResponse{
			Status: "dead",
		})
		return
	}

	// 提取令牌
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		log.Printf("[Heartbeat] REJECTED: Invalid Authorization format")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(MockHeartbeatResponse{
			Status: "dead",
		})
		return
	}

	// 验证令牌
	if !activeTokens[token] {
		log.Printf("[Heartbeat] REJECTED: Invalid token")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(MockHeartbeatResponse{
			Status: "dead",
		})
		return
	}

	log.Printf("[Heartbeat] OK: Token %s...", token[:20])

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MockHeartbeatResponse{
		Status: "alive",
	})
}

// generateMockToken 生成Mock JWT令牌（仅用于测试）
func generateMockToken(key, hwid string) string {
	// 在实际应用中，这里应该使用真正的JWT库
	// 例如：github.com/golang-jwt/jwt
	return "mock-jwt-token-" + key + "-" + hwid[:8]
}
