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

## 生产环境部署指南

### 📦 服务器端部署

#### 方式一：直接运行（推荐用于开发/测试）

1. **编译服务器端**
```bash
cd server
go build -o license-server
```

2. **配置环境变量**
```bash
export PORT=8080              # 监听端口
export DB_PATH=./licenses.db  # 数据库文件路径
export JWT_SECRET=your-secret-key-here  # JWT密钥（重要！）
```

3. **启动服务器**
```bash
./license-server
```

#### 方式二：Docker 部署（推荐用于生产）

1. **创建 Dockerfile**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN cd server && go mod download
RUN cd server && go build -ldflags="-s -w" -o license-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server/license-server .
EXPOSE 8080
CMD ["./license-server"]
```

2. **构建并运行**
```bash
docker build -t license-server:latest .
docker run -d \
  -p 8080:8080 \
  -e JWT_SECRET=your-secret-key \
  -v $(pwd)/data:/root/data \
  --name license-server \
  license-server:latest
```

#### 方式三：systemd 服务（Linux 服务器）

1. **创建服务文件** `/etc/systemd/system/license-server.service`
```ini
[Unit]
Description=License Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/license-server
ExecStart=/opt/license-server/license-server
Environment="PORT=8080"
Environment="DB_PATH=/var/lib/license-server/licenses.db"
Environment="JWT_SECRET=your-secret-key-here"
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

2. **启动服务**
```bash
sudo systemctl daemon-reload
sudo systemctl enable license-server
sudo systemctl start license-server
sudo systemctl status license-server
```

#### 反向代理配置（Nginx + SSL）

```nginx
server {
    listen 443 ssl http2;
    server_name license.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 🖥️ 客户端集成指南

#### 集成到现有 EXE 程序的三种方式

##### 方式一：作为独立进程（推荐 - 最简单）

**原理**：你的主程序在启动时先调用验证程序，验证通过后才继续运行。

1. **编译验证客户端**
```bash
# Windows 64位
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o drm-validator.exe

# 压缩（可选）
upx --best drm-validator.exe
```

2. **在你的程序中调用**（任何语言都可以）

**C# 示例**：
```csharp
using System;
using System.Diagnostics;

class Program {
    static void Main() {
        // 调用验证程序
        var process = new Process {
            StartInfo = new ProcessStartInfo {
                FileName = "drm-validator.exe",
                UseShellExecute = false,
                RedirectStandardOutput = true,
                CreateNoWindow = true
            }
        };

        process.Start();
        process.WaitForExit();

        if (process.ExitCode != 0) {
            Console.WriteLine("License validation failed!");
            Environment.Exit(1);
        }

        // 验证通过，继续你的程序逻辑
        Console.WriteLine("License valid! Starting main application...");
        RunYourApp();
    }
}
```

**Python 示例**：
```python
import subprocess
import sys

# 调用验证程序
result = subprocess.run(['drm-validator.exe'], capture_output=True)

if result.returncode != 0:
    print("License validation failed!")
    sys.exit(1)

# 验证通过
print("License valid! Starting main application...")
run_your_app()
```

**C++ 示例**：
```cpp
#include <windows.h>
#include <iostream>

int main() {
    STARTUPINFO si = {sizeof(si)};
    PROCESS_INFORMATION pi;

    if (!CreateProcess("drm-validator.exe", NULL, NULL, NULL, FALSE,
                       0, NULL, NULL, &si, &pi)) {
        std::cerr << "Failed to start validator" << std::endl;
        return 1;
    }

    WaitForSingleObject(pi.hProcess, INFINITE);

    DWORD exitCode;
    GetExitCodeProcess(pi.hProcess, &exitCode);
    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);

    if (exitCode != 0) {
        std::cerr << "License validation failed!" << std::endl;
        return 1;
    }

    // 验证通过
    std::cout << "License valid! Starting main application..." << std::endl;
    RunYourApp();
    return 0;
}
```

##### 方式二：作为 DLL/动态链接库

1. **将 Go 代码编译为 C 兼容的 DLL**

修改 `main.go`，导出 C 函数：
```go
package main

import "C"
import (
    "github.com/Lazywords2006/web/auth"
    "github.com/Lazywords2006/web/hwid"
    "github.com/Lazywords2006/web/heartbeat"
)

var monitor *heartbeat.Monitor

//export ValidateLicense
func ValidateLicense(serverURL *C.char, licenseKey *C.char) C.int {
    // 转换C字符串
    url := C.GoString(serverURL)
    key := C.GoString(licenseKey)

    // 执行验证逻辑
    client := auth.NewClient(url)
    hwid, _ := hwid.GetHardwareID()

    if err := client.Activate(key, hwid); err != nil {
        return 0 // 失败
    }

    // 启动心跳
    monitor = heartbeat.NewMonitor(client, 30, 3, 2)
    go monitor.Start()

    return 1 // 成功
}

//export StopValidation
func StopValidation() {
    if monitor != nil {
        // 停止监控（需要添加Stop方法）
    }
}

func main() {}
```

2. **编译为 DLL**
```bash
go build -buildmode=c-shared -o drm-validator.dll
```

3. **在你的程序中调用**

**C# 示例**：
```csharp
using System.Runtime.InteropServices;

class DRMValidator {
    [DllImport("drm-validator.dll")]
    private static extern int ValidateLicense(string serverURL, string licenseKey);

    [DllImport("drm-validator.dll")]
    private static extern void StopValidation();

    public static bool Validate(string serverURL, string key) {
        return ValidateLicense(serverURL, key) == 1;
    }
}

// 使用
if (!DRMValidator.Validate("https://license.yourdomain.com", "YOUR-KEY")) {
    Console.WriteLine("License validation failed!");
    Environment.Exit(1);
}
```

##### 方式三：嵌入到主程序（最隐蔽）

将验证程序作为资源嵌入到你的 EXE 中：

1. **将 drm-validator.exe 转换为 Base64 或二进制资源**
```bash
# PowerShell
$bytes = [System.IO.File]::ReadAllBytes("drm-validator.exe")
[System.Convert]::ToBase64String($bytes) > validator.b64
```

2. **在运行时解压并执行**
```csharp
// 从资源中提取验证器
byte[] validatorBytes = Convert.FromBase64String(Properties.Resources.ValidatorBase64);
string tempPath = Path.Combine(Path.GetTempPath(), "drm-validator.exe");
File.WriteAllBytes(tempPath, validatorBytes);

// 执行验证
var process = Process.Start(tempPath);
process.WaitForExit();

// 清理临时文件
File.Delete(tempPath);

if (process.ExitCode != 0) {
    Environment.Exit(1);
}
```

### 🔧 配置客户端

在你的 EXE 同目录创建 `config.json`：
```json
{
  "server_url": "https://license.yourdomain.com",
  "license_key": "",
  "heartbeat_interval_seconds": 300,
  "max_retries": 3,
  "retry_delay_seconds": 2
}
```

或使用环境变量：
```bash
set LICENSE_SERVER=https://license.yourdomain.com
set LICENSE_KEY=YOUR-KEY-HERE
```

### 🔐 安全加固（生产必须！）

#### 1. 启用 SSL Pinning

编辑 `auth/auth.go:35`：
```go
// 加载证书
certPool := x509.NewCertPool()
cert, _ := ioutil.ReadFile("server.crt")
certPool.AppendCertsFromPEM(cert)

Transport: &http.Transport{
    TLSClientConfig: &tls.Config{
        RootCAs:      certPool,
        MinVersion:   tls.VersionTLS12,
    },
}
```

#### 2. 添加请求签名

编辑 `auth/auth.go:67` 和 `:121`：
```go
import "crypto/hmac"
import "crypto/sha256"

func generateHMAC(data []byte, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write(data)
    return hex.EncodeToString(h.Sum(nil))
}

// 在发送请求前
signature := generateHMAC(jsonData, "your-shared-secret")
req.Header.Set("X-Request-Signature", signature)
req.Header.Set("X-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
```

#### 3. 代码混淆

```bash
# 安装 garble
go install mvdan.cc/garble@latest

# 混淆编译
garble -literals -tiny build -ldflags="-s -w" -o secure-client.exe
```

#### 4. 编译优化

```bash
# 最小化二进制
go build -ldflags="-s -w" -o secure-client.exe

# UPX 压缩
upx --best --ultra-brute secure-client.exe
```

### 📊 管理许可证

#### 使用 API 生成许可证

```bash
# 生成新许可证
curl -X POST https://license.yourdomain.com/api/admin/license \
  -H "Content-Type: application/json" \
  -d '{
    "key": "CUSTOM-KEY-2024-001",
    "max_devices": 3,
    "expiry_date": "2025-12-31T23:59:59Z",
    "note": "Customer: John Doe"
  }'

# 查询许可证
curl "https://license.yourdomain.com/api/admin/license?key=CUSTOM-KEY-2024-001"

# 获取统计
curl "https://license.yourdomain.com/api/admin/stats"
```

#### 管理前端（可选）

将前端文件放到 `server/frontend/` 目录，通过浏览器访问：
```
http://license.yourdomain.com/
```

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
