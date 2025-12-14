package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Lazywords2006/web/server/database"
	"github.com/Lazywords2006/web/server/models"
	"github.com/Lazywords2006/web/server/utils"
)

// ActivateRequest 激活请求
type ActivateRequest struct {
	Key  string `json:"key"`
	HWID string `json:"hwid"`
}

// ActivateResponse 激活响应
type ActivateResponse struct {
	Status string `json:"status"`
	Token  string `json:"token,omitempty"`
	Error  string `json:"error,omitempty"`
}

// HandleActivate 处理许可证激活
func HandleActivate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求
	var req ActivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logActivation(req.Key, req.HWID, "activate", r, false, "Invalid request format")
		respondError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	log.Printf("[Activate] Request: key=%s, hwid=%s...", req.Key, truncate(req.HWID, 16))

	// 查询许可证
	var license models.License
	var validityDays int
	var expiresAt sql.NullTime
	var hwid sql.NullString
	var userID sql.NullInt64

	err := database.DB.QueryRow(`
		SELECT id, license_key, product_name, hwid, status, max_devices, validity_days, expires_at, user_id
		FROM licenses WHERE license_key = ?
	`, req.Key).Scan(
		&license.ID, &license.LicenseKey, &license.ProductName,
		&hwid, &license.Status, &license.MaxDevices,
		&validityDays, &expiresAt, &userID,
	)

	// 处理 NULL hwid
	if hwid.Valid {
		license.HWID = hwid.String
	}
	// 处理 NULL user_id
	if userID.Valid {
		license.UserID = userID.Int64
	}

	if err == sql.ErrNoRows {
		log.Printf("[Activate] REJECTED: License not found")
		logActivation(req.Key, req.HWID, "activate", r, false, "License not found")
		respondError(w, "Invalid license key", http.StatusForbidden)
		return
	}

	if err != nil {
		log.Printf("[Activate] ERROR: Database error: %v", err)
		logActivation(req.Key, req.HWID, "activate", r, false, "Database error")
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 如果已有过期时间,使用它;否则在激活时计算
	if expiresAt.Valid {
		license.ExpiresAt = expiresAt.Time
	}

	// 验证许可证状态
	if license.Status == "banned" {
		log.Printf("[Activate] REJECTED: License banned")
		logActivation(req.Key, req.HWID, "activate", r, false, "License banned")
		respondError(w, "License has been banned", http.StatusForbidden)
		return
	}

	if license.Status == "expired" {
		log.Printf("[Activate] REJECTED: License expired")
		logActivation(req.Key, req.HWID, "activate", r, false, "License expired")
		respondError(w, "License has expired", http.StatusForbidden)
		return
	}

	// 检查过期时间 (只有已激活且设置了expires_at的才检查)
	if expiresAt.Valid && time.Now().After(license.ExpiresAt) {
		// 更新状态为expired
		database.DB.Exec("UPDATE licenses SET status = 'expired' WHERE id = ?", license.ID)
		log.Printf("[Activate] REJECTED: License expired (expires_at check)")
		logActivation(req.Key, req.HWID, "activate", r, false, "License expired")
		respondError(w, "License has expired", http.StatusForbidden)
		return
	}

	// 检查硬件绑定
	if license.Status == "active" {
		if license.HWID != req.HWID {
			log.Printf("[Activate] REJECTED: HWID mismatch (expected %s, got %s)", truncate(license.HWID, 16), truncate(req.HWID, 16))
			logActivation(req.Key, req.HWID, "activate", r, false, "HWID mismatch")
			respondError(w, "License already activated on another device", http.StatusForbidden)
			return
		}
	}

	// 激活许可证 (首次激活)
	if license.Status == "unused" {
		// 计算过期时间: 当前时间 + validity_days 天
		now := time.Now()
		expiresAt := now.AddDate(0, 0, validityDays)
		license.ExpiresAt = expiresAt

		_, err = database.DB.Exec(`
			UPDATE licenses
			SET hwid = ?, status = 'active', activated_at = ?, expires_at = ?, updated_at = ?
			WHERE id = ?
		`, req.HWID, now, expiresAt, now, license.ID)

		if err != nil {
			log.Printf("[Activate] ERROR: Failed to update license: %v", err)
			logActivation(req.Key, req.HWID, "activate", r, false, "Failed to update license")
			respondError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("[Activate] License activated successfully, expires_at=%s (%d days)", expiresAt.Format("2006-01-02"), validityDays)
	} else {
		// 已激活，验证HWID
		log.Printf("[Activate] License already active, validating HWID")
	}

	// 生成JWT令牌
	token, err := utils.GenerateJWT(license.LicenseKey, req.HWID, license.ExpiresAt)
	if err != nil {
		log.Printf("[Activate] ERROR: Failed to generate token: %v", err)
		logActivation(req.Key, req.HWID, "activate", r, false, "Failed to generate token")
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[Activate] SUCCESS: Token issued")
	logActivation(req.Key, req.HWID, "activate", r, true, "")

	// 返回成功响应
	respondJSON(w, ActivateResponse{
		Status: "success",
		Token:  token,
	}, http.StatusOK)
}

// HeartbeatResponse 心跳响应
type HeartbeatResponse struct {
	Status string `json:"status"`
}

// HandleHeartbeat 处理心跳请求
func HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 提取Authorization头
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		log.Printf("[Heartbeat] REJECTED: No valid Authorization header")
		respondJSON(w, HeartbeatResponse{Status: "dead"}, http.StatusUnauthorized)
		return
	}

	token := authHeader[7:]

	// 验证JWT
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		log.Printf("[Heartbeat] REJECTED: Invalid token: %v", err)
		respondJSON(w, HeartbeatResponse{Status: "dead"}, http.StatusUnauthorized)
		return
	}

	licenseKey, ok := (*claims)["license_key"].(string)
	if !ok {
		log.Printf("[Heartbeat] REJECTED: Invalid token claims")
		respondJSON(w, HeartbeatResponse{Status: "dead"}, http.StatusUnauthorized)
		return
	}

	hwid, _ := (*claims)["hwid"].(string)

	// 查询许可证状态
	var status string
	var expiresAt time.Time
	err = database.DB.QueryRow(`
		SELECT status, expires_at FROM licenses WHERE license_key = ?
	`, licenseKey).Scan(&status, &expiresAt)

	if err != nil {
		log.Printf("[Heartbeat] REJECTED: License not found or error: %v", err)
		logActivation(licenseKey, hwid, "heartbeat", r, false, "License not found")
		respondJSON(w, HeartbeatResponse{Status: "dead"}, http.StatusForbidden)
		return
	}

	// 检查状态
	if status != "active" {
		log.Printf("[Heartbeat] REJECTED: License status is %s", status)
		logActivation(licenseKey, hwid, "heartbeat", r, false, "License not active")
		respondJSON(w, HeartbeatResponse{Status: "dead"}, http.StatusForbidden)
		return
	}

	// 检查过期时间
	if time.Now().After(expiresAt) {
		database.DB.Exec("UPDATE licenses SET status = 'expired' WHERE license_key = ?", licenseKey)
		log.Printf("[Heartbeat] REJECTED: License expired")
		logActivation(licenseKey, hwid, "heartbeat", r, false, "License expired")
		respondJSON(w, HeartbeatResponse{Status: "dead"}, http.StatusForbidden)
		return
	}

	// 更新最后心跳时间
	database.DB.Exec(`
		UPDATE licenses SET last_heartbeat = ? WHERE license_key = ?
	`, time.Now(), licenseKey)

	log.Printf("[Heartbeat] OK: %s", truncate(licenseKey, 20))
	logActivation(licenseKey, hwid, "heartbeat", r, true, "")

	// 返回成功
	respondJSON(w, HeartbeatResponse{Status: "alive"}, http.StatusOK)
}

// 辅助函数

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, map[string]string{"error": message}, statusCode)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func logActivation(licenseKey, hwid, action string, r *http.Request, success bool, errorMsg string) {
	ip := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")

	_, err := database.DB.Exec(`
		INSERT INTO activation_logs (license_key, hwid, action, ip_address, user_agent, success, error_msg)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, licenseKey, hwid, action, ip, userAgent, success, errorMsg)

	if err != nil {
		log.Printf("[Log] Warning: Failed to log activation: %v", err)
	}
}
