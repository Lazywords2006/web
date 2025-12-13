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

// HandleGenerateLicense 生成新许可证
func HandleGenerateLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: 添加管理员认证中间件

	var req struct {
		ProductName string `json:"product_name"`
		Duration    int    `json:"duration"` // 天数
		MaxDevices  int    `json:"max_devices"`
		UserID      int64  `json:"user_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 生成许可证密钥
	licenseKey, err := utils.GenerateLicenseKey()
	if err != nil {
		log.Printf("[Admin] Failed to generate license key: %v", err)
		respondError(w, "Failed to generate license", http.StatusInternalServerError)
		return
	}

	// 计算过期时间
	expiresAt := time.Now().AddDate(0, 0, req.Duration)

	// 插入数据库
	_, err = database.DB.Exec(`
		INSERT INTO licenses (license_key, product_name, status, max_devices, expires_at, user_id)
		VALUES (?, ?, 'unused', ?, ?, ?)
	`, licenseKey, req.ProductName, req.MaxDevices, expiresAt, req.UserID)

	if err != nil {
		log.Printf("[Admin] Failed to insert license: %v", err)
		respondError(w, "Failed to save license", http.StatusInternalServerError)
		return
	}

	log.Printf("[Admin] Generated license: %s for product %s", licenseKey, req.ProductName)

	respondJSON(w, map[string]interface{}{
		"license_key":  licenseKey,
		"product_name": req.ProductName,
		"expires_at":   expiresAt,
		"max_devices":  req.MaxDevices,
	}, http.StatusOK)
}

// HandleListLicenses 列出所有许可证
func HandleListLicenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: 添加管理员认证

	// 查询参数
	status := r.URL.Query().Get("status")
	userID := r.URL.Query().Get("user_id")

	query := "SELECT id, license_key, product_name, hwid, status, max_devices, expires_at, activated_at, created_at, user_id, last_heartbeat FROM licenses WHERE 1=1"
	args := []interface{}{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if userID != "" {
		query += " AND user_id = ?"
		args = append(args, userID)
	}

	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("[Admin] Failed to query licenses: %v", err)
		respondError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	licenses := []models.License{}
	for rows.Next() {
		var license models.License
		var hwid, lastHeartbeat sql.NullString
		var activatedAt sql.NullTime

		err := rows.Scan(
			&license.ID, &license.LicenseKey, &license.ProductName,
			&hwid, &license.Status, &license.MaxDevices,
			&license.ExpiresAt, &activatedAt, &license.CreatedAt,
			&license.UserID, &lastHeartbeat,
		)

		if err != nil {
			log.Printf("[Admin] Failed to scan row: %v", err)
			continue
		}

		if hwid.Valid {
			license.HWID = hwid.String
		}
		if activatedAt.Valid {
			license.ActivatedAt = activatedAt.Time
		}

		licenses = append(licenses, license)
	}

	respondJSON(w, map[string]interface{}{
		"licenses": licenses,
		"count":    len(licenses),
	}, http.StatusOK)
}

// HandleGetLicense 获取单个许可证详情
func HandleGetLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	licenseKey := r.URL.Query().Get("key")
	if licenseKey == "" {
		respondError(w, "Missing license key", http.StatusBadRequest)
		return
	}

	var license models.License
	var hwid sql.NullString
	var activatedAt sql.NullTime

	err := database.DB.QueryRow(`
		SELECT id, license_key, product_name, hwid, status, max_devices,
		       expires_at, activated_at, created_at, user_id, order_id, last_heartbeat
		FROM licenses WHERE license_key = ?
	`, licenseKey).Scan(
		&license.ID, &license.LicenseKey, &license.ProductName,
		&hwid, &license.Status, &license.MaxDevices,
		&license.ExpiresAt, &activatedAt, &license.CreatedAt,
		&license.UserID, &license.OrderID, &license.LastHeartbeat,
	)

	if err == sql.ErrNoRows {
		respondError(w, "License not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("[Admin] Database error: %v", err)
		respondError(w, "Database error", http.StatusInternalServerError)
		return
	}

	if hwid.Valid {
		license.HWID = hwid.String
	}
	if activatedAt.Valid {
		license.ActivatedAt = activatedAt.Time
	}

	// 获取激活日志
	logs, _ := getActivationLogs(licenseKey)

	respondJSON(w, map[string]interface{}{
		"license": license,
		"logs":    logs,
	}, http.StatusOK)
}

// HandleUpdateLicense 更新许可证状态
func HandleUpdateLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		LicenseKey string    `json:"license_key"`
		Key        string    `json:"key"` // 兼容前端使用的 key 字段
		Status     string    `json:"status,omitempty"`
		ExpiryDate string    `json:"expiry_date,omitempty"`
		MaxDevices int       `json:"max_devices,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 兼容 key 和 license_key 两种字段
	licenseKey := req.LicenseKey
	if licenseKey == "" {
		licenseKey = req.Key
	}

	if licenseKey == "" {
		respondError(w, "Missing license key", http.StatusBadRequest)
		return
	}

	// 构建动态 UPDATE 语句
	updates := []string{}
	args := []interface{}{}

	if req.Status != "" {
		// 验证状态
		validStatuses := map[string]bool{"active": true, "unused": true, "expired": true, "banned": true}
		if !validStatuses[req.Status] {
			respondError(w, "Invalid status", http.StatusBadRequest)
			return
		}
		updates = append(updates, "status = ?")
		args = append(args, req.Status)
	}

	if req.ExpiryDate != "" {
		// 解析并验证时间格式
		expiryTime, err := time.Parse(time.RFC3339, req.ExpiryDate)
		if err != nil {
			respondError(w, "Invalid expiry_date format", http.StatusBadRequest)
			return
		}
		updates = append(updates, "expires_at = ?")
		args = append(args, expiryTime)
	}

	if req.MaxDevices > 0 {
		updates = append(updates, "max_devices = ?")
		args = append(args, req.MaxDevices)
	}

	if len(updates) == 0 {
		respondError(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// 构建查询，添加 license_key 参数
	args = append(args, licenseKey)

	query := "UPDATE licenses SET " + updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE license_key = ?"

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		log.Printf("[Admin] Failed to update license: %v", err)
		respondError(w, "Failed to update license", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, "License not found", http.StatusNotFound)
		return
	}

	log.Printf("[Admin] Updated license %s", licenseKey)

	respondJSON(w, map[string]string{
		"message": "License updated successfully",
	}, http.StatusOK)
}

// HandleDeleteLicense 删除许可证
func HandleDeleteLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	licenseKey := r.URL.Query().Get("key")
	if licenseKey == "" {
		respondError(w, "Missing license key", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("DELETE FROM licenses WHERE license_key = ?", licenseKey)
	if err != nil {
		log.Printf("[Admin] Failed to delete license: %v", err)
		respondError(w, "Failed to delete license", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, "License not found", http.StatusNotFound)
		return
	}

	log.Printf("[Admin] Deleted license: %s", licenseKey)

	respondJSON(w, map[string]string{
		"message": "License deleted successfully",
	}, http.StatusOK)
}

// HandleGetStats 获取统计数据
func HandleGetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := make(map[string]interface{})

	// 总许可证数
	var total, active, unused, expired, banned int
	database.DB.QueryRow("SELECT COUNT(*) FROM licenses").Scan(&total)
	database.DB.QueryRow("SELECT COUNT(*) FROM licenses WHERE status = 'active'").Scan(&active)
	database.DB.QueryRow("SELECT COUNT(*) FROM licenses WHERE status = 'unused'").Scan(&unused)
	database.DB.QueryRow("SELECT COUNT(*) FROM licenses WHERE status = 'expired'").Scan(&expired)
	database.DB.QueryRow("SELECT COUNT(*) FROM licenses WHERE status = 'banned'").Scan(&banned)

	stats["licenses"] = map[string]int{
		"total":   total,
		"active":  active,
		"unused":  unused,
		"expired": expired,
		"banned":  banned,
	}

	// 用户数
	var userCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	stats["users"] = userCount

	// 今日激活数
	var todayActivations int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM activation_logs
		WHERE action = 'activate' AND success = 1
		AND DATE(created_at) = DATE('now')
	`).Scan(&todayActivations)
	stats["today_activations"] = todayActivations

	respondJSON(w, stats, http.StatusOK)
}

// 辅助函数：获取激活日志
func getActivationLogs(licenseKey string) ([]models.ActivationLog, error) {
	rows, err := database.DB.Query(`
		SELECT id, license_key, hwid, action, ip_address, user_agent, success, error_msg, created_at
		FROM activation_logs WHERE license_key = ?
		ORDER BY created_at DESC LIMIT 50
	`, licenseKey)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []models.ActivationLog{}
	for rows.Next() {
		var log models.ActivationLog
		var hwid, errorMsg sql.NullString

		err := rows.Scan(
			&log.ID, &log.LicenseKey, &hwid, &log.Action,
			&log.IPAddress, &log.UserAgent, &log.Success,
			&errorMsg, &log.CreatedAt,
		)

		if err != nil {
			continue
		}

		if hwid.Valid {
			log.HWID = hwid.String
		}
		if errorMsg.Valid {
			log.ErrorMsg = errorMsg.String
		}

		logs = append(logs, log)
	}

	return logs, nil
}
