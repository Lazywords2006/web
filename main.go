package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"网络验证/auth"
	"网络验证/heartbeat"
	"网络验证/hwid"
)

// Config 应用程序配置
type Config struct {
	ServerURL     string `json:"server_url"`
	LicenseKey    string `json:"license_key,omitempty"`
	HeartbeatSec  int    `json:"heartbeat_interval_seconds"`
	MaxRetries    int    `json:"max_retries"`
	RetryDelaySec int    `json:"retry_delay_seconds"`
}

const (
	configFile = "config.json"
	appVersion = "1.0.0"
)

func main() {
	log.Printf("=== Secure Always-Online Client v%s ===\n", appVersion)

	// 1. 加载配置
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 生成硬件ID
	log.Println("[HWID] Generating hardware identifier...")
	hwID, err := hwid.GetHardwareID()
	if err != nil {
		log.Fatalf("Failed to get hardware ID: %v", err)
	}
	log.Printf("[HWID] Generated: %s\n", hwID[:16]+"...") // 只显示前16位

	// 3. 获取许可证密钥（从配置或用户输入）
	licenseKey := config.LicenseKey
	if licenseKey == "" {
		licenseKey = promptLicenseKey()
	}

	// 4. 创建认证客户端
	authClient := auth.NewClient(config.ServerURL)

	// 5. 激活许可证
	log.Println("[Auth] Activating license...")
	if err := authClient.Activate(licenseKey, hwID); err != nil {
		log.Fatalf("License activation failed: %v", err)
	}
	log.Println("[Auth] ✓ License activated successfully")
	log.Printf("[Auth] Token received: %s...\n", authClient.GetToken()[:20])

	// 6. 启动心跳监控
	log.Println("[Heartbeat] Starting background monitor...")
	hbConfig := &heartbeat.Config{
		Interval:   time.Duration(config.HeartbeatSec) * time.Second,
		MaxRetries: config.MaxRetries,
		RetryDelay: time.Duration(config.RetryDelaySec) * time.Second,
		ErrorCallback: func(err error) {
			log.Printf("[Heartbeat] Critical error callback: %v", err)
		},
	}
	monitor := heartbeat.NewMonitor(authClient, hbConfig)
	monitor.Start()

	// 7. 运行主业务逻辑
	log.Println("\n[App] All security checks passed. Starting main application...")
	log.Println("[App] ========================================")
	RunMainApp()
}

// loadConfig 加载配置文件
func loadConfig() (*Config, error) {
	// 默认配置
	defaultConfig := &Config{
		ServerURL:     "http://localhost:8080",
		HeartbeatSec:  30,
		MaxRetries:    3,
		RetryDelaySec: 2,
	}

	// 尝试读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		// 如果文件不存在，使用默认配置
		if os.IsNotExist(err) {
			log.Printf("[Config] No config file found, using defaults")
			return defaultConfig, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// 使用默认值填充缺失的字段
	if config.ServerURL == "" {
		config.ServerURL = defaultConfig.ServerURL
	}
	if config.HeartbeatSec == 0 {
		config.HeartbeatSec = defaultConfig.HeartbeatSec
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = defaultConfig.MaxRetries
	}
	if config.RetryDelaySec == 0 {
		config.RetryDelaySec = defaultConfig.RetryDelaySec
	}

	log.Printf("[Config] Loaded from %s", configFile)
	log.Printf("[Config] Server: %s", config.ServerURL)
	log.Printf("[Config] Heartbeat: %ds | Retries: %d | Delay: %ds",
		config.HeartbeatSec, config.MaxRetries, config.RetryDelaySec)

	return &config, nil
}

// promptLicenseKey 提示用户输入许可证密钥
func promptLicenseKey() string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your license key: ")
		key, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		key = strings.TrimSpace(key)
		if key != "" {
			return key
		}

		fmt.Println("License key cannot be empty. Please try again.")
	}
}

// RunMainApp 主业务逻辑占位符
// 这里是您实际应用程序的入口点
// 在生产环境中，将此函数替换为您的业务代码
func RunMainApp() {
	log.Println("[App] Main business logic is running...")
	log.Println("[App] (This is a placeholder - replace with your actual application code)")

	// 模拟应用程序运行
	// 在实际应用中，这里应该是您的主要业务逻辑
	// 例如：启动Web服务器、GUI界面、数据处理等

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	counter := 0
	for {
		select {
		case <-ticker.C:
			counter++
			log.Printf("[App] Business logic tick #%d - Application is running normally", counter)

			// 示例：模拟一些工作
			performBusinessLogic(counter)
		}
	}
}

// performBusinessLogic 执行业务逻辑示例
func performBusinessLogic(iteration int) {
	// 这里放置您的实际业务代码
	// 例如：
	// - 处理用户请求
	// - 执行数据分析
	// - 渲染UI
	// - 处理文件
	// etc.

	log.Printf("[App] Processing task #%d...", iteration)

	// 模拟一些工作
	time.Sleep(500 * time.Millisecond)

	log.Printf("[App] Task #%d completed", iteration)
}
