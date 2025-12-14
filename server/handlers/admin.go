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
		Key          string `json:"key"`           // 许可证密钥
		MaxDevices   int    `json:"max_devices"`   // 最大设备数
		ValidityDays int    `json:"validity_days"` // 有效期天数
		Note         string `json:"note"`          // 备注
		ProductName  string `json:"product_name"`  // 产品名称(可选)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if req.Key == "" {
		respondError(w, "License key is required", http.StatusBadRequest)
		return
	}

	if req.MaxDevices <= 0 {
		req.MaxDevices = 1
	}

	if req.ValidityDays <= 0 {
		req.ValidityDays = 365 // 默认1年
	}

	if req.ProductName == "" {
		req.ProductName = "Default Product"
	}

	// 插入数据库 (不设置 expires_at,等激活时再计算)
	_, err := database.DB.Exec(`
		INSERT INTO licenses (license_key, product_name, status, max_devices, validity_days, note)
		VALUES (?, ?, 'unused', ?, ?, ?)
	`, req.Key, req.ProductName, req.MaxDevices, req.ValidityDays, req.Note)

	if err != nil {
		log.Printf("[Admin] Failed to insert license: %v", err)
		respondError(w, "Failed to save license", http.StatusInternalServerError)
		return
	}

	log.Printf("[Admin] Generated license: %s (validity: %d days)", req.Key, req.ValidityDays)

	respondJSON(w, map[string]interface{}{
		"license_key":   req.Key,
		"product_name":  req.ProductName,
		"validity_days": req.ValidityDays,
		"max_devices":   req.MaxDevices,
		"note":          req.Note,
		"status":        "unused",
	}, http.StatusOK)
}

// HandleBatchGenerateLicense 批量生成许可证
func HandleBatchGenerateLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Count        int    `json:"count"`         // 生成数量
		Prefix       string `json:"prefix"`        // 密钥前缀
		MaxDevices   int    `json:"max_devices"`   // 最大设备数
		ValidityDays int    `json:"validity_days"` // 有效期天数
		Note         string `json:"note"`          // 备注
		ProductName  string `json:"product_name"`  // 产品名称
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.Count <= 0 || req.Count > 1000 {
		respondError(w, "Count must be between 1 and 1000", http.StatusBadRequest)
		return
	}

	if req.MaxDevices <= 0 {
		req.MaxDevices = 1
	}

	if req.ValidityDays <= 0 {
		req.ValidityDays = 365
	}

	if req.ProductName == "" {
		req.ProductName = "Default Product"
	}

	if req.Prefix == "" {
		req.Prefix = "LICENSE"
	}

	// 生成许可证
	generated := []map[string]interface{}{}
	failed := 0

	for i := 0; i < req.Count; i++ {
		// 生成唯一密钥
		key, err := utils.GenerateLicenseKey()
		if err != nil {
			failed++
			continue
		}

		// 添加前缀
		if req.Prefix != "" && req.Prefix != "LICENSE" {
			key = req.Prefix + "-" + key[8:] // 替换默认前缀
		}

		// 插入数据库
		_, err = database.DB.Exec(`
			INSERT INTO licenses (license_key, product_name, status, max_devices, validity_days, note)
			VALUES (?, ?, 'unused', ?, ?, ?)
		`, key, req.ProductName, req.MaxDevices, req.ValidityDays, req.Note)

		if err != nil {
			log.Printf("[Admin] Failed to insert batch license: %v", err)
			failed++
			continue
		}

		generated = append(generated, map[string]interface{}{
			"license_key": key,
		})
	}

	log.Printf("[Admin] Batch generated %d licenses (failed: %d)", len(generated), failed)

	respondJSON(w, map[string]interface{}{
		"success":       len(generated),
		"failed":        failed,
		"total":         req.Count,
		"licenses":      generated,
		"max_devices":   req.MaxDevices,
		"validity_days": req.ValidityDays,
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

	query := "SELECT id, license_key, product_name, hwid, status, max_devices, validity_days, expires_at, activated_at, created_at, updated_at, user_id, last_heartbeat, note FROM licenses WHERE 1=1"
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

	licenses := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var licenseKey, productName, status string
		var maxDevices, validityDays int
		var hwid, note, lastHeartbeat sql.NullString
		var expiresAt, activatedAt, createdAt, updatedAt sql.NullTime
		var userID sql.NullInt64

		err := rows.Scan(
			&id, &licenseKey, &productName,
			&hwid, &status, &maxDevices, &validityDays,
			&expiresAt, &activatedAt, &createdAt, &updatedAt,
			&userID, &lastHeartbeat, &note,
		)

		if err != nil {
			log.Printf("[Admin] Failed to scan row: %v", err)
			continue
		}

		// 计算已激活设备数
		activatedDevices := 0
		if status == "active" && hwid.Valid && hwid.String != "" {
			activatedDevices = 1
		}

		// 构建许可证对象
		license := map[string]interface{}{
			"id":                id,
			"key":               licenseKey,
			"license_key":       licenseKey,
			"product_name":      productName,
			"status":            status,
			"max_devices":       maxDevices,
			"validity_days":     validityDays,
			"activated_devices": activatedDevices,
		}

		if hwid.Valid {
			license["hwid"] = hwid.String
		}
		if note.Valid {
			license["note"] = note.String
		} else {
			license["note"] = ""
		}
		if expiresAt.Valid {
			license["expires_at"] = expiresAt.Time
			license["expiry_date"] = expiresAt.Time
		}
		if activatedAt.Valid {
			license["activated_at"] = activatedAt.Time
		}
		if createdAt.Valid {
			license["created_at"] = createdAt.Time
		}
		if updatedAt.Valid {
			license["updated_at"] = updatedAt.Time
		}
		if userID.Valid {
			license["user_id"] = userID.Int64
		}
		if lastHeartbeat.Valid {
			license["last_heartbeat"] = lastHeartbeat.String
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
	var hwid, orderID sql.NullString
	var activatedAt, lastHeartbeat sql.NullTime

	err := database.DB.QueryRow(`
		SELECT id, license_key, product_name, hwid, status, max_devices,
		       expires_at, activated_at, created_at, updated_at, user_id, order_id, last_heartbeat
		FROM licenses WHERE license_key = ?
	`, licenseKey).Scan(
		&license.ID, &license.LicenseKey, &license.ProductName,
		&hwid, &license.Status, &license.MaxDevices,
		&license.ExpiresAt, &activatedAt, &license.CreatedAt, &license.UpdatedAt,
		&license.UserID, &orderID, &lastHeartbeat,
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
	if orderID.Valid {
		license.OrderID = orderID.String
	}
	if lastHeartbeat.Valid {
		license.LastHeartbeat = lastHeartbeat.Time
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
