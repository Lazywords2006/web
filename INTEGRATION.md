# 客户端集成指南 - Client Integration Guide

本指南将帮助您将许可证验证系统集成到您的实际应用程序中。

## 📋 目录

1. [快速开始](#快速开始)
2. [集成步骤](#集成步骤)
3. [代码示例](#代码示例)
4. [最佳实践](#最佳实践)
5. [常见问题](#常见问题)

---

## 快速开始

### 1. 复制核心模块

将以下目录复制到您的项目中：

```
your-project/
├── auth/      # 复制自: 网络验证/auth/
├── hwid/      # 复制自: 网络验证/hwid/
└── heartbeat/ # 复制自: 网络验证/heartbeat/
```

### 2. 更新go.mod

在您的项目根目录：

```bash
# 初始化Go模块（如果还没有）
go mod init your-project-name

# 添加依赖
go get github.com/golang-jwt/jwt/v5
```

### 3. 修改模块导入路径

在复制的文件中，将所有导入路径从 `网络验证/xxx` 改为 `your-project-name/xxx`

**示例：**

```go
// 修改前
import "网络验证/auth"

// 修改后
import "your-project-name/auth"
```

---

## 集成步骤

### 方式一：独立程序（推荐用于桌面应用）

如果您的应用是一个独立的可执行程序，可以直接使用完整的验证流程：

#### main.go 示例

```go
package main

import (
    "log"
    "time"
    "your-project-name/auth"
    "your-project-name/heartbeat"
    "your-project-name/hwid"
)

func main() {
    log.Println("=== 您的应用程序启动 ===")

    // 1. 生成硬件ID
    hwID, err := hwid.GetHardwareID()
    if err != nil {
        log.Fatalf("无法获取硬件ID: %v", err)
    }

    // 2. 获取许可证密钥（从配置文件或用户输入）
    licenseKey := getLicenseKeyFromConfig() // 您需要实现这个函数

    // 3. 创建认证客户端
    authClient := auth.NewClient("https://your-license-server.com")

    // 4. 激活许可证
    if err := authClient.Activate(licenseKey, hwID); err != nil {
        log.Fatalf("许可证激活失败: %v", err)
    }
    log.Println("✓ 许可证验证成功")

    // 5. 启动心跳监控
    hbConfig := &heartbeat.Config{
        Interval:   30 * time.Second,
        MaxRetries: 3,
        RetryDelay: 2 * time.Second,
    }
    monitor := heartbeat.NewMonitor(authClient, hbConfig)
    monitor.Start()

    // 6. 运行您的实际业务逻辑
    runYourApplication()
}

func runYourApplication() {
    // 在这里编写您的实际程序逻辑
    log.Println("应用程序主逻辑运行中...")

    // 示例：Web服务器
    // http.ListenAndServe(":8080", nil)

    // 示例：GUI应用
    // startGUI()

    // 示例：命令行工具
    // processCommands()

    // 保持程序运行
    select {}
}

func getLicenseKeyFromConfig() string {
    // TODO: 从配置文件读取许可证密钥
    // 或提示用户输入
    return "YOUR-LICENSE-KEY-HERE"
}
```

### 方式二：库模式（推荐用于SDK/库）

如果您的应用是一个库或SDK，可以将许可证验证作为初始化的一部分：

#### license_wrapper.go

```go
package yourlibrary

import (
    "fmt"
    "time"
    "your-project-name/auth"
    "your-project-name/heartbeat"
    "your-project-name/hwid"
)

type LicenseManager struct {
    client  *auth.Client
    monitor *heartbeat.Monitor
    isValid bool
}

// NewLicenseManager 创建许可证管理器
func NewLicenseManager(serverURL, licenseKey string) (*LicenseManager, error) {
    // 生成硬件ID
    hwID, err := hwid.GetHardwareID()
    if err != nil {
        return nil, fmt.Errorf("获取硬件ID失败: %w", err)
    }

    // 创建认证客户端
    client := auth.NewClient(serverURL)

    // 激活许可证
    if err := client.Activate(licenseKey, hwID); err != nil {
        return nil, fmt.Errorf("激活失败: %w", err)
    }

    // 创建许可证管理器
    lm := &LicenseManager{
        client:  client,
        isValid: true,
    }

    // 启动心跳监控
    hbConfig := &heartbeat.Config{
        Interval:   30 * time.Second,
        MaxRetries: 3,
        RetryDelay: 2 * time.Second,
        ErrorCallback: func(err error) {
            lm.isValid = false
        },
    }
    lm.monitor = heartbeat.NewMonitor(client, hbConfig)
    lm.monitor.Start()

    return lm, nil
}

// IsValid 检查许可证是否有效
func (lm *LicenseManager) IsValid() bool {
    return lm.isValid
}

// Shutdown 优雅关闭
func (lm *LicenseManager) Shutdown() {
    if lm.monitor != nil {
        lm.monitor.Stop()
    }
}
```

#### 使用示例

```go
package main

import (
    "log"
    "yourlibrary"
)

func main() {
    // 初始化许可证
    license, err := yourlibrary.NewLicenseManager(
        "https://your-server.com",
        "YOUR-LICENSE-KEY",
    )
    if err != nil {
        log.Fatalf("许可证验证失败: %v", err)
    }
    defer license.Shutdown()

    // 在关键功能前检查许可证
    if !license.IsValid() {
        log.Fatal("许可证已失效")
    }

    // 运行您的功能
    yourLibraryFunction()
}
```

---

## 代码示例

### 示例1：带许可证验证的Web服务器

```go
package main

import (
    "log"
    "net/http"
    "time"
    "your-project/auth"
    "your-project/heartbeat"
    "your-project/hwid"
)

func main() {
    // 初始化许可证验证
    if err := initLicense(); err != nil {
        log.Fatalf("许可证验证失败: %v", err)
    }

    // 启动Web服务器
    http.HandleFunc("/", handleHome)
    http.HandleFunc("/api/data", handleData)

    log.Println("服务器启动在 :8080")
    http.ListenAndServe(":8080", nil)
}

func initLicense() error {
    hwID, err := hwid.GetHardwareID()
    if err != nil {
        return err
    }

    client := auth.NewClient("https://your-server.com")
    if err := client.Activate("YOUR-KEY", hwID); err != nil {
        return err
    }

    monitor := heartbeat.NewMonitor(client, &heartbeat.Config{
        Interval:   30 * time.Second,
        MaxRetries: 3,
        RetryDelay: 2 * time.Second,
    })
    monitor.Start()

    return nil
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to licensed application!"))
}

func handleData(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Protected data"))
}
```

### 示例2：命令行工具

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "time"
    "your-project/auth"
    "your-project/heartbeat"
    "your-project/hwid"
)

var (
    licenseKey = flag.String("license", "", "许可证密钥")
    serverURL  = flag.String("server", "https://your-server.com", "服务器地址")
)

func main() {
    flag.Parse()

    if *licenseKey == "" {
        log.Fatal("请提供许可证密钥: -license YOUR-KEY")
    }

    // 验证许可证
    if err := verifyLicense(*serverURL, *licenseKey); err != nil {
        log.Fatalf("许可证验证失败: %v", err)
    }

    // 运行工具
    runTool()
}

func verifyLicense(serverURL, licenseKey string) error {
    hwID, err := hwid.GetHardwareID()
    if err != nil {
        return fmt.Errorf("获取硬件ID失败: %w", err)
    }

    client := auth.NewClient(serverURL)
    if err := client.Activate(licenseKey, hwID); err != nil {
        return fmt.Errorf("激活失败: %w", err)
    }

    // 后台心跳
    monitor := heartbeat.NewMonitor(client, &heartbeat.Config{
        Interval:   60 * time.Second,
        MaxRetries: 3,
        RetryDelay: 5 * time.Second,
    })
    monitor.Start()

    log.Println("✓ 许可证验证成功")
    return nil
}

func runTool() {
    fmt.Println("工具运行中...")
    // 您的工具逻辑
}
```

---

## 最佳实践

### 1. 配置文件管理

创建 `config.json`：

```json
{
  "server_url": "https://your-license-server.com",
  "license_key": "",
  "heartbeat_interval": 30,
  "max_retries": 3
}
```

读取配置：

```go
type Config struct {
    ServerURL        string `json:"server_url"`
    LicenseKey       string `json:"license_key"`
    HeartbeatInterval int   `json:"heartbeat_interval"`
    MaxRetries       int    `json:"max_retries"`
}

func loadConfig() (*Config, error) {
    data, err := os.ReadFile("config.json")
    if err != nil {
        return nil, err
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

### 2. 错误处理

```go
func main() {
    // 捕获panic以便优雅退出
    defer func() {
        if r := recover(); r != nil {
            log.Printf("程序崩溃: %v", r)
            os.Exit(1)
        }
    }()

    // 验证许可证
    if err := initLicense(); err != nil {
        log.Printf("许可证错误: %v", err)
        // 给用户友好的提示
        fmt.Println("\n❌ 许可证验证失败")
        fmt.Println("请检查您的许可证密钥是否正确")
        fmt.Println("如需帮助，请联系: support@yourcompany.com")
        os.Exit(1)
    }

    // 继续运行程序
    run()
}
```

### 3. 日志管理

```go
func setupLogging() {
    // 写入日志文件
    f, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("打开日志文件失败: %v", err)
    }

    log.SetOutput(f)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}
```

### 4. 许可证密钥保护

```go
// 不要硬编码密钥
// ❌ 错误
const LICENSE_KEY = "AAAA-BBBB-CCCC-DDDD-EEEE"

// ✅ 正确：从环境变量读取
licenseKey := os.Getenv("LICENSE_KEY")

// ✅ 正确：从加密配置读取
licenseKey := readEncryptedConfig()

// ✅ 正确：首次运行时提示用户输入并保存
licenseKey := promptAndSaveLicense()
```

### 5. 多进程/多线程安全

```go
var (
    licenseClient *auth.Client
    licenseMutex  sync.RWMutex
)

func getLicenseClient() *auth.Client {
    licenseMutex.RLock()
    defer licenseMutex.RUnlock()
    return licenseClient
}

func setLicenseClient(client *auth.Client) {
    licenseMutex.Lock()
    defer licenseMutex.Unlock()
    licenseClient = client
}
```

---

## 常见问题

### Q1: 如何在离线环境测试？

**A:** 修改 `auth.Client` 的HTTP超时，或使用Mock服务器：

```go
// 开发模式：禁用许可证验证
if os.Getenv("DEV_MODE") == "true" {
    log.Println("开发模式：跳过许可证验证")
    return nil
}
```

### Q2: 如何自定义心跳间隔？

**A:** 在创建Monitor时修改配置：

```go
monitor := heartbeat.NewMonitor(client, &heartbeat.Config{
    Interval:   5 * time.Minute,  // 5分钟一次
    MaxRetries: 5,                // 重试5次
    RetryDelay: 10 * time.Second, // 10秒延迟
})
```

### Q3: 如何处理许可证到期？

**A:** 服务器会在心跳时检查过期，客户端会自动退出。您也可以添加自定义回调：

```go
hbConfig := &heartbeat.Config{
    Interval: 30 * time.Second,
    ErrorCallback: func(err error) {
        log.Printf("许可证错误: %v", err)
        // 显示友好提示
        showExpirationDialog()
        // 或优雅关闭
        gracefulShutdown()
    },
}
```

### Q4: 如何在多个程序间共享许可证？

**A:** 服务器的 `max_devices` 字段控制同时激活的设备数。同一硬件ID可以重复激活。

### Q5: 如何防止用户绕过许可证验证？

**A:**
1. 使用代码混淆工具（garble）
2. 在关键功能前多次检查许可证状态
3. 将部分核心逻辑放在服务器端
4. 使用二进制加壳工具

---

## 下一步

- 📖 阅读 [服务器部署指南](DEPLOYMENT.md)
- 🔒 了解 [安全加固方案](SECURITY.md)
- 📞 获取技术支持: support@yourcompany.com

---

**提示：** 本示例仅供学习参考，生产环境请务必：
- 使用HTTPS
- 启用SSL Pinning
- 添加请求签名
- 使用bcrypt加密密码
- 定期更新依赖库
