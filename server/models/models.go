package models

import (
	"time"
)

// License 许可证模型
type License struct {
	ID           int64     `json:"id" db:"id"`
	LicenseKey   string    `json:"license_key" db:"license_key"`
	ProductName  string    `json:"product_name" db:"product_name"`
	HWID         string    `json:"hwid,omitempty" db:"hwid"`
	Status       string    `json:"status" db:"status"` // active, expired, banned, unused
	MaxDevices   int       `json:"max_devices" db:"max_devices"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	ActivatedAt  time.Time `json:"activated_at,omitempty" db:"activated_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	UserID       int64     `json:"user_id,omitempty" db:"user_id"`
	OrderID      string    `json:"order_id,omitempty" db:"order_id"`
	LastHeartbeat time.Time `json:"last_heartbeat,omitempty" db:"last_heartbeat"`
}

// User 用户模型
type User struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // 不返回到JSON
	Name      string    `json:"name" db:"name"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Order 订单模型
type Order struct {
	ID            string    `json:"id" db:"id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	ProductName   string    `json:"product_name" db:"product_name"`
	Amount        float64   `json:"amount" db:"amount"`
	Currency      string    `json:"currency" db:"currency"`
	Status        string    `json:"status" db:"status"` // pending, paid, cancelled
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	LicenseKey    string    `json:"license_key,omitempty" db:"license_key"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// ActivationLog 激活日志
type ActivationLog struct {
	ID         int64     `json:"id" db:"id"`
	LicenseKey string    `json:"license_key" db:"license_key"`
	HWID       string    `json:"hwid" db:"hwid"`
	Action     string    `json:"action" db:"action"` // activate, heartbeat, deactivate
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	UserAgent  string    `json:"user_agent" db:"user_agent"`
	Success    bool      `json:"success" db:"success"`
	ErrorMsg   string    `json:"error_msg,omitempty" db:"error_msg"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Product 产品模型
type Product struct {
	ID          int64   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	Currency    string  `json:"currency" db:"currency"`
	Duration    int     `json:"duration" db:"duration"` // 有效期（天）
	MaxDevices  int     `json:"max_devices" db:"max_devices"`
	IsActive    bool    `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
