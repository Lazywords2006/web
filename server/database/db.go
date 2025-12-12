package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// 测试连接
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// 创建表
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("[DB] Database initialized successfully")
	return nil
}

// createTables 创建所有数据表
func createTables() error {
	schema := `
	-- 用户表
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		name TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- 产品表
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		currency TEXT DEFAULT 'USD',
		duration INTEGER NOT NULL,
		max_devices INTEGER DEFAULT 1,
		is_active BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- 许可证表
	CREATE TABLE IF NOT EXISTS licenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		license_key TEXT UNIQUE NOT NULL,
		product_name TEXT NOT NULL,
		hwid TEXT,
		status TEXT DEFAULT 'unused',
		max_devices INTEGER DEFAULT 1,
		expires_at DATETIME,
		activated_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER,
		order_id TEXT,
		last_heartbeat DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	-- 订单表
	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		product_name TEXT NOT NULL,
		amount REAL NOT NULL,
		currency TEXT DEFAULT 'USD',
		status TEXT DEFAULT 'pending',
		payment_method TEXT,
		license_key TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (license_key) REFERENCES licenses(license_key)
	);

	-- 激活日志表
	CREATE TABLE IF NOT EXISTS activation_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		license_key TEXT NOT NULL,
		hwid TEXT,
		action TEXT NOT NULL,
		ip_address TEXT,
		user_agent TEXT,
		success BOOLEAN DEFAULT 0,
		error_msg TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- 创建索引
	CREATE INDEX IF NOT EXISTS idx_licenses_key ON licenses(license_key);
	CREATE INDEX IF NOT EXISTS idx_licenses_hwid ON licenses(hwid);
	CREATE INDEX IF NOT EXISTS idx_licenses_status ON licenses(status);
	CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
	CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
	CREATE INDEX IF NOT EXISTS idx_logs_license ON activation_logs(license_key);
	CREATE INDEX IF NOT EXISTS idx_logs_action ON activation_logs(action);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// 插入默认管理员账户（仅在不存在时）
	createDefaultAdmin()

	// 插入示例产品
	createDefaultProducts()

	return nil
}

// createDefaultAdmin 创建默认管理员
func createDefaultAdmin() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = 1").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	// 密码: admin123 (实际应该使用bcrypt加密)
	_, err = DB.Exec(`
		INSERT INTO users (email, password, name, is_admin)
		VALUES ('admin@example.com', 'admin123', 'Admin', 1)
	`)
	if err != nil {
		log.Printf("[DB] Warning: Failed to create default admin: %v", err)
	} else {
		log.Println("[DB] Default admin created: admin@example.com / admin123")
	}
}

// createDefaultProducts 创建默认产品
func createDefaultProducts() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	products := []struct {
		name        string
		description string
		price       float64
		duration    int
		maxDevices  int
	}{
		{"Basic License", "基础版许可证 - 1个月", 9.99, 30, 1},
		{"Standard License", "标准版许可证 - 3个月", 24.99, 90, 2},
		{"Premium License", "高级版许可证 - 1年", 79.99, 365, 5},
		{"Lifetime License", "终身版许可证", 199.99, 36500, 10},
	}

	for _, p := range products {
		_, err := DB.Exec(`
			INSERT INTO products (name, description, price, duration, max_devices)
			VALUES (?, ?, ?, ?, ?)
		`, p.name, p.description, p.price, p.duration, p.maxDevices)
		if err != nil {
			log.Printf("[DB] Warning: Failed to create product %s: %v", p.name, err)
		}
	}

	log.Println("[DB] Default products created")
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
