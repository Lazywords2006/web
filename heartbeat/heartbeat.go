package heartbeat

import (
	"fmt"
	"log"
	"os"
	"time"
)

// AuthClient 认证客户端接口（用于解耦）
type AuthClient interface {
	Heartbeat() error
}

// Monitor 心跳监控器
type Monitor struct {
	client        AuthClient
	interval      time.Duration
	maxRetries    int
	retryDelay    time.Duration
	stopChan      chan struct{}
	errorCallback func(error)
}

// Config 心跳监控配置
type Config struct {
	Interval      time.Duration // 心跳间隔（默认30秒）
	MaxRetries    int           // 最大重试次数（默认3次）
	RetryDelay    time.Duration // 重试延迟（默认2秒）
	ErrorCallback func(error)   // 错误回调（可选）
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Interval:   30 * time.Second,
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
	}
}

// NewMonitor 创建新的心跳监控器
func NewMonitor(client AuthClient, config *Config) *Monitor {
	if config == nil {
		config = DefaultConfig()
	}

	return &Monitor{
		client:        client,
		interval:      config.Interval,
		maxRetries:    config.MaxRetries,
		retryDelay:    config.RetryDelay,
		stopChan:      make(chan struct{}),
		errorCallback: config.ErrorCallback,
	}
}

// Start 启动心跳监控（在后台Goroutine中运行）
// 此方法会立即返回，心跳检查在后台持续进行
// 如果心跳失败达到最大重试次数，将触发ForceExit终止程序
func (m *Monitor) Start() {
	log.Println("[Heartbeat] Monitor started")

	go m.run()
}

// run 心跳监控主循环（私有方法）
func (m *Monitor) run() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := m.sendHeartbeatWithRetry(); err != nil {
				log.Printf("[Heartbeat] CRITICAL: All retries failed - %v", err)

				// 调用错误回调（如果配置）
				if m.errorCallback != nil {
					m.errorCallback(err)
				}

				// 触发强制退出
				ForceExit(fmt.Sprintf("Heartbeat validation failed: %v", err))
			} else {
				log.Println("[Heartbeat] OK")
			}

		case <-m.stopChan:
			log.Println("[Heartbeat] Monitor stopped")
			return
		}
	}
}

// sendHeartbeatWithRetry 发送心跳并在失败时重试
func (m *Monitor) sendHeartbeatWithRetry() error {
	var lastErr error

	for attempt := 1; attempt <= m.maxRetries; attempt++ {
		err := m.client.Heartbeat()
		if err == nil {
			// 心跳成功
			if attempt > 1 {
				log.Printf("[Heartbeat] Recovered after %d attempts", attempt)
			}
			return nil
		}

		lastErr = err
		log.Printf("[Heartbeat] Attempt %d/%d failed: %v", attempt, m.maxRetries, err)

		// 如果不是最后一次尝试，等待后重试
		if attempt < m.maxRetries {
			time.Sleep(m.retryDelay)
		}
	}

	return fmt.Errorf("heartbeat failed after %d attempts: %w", m.maxRetries, lastErr)
}

// Stop 停止心跳监控
func (m *Monitor) Stop() {
	close(m.stopChan)
}

// ForceExit 强制终止程序
// 这是一个紧急退出函数，用于在许可证验证失败时立即终止应用程序
// 注意：此函数会调用os.Exit(1)，不会执行defer语句
func ForceExit(reason string) {
	log.Printf("[KILL SWITCH] Force exit triggered: %s", reason)
	log.Println("[KILL SWITCH] Application will terminate immediately")

	// TODO: 生产环境中可以添加以下安全措施：
	// 1. 清理敏感数据（内存清零）
	// 2. 关闭所有网络连接
	// 3. 写入审计日志到远程服务器
	// 4. 触发防篡改机制（如删除临时文件、锁定资源等）

	// 立即终止程序
	// Exit code 1 表示异常终止
	os.Exit(1)
}

// GracefulShutdown 优雅关闭（正常退出时使用）
// 与ForceExit不同，此函数允许defer语句执行
func GracefulShutdown(reason string) {
	log.Printf("[Shutdown] Graceful shutdown initiated: %s", reason)

	// 可以在这里添加清理逻辑
	// 例如：保存状态、关闭数据库连接等

	os.Exit(0)
}
