# 集成到 EXE 的完整指南

## 📋 目录
1. [Go 程序集成](#go-程序集成)
2. [其他语言集成](#其他语言集成)
3. [完整示例代码](#完整示例代码)
4. [编译和部署](#编译和部署)

---

## 🚀 Go 程序集成 (推荐)

### 方式 A: 使用现有的客户端模块

你的项目中已经有完整的客户端代码:
- `auth/auth.go` - 许可证激活和验证
- `hwid/hwid.go` - 硬件ID生成
- `heartbeat/heartbeat.go` - 心跳监控

### 步骤 1: 创建你的应用程序

```go
// myapp.go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/Lazywords2006/web/auth"
    "github.com/Lazywords2006/web/heartbeat"
    "github.com/Lazywords2006/web/hwid"
)

func main() {
    fmt.Println("=== 我的应用程序启动 ===")

    // 1. 生成硬件ID
    hwidStr, err := hwid.GetHardwareID()
    if err != nil {
        log.Fatal("❌ 无法获取硬件ID:", err)
    }
    fmt.Printf("🔑 硬件ID: %s...\n", hwidStr[:16])

    // 2. 创建认证客户端
    serverURL := "http://localhost:8080" // 修改为你的服务器地址
    client := auth.NewClient(serverURL)

    // 3. 激活许可证
    fmt.Println("\n📝 请输入许可证密钥:")
    var licenseKey string
    fmt.Scanln(&licenseKey)

    fmt.Println("🔄 正在激活许可证...")
    err = client.Activate(licenseKey, hwidStr)
    if err != nil {
        log.Fatal("❌ 许可证激活失败:", err)
    }
    fmt.Println("✅ 许可证激活成功!")

    // 4. 启动心跳监控
    fmt.Println("💓 启动心跳监控...")
    monitor := heartbeat.NewMonitor(
        client,
        30*time.Second, // 心跳间隔
        3,              // 最大重试次数
        2*time.Second,  // 重试延迟
    )

    // 设置心跳失败回调
    monitor.SetOnFailure(func(err error) {
        log.Printf("⚠️ 心跳失败: %v", err)
    })

    monitor.Start()
    fmt.Println("✅ 心跳监控已启动")

    // 5. 运行你的应用逻辑
    fmt.Println("\n🎉 应用程序正在运行...")
    runYourApplication()
}

// 这里是你的实际应用逻辑
func runYourApplication() {
    // 示例: 持续运行
    for {
        fmt.Println("⚙️  应用正在工作...")
        time.Sleep(10 * time.Second)
    }
}
```

### 步骤 2: 编译成 EXE

#### Windows:
```bash
# 编译 Windows 64位
GOOS=windows GOARCH=amd64 go build -o myapp.exe myapp.go

# 编译 Windows 32位
GOOS=windows GOARCH=386 go build -o myapp_x86.exe myapp.go

# 隐藏控制台窗口 (可选)
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o myapp.exe myapp.go
```

#### macOS:
```bash
GOOS=darwin GOARCH=amd64 go build -o myapp myapp.go
```

#### Linux:
```bash
GOOS=linux GOARCH=amd64 go build -o myapp myapp.go
```

---

## 🔧 高级集成: 自动保存配置

```go
// myapp_advanced.go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"

    "github.com/Lazywords2006/web/auth"
    "github.com/Lazywords2006/web/heartbeat"
    "github.com/Lazywords2006/web/hwid"
)

// Config 配置文件结构
type Config struct {
    ServerURL   string `json:"server_url"`
    LicenseKey  string `json:"license_key"`
    LastChecked string `json:"last_checked"`
}

func main() {
    fmt.Println("=== 我的应用程序启动 ===")

    // 1. 加载配置
    config := loadConfig()

    // 2. 生成硬件ID
    hwidStr, err := hwid.GetHardwareID()
    if err != nil {
        log.Fatal("❌ 无法获取硬件ID:", err)
    }

    // 3. 创建认证客户端
    client := auth.NewClient(config.ServerURL)

    // 4. 检查许可证
    if config.LicenseKey == "" {
        // 首次运行,需要激活
        fmt.Println("📝 首次运行,请输入许可证密钥:")
        fmt.Scanln(&config.LicenseKey)

        fmt.Println("🔄 正在激活许可证...")
        err = client.Activate(config.LicenseKey, hwidStr)
        if err != nil {
            log.Fatal("❌ 许可证激活失败:", err)
        }
        fmt.Println("✅ 许可证激活成功!")

        // 保存配置
        config.LastChecked = time.Now().Format(time.RFC3339)
        saveConfig(config)
    } else {
        // 验证现有许可证
        fmt.Println("🔄 验证许可证...")
        err = client.Activate(config.LicenseKey, hwidStr)
        if err != nil {
            log.Fatal("❌ 许可证验证失败:", err)
        }
        fmt.Println("✅ 许可证有效")
    }

    // 5. 启动心跳监控
    monitor := heartbeat.NewMonitor(client, 30*time.Second, 3, 2*time.Second)
    monitor.SetOnFailure(func(err error) {
        log.Printf("⚠️ 许可证验证失败,应用即将退出: %v", err)
    })
    monitor.Start()

    // 6. 运行应用
    runYourApplication()
}

func loadConfig() *Config {
    configPath := getConfigPath()

    // 如果配置文件不存在,返回默认配置
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        return &Config{
            ServerURL: "http://localhost:8080",
        }
    }

    // 读取配置文件
    data, err := os.ReadFile(configPath)
    if err != nil {
        return &Config{ServerURL: "http://localhost:8080"}
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return &Config{ServerURL: "http://localhost:8080"}
    }

    return &config
}

func saveConfig(config *Config) {
    configPath := getConfigPath()

    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        log.Printf("⚠️ 无法保存配置: %v", err)
        return
    }

    os.WriteFile(configPath, data, 0644)
}

func getConfigPath() string {
    // 获取可执行文件所在目录
    exePath, _ := os.Executable()
    exeDir := filepath.Dir(exePath)
    return filepath.Join(exeDir, "config.json")
}

func runYourApplication() {
    // 你的应用逻辑
    for {
        fmt.Println("⚙️  应用正在工作...")
        time.Sleep(10 * time.Second)
    }
}
```

---

## 🌐 其他语言集成 (C++/C#/Python 等)

如果你的 exe 不是 Go 开发的,可以通过 HTTP API 调用许可证服务器。

### C++ 示例 (使用 cURL)

```cpp
// license_client.cpp
#include <iostream>
#include <string>
#include <curl/curl.h>
#include <json/json.h>

class LicenseClient {
private:
    std::string serverURL;
    std::string token;
    std::string licenseKey;
    std::string hwid;

    // 获取硬件ID (简化版)
    std::string getHardwareID() {
        // 这里需要实现获取CPU/主板序列号的逻辑
        // Windows: 使用 WMI
        // Linux: 读取 /proc/cpuinfo
        // macOS: 使用 IOKit
        return "YOUR-HARDWARE-ID";
    }

    // HTTP POST 请求
    std::string httpPost(const std::string& url, const std::string& jsonData) {
        CURL* curl = curl_easy_init();
        std::string response;

        if (curl) {
            struct curl_slist* headers = NULL;
            headers = curl_slist_append(headers, "Content-Type: application/json");
            if (!token.empty()) {
                std::string auth = "Authorization: Bearer " + token;
                headers = curl_slist_append(headers, auth.c_str());
            }

            curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
            curl_easy_setopt(curl, CURLOPT_POSTFIELDS, jsonData.c_str());
            curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
            curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
            curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);

            CURLcode res = curl_easy_perform(curl);
            curl_easy_cleanup(curl);
            curl_slist_free_all(headers);

            if (res != CURLE_OK) {
                throw std::runtime_error("HTTP request failed");
            }
        }

        return response;
    }

    static size_t WriteCallback(void* contents, size_t size, size_t nmemb, void* userp) {
        ((std::string*)userp)->append((char*)contents, size * nmemb);
        return size * nmemb;
    }

public:
    LicenseClient(const std::string& url) : serverURL(url) {
        hwid = getHardwareID();
    }

    // 激活许可证
    bool activate(const std::string& key) {
        licenseKey = key;

        Json::Value root;
        root["key"] = licenseKey;
        root["hwid"] = hwid;

        Json::StreamWriterBuilder writer;
        std::string jsonData = Json::writeString(writer, root);

        try {
            std::string response = httpPost(serverURL + "/api/activate", jsonData);

            Json::CharReaderBuilder readerBuilder;
            Json::Value jsonResponse;
            std::istringstream s(response);
            std::string errs;

            if (Json::parseFromStream(readerBuilder, s, &jsonResponse, &errs)) {
                if (jsonResponse["status"].asString() == "success") {
                    token = jsonResponse["token"].asString();
                    return true;
                }
            }
        } catch (...) {
            return false;
        }

        return false;
    }

    // 心跳验证
    bool heartbeat() {
        Json::Value root;
        root["key"] = licenseKey;
        root["hwid"] = hwid;

        Json::StreamWriterBuilder writer;
        std::string jsonData = Json::writeString(writer, root);

        try {
            std::string response = httpPost(serverURL + "/api/heartbeat", jsonData);

            Json::CharReaderBuilder readerBuilder;
            Json::Value jsonResponse;
            std::istringstream s(response);
            std::string errs;

            if (Json::parseFromStream(readerBuilder, s, &jsonResponse, &errs)) {
                return jsonResponse["status"].asString() == "alive";
            }
        } catch (...) {
            return false;
        }

        return false;
    }
};

int main() {
    std::cout << "=== 我的应用程序启动 ===" << std::endl;

    LicenseClient client("http://localhost:8080");

    std::string licenseKey;
    std::cout << "请输入许可证密钥: ";
    std::cin >> licenseKey;

    if (client.activate(licenseKey)) {
        std::cout << "✅ 许可证激活成功!" << std::endl;

        // 启动心跳线程
        std::thread([&client]() {
            while (true) {
                std::this_thread::sleep_for(std::chrono::seconds(30));
                if (!client.heartbeat()) {
                    std::cerr << "❌ 许可证验证失败,程序退出" << std::endl;
                    exit(1);
                }
            }
        }).detach();

        // 运行应用
        runYourApplication();
    } else {
        std::cerr << "❌ 许可证激活失败" << std::endl;
        return 1;
    }

    return 0;
}
```

### C# 示例 (.NET)

```csharp
// LicenseClient.cs
using System;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using System.Management;

public class LicenseClient
{
    private readonly string serverURL;
    private string token;
    private string licenseKey;
    private string hwid;
    private readonly HttpClient httpClient;

    public LicenseClient(string url)
    {
        serverURL = url;
        hwid = GetHardwareID();
        httpClient = new HttpClient();
    }

    // 获取硬件ID
    private string GetHardwareID()
    {
        try
        {
            string cpuId = "";
            ManagementObjectSearcher searcher = new ManagementObjectSearcher("SELECT ProcessorId FROM Win32_Processor");
            foreach (ManagementObject obj in searcher.Get())
            {
                cpuId = obj["ProcessorId"].ToString();
                break;
            }
            return cpuId;
        }
        catch
        {
            return "UNKNOWN-HWID";
        }
    }

    // 激活许可证
    public async Task<bool> ActivateAsync(string key)
    {
        licenseKey = key;

        var data = new
        {
            key = licenseKey,
            hwid = hwid
        };

        var json = JsonSerializer.Serialize(data);
        var content = new StringContent(json, Encoding.UTF8, "application/json");

        try
        {
            var response = await httpClient.PostAsync($"{serverURL}/api/activate", content);
            var result = await response.Content.ReadAsStringAsync();
            var jsonDoc = JsonDocument.Parse(result);

            if (jsonDoc.RootElement.TryGetProperty("status", out var status) &&
                status.GetString() == "success")
            {
                token = jsonDoc.RootElement.GetProperty("token").GetString();
                return true;
            }
        }
        catch
        {
            return false;
        }

        return false;
    }

    // 心跳验证
    public async Task<bool> HeartbeatAsync()
    {
        var data = new
        {
            key = licenseKey,
            hwid = hwid
        };

        var json = JsonSerializer.Serialize(data);
        var content = new StringContent(json, Encoding.UTF8, "application/json");

        if (!string.IsNullOrEmpty(token))
        {
            httpClient.DefaultRequestHeaders.Authorization =
                new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", token);
        }

        try
        {
            var response = await httpClient.PostAsync($"{serverURL}/api/heartbeat", content);
            var result = await response.Content.ReadAsStringAsync();
            var jsonDoc = JsonDocument.Parse(result);

            return jsonDoc.RootElement.TryGetProperty("status", out var status) &&
                   status.GetString() == "alive";
        }
        catch
        {
            return false;
        }
    }

    // 启动心跳监控
    public void StartHeartbeat()
    {
        Task.Run(async () =>
        {
            while (true)
            {
                await Task.Delay(30000); // 30秒
                if (!await HeartbeatAsync())
                {
                    Console.WriteLine("❌ 许可证验证失败,程序退出");
                    Environment.Exit(1);
                }
            }
        });
    }
}

// Program.cs
class Program
{
    static async Task Main(string[] args)
    {
        Console.WriteLine("=== 我的应用程序启动 ===");

        var client = new LicenseClient("http://localhost:8080");

        Console.Write("请输入许可证密钥: ");
        string licenseKey = Console.ReadLine();

        if (await client.ActivateAsync(licenseKey))
        {
            Console.WriteLine("✅ 许可证激活成功!");

            // 启动心跳监控
            client.StartHeartbeat();

            // 运行应用
            RunYourApplication();
        }
        else
        {
            Console.WriteLine("❌ 许可证激活失败");
            return;
        }
    }

    static void RunYourApplication()
    {
        while (true)
        {
            Console.WriteLine("⚙️  应用正在工作...");
            Thread.Sleep(10000);
        }
    }
}
```

### Python 示例

```python
# license_client.py
import requests
import time
import hashlib
import uuid
import threading

class LicenseClient:
    def __init__(self, server_url):
        self.server_url = server_url
        self.token = None
        self.license_key = None
        self.hwid = self.get_hardware_id()

    def get_hardware_id(self):
        """获取硬件ID"""
        # 使用 MAC 地址作为硬件ID (简化版)
        mac = ':'.join(['{:02x}'.format((uuid.getnode() >> elements) & 0xff)
                       for elements in range(0, 48, 8)][::-1])
        return hashlib.sha256(mac.encode()).hexdigest()

    def activate(self, license_key):
        """激活许可证"""
        self.license_key = license_key

        try:
            response = requests.post(
                f"{self.server_url}/api/activate",
                json={
                    "key": license_key,
                    "hwid": self.hwid
                }
            )

            if response.status_code == 200:
                data = response.json()
                if data.get("status") == "success":
                    self.token = data.get("token")
                    return True
        except Exception as e:
            print(f"激活失败: {e}")
            return False

        return False

    def heartbeat(self):
        """心跳验证"""
        try:
            headers = {}
            if self.token:
                headers["Authorization"] = f"Bearer {self.token}"

            response = requests.post(
                f"{self.server_url}/api/heartbeat",
                json={
                    "key": self.license_key,
                    "hwid": self.hwid
                },
                headers=headers
            )

            if response.status_code == 200:
                data = response.json()
                return data.get("status") == "alive"
        except Exception as e:
            print(f"心跳失败: {e}")
            return False

        return False

    def start_heartbeat(self):
        """启动心跳监控线程"""
        def heartbeat_loop():
            while True:
                time.sleep(30)  # 30秒
                if not self.heartbeat():
                    print("❌ 许可证验证失败,程序退出")
                    exit(1)

        thread = threading.Thread(target=heartbeat_loop, daemon=True)
        thread.start()

# main.py
if __name__ == "__main__":
    print("=== 我的应用程序启动 ===")

    client = LicenseClient("http://localhost:8080")

    license_key = input("请输入许可证密钥: ")

    if client.activate(license_key):
        print("✅ 许可证激活成功!")

        # 启动心跳监控
        client.start_heartbeat()

        # 运行应用
        while True:
            print("⚙️  应用正在工作...")
            time.sleep(10)
    else:
        print("❌ 许可证激活失败")
```

---

## 📦 编译和部署

### Go 项目编译

```bash
# 1. 确保依赖正确
go mod tidy

# 2. 编译 Windows exe (在任何平台)
GOOS=windows GOARCH=amd64 go build -o myapp.exe

# 3. 减小文件大小 (可选)
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o myapp.exe

# 4. 隐藏控制台窗口 (Windows GUI 应用)
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o myapp.exe

# 5. 添加图标和版本信息 (Windows, 需要 go-winres)
go install github.com/tc-hib/go-winres@latest
go-winres make --product-version=1.0.0 --file-version=1.0.0
go build -ldflags="-H windowsgui -s -w" -o myapp.exe
```

### C++ 项目编译

```bash
# Linux/macOS
g++ -o myapp license_client.cpp -lcurl -ljsoncpp -lpthread

# Windows (MinGW)
g++ -o myapp.exe license_client.cpp -lcurl -ljsoncpp -lws2_32
```

### C# 项目编译

```bash
# 发布单文件 exe
dotnet publish -c Release -r win-x64 --self-contained -p:PublishSingleFile=true

# 结果在 bin/Release/net6.0/win-x64/publish/myapp.exe
```

### Python 打包成 exe

```bash
# 安装 PyInstaller
pip install pyinstaller

# 打包成单个 exe
pyinstaller --onefile --name=myapp main.py

# 带图标
pyinstaller --onefile --icon=app.ico --name=myapp main.py

# 隐藏控制台窗口
pyinstaller --onefile --noconsole --name=myapp main.py
```

---

## 🔧 部署清单

### 客户端部署
1. ✅ 编译好的 exe 文件
2. ✅ config.json (可选,存储配置)
3. ✅ 许可证服务器地址

### 服务器部署
1. ✅ 许可证服务器运行在固定地址
2. ✅ 开放端口 (默认 8080)
3. ✅ HTTPS 证书 (生产环境必须)
4. ✅ 数据库备份策略

---

## 🛡️ 安全建议

### 客户端
1. **混淆代码** - 使用 UPX 或其他工具压缩/混淆 exe
2. **加密通信** - 使用 HTTPS 而不是 HTTP
3. **防止调试** - 添加反调试代码 (可选)
4. **代码签名** - 对 exe 进行数字签名,增加可信度

### 服务器
1. **启用 HTTPS** - Let's Encrypt 免费证书
2. **速率限制** - 防止暴力破解
3. **日志监控** - 监控异常激活行为
4. **定期备份** - 备份 licenses.db

---

## 📊 完整工作流程

```
┌─────────────┐
│ 客户购买    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 管理员生成  │ → 批量生成或单个生成
│ 许可证密钥  │   validity_days = 365
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 发送给客户  │ → 邮件/其他方式
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 客户运行exe │
│ 输入密钥    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ exe 调用    │ → POST /api/activate
│ 激活 API    │   {key, hwid}
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 服务器验证  │ → 计算 expires_at
│ 并激活      │   = now + 365 days
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 返回 token  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ exe 启动    │ → 每 30 秒一次
│ 心跳监控    │   POST /api/heartbeat
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 应用正常    │
│ 运行        │
└─────────────┘
```

---

## 🎯 测试流程

### 1. 测试激活
```bash
# 启动服务器
cd server && ./server

# 生成测试许可证
curl -X POST http://localhost:8080/api/admin/license \
  -H "Content-Type: application/json" \
  -d '{
    "key": "TEST-2025-001",
    "max_devices": 3,
    "validity_days": 365,
    "note": "测试许可证"
  }'

# 运行你的 exe
./myapp.exe
# 输入: TEST-2025-001
```

### 2. 测试心跳
```bash
# 观察服务器日志
tail -f server/server.log

# 应该每 30 秒看到心跳请求
```

### 3. 测试过期
```bash
# 修改数据库中的过期时间为过去
sqlite3 server/licenses.db "UPDATE licenses SET expires_at = '2020-01-01' WHERE license_key='TEST-2025-001';"

# 重新运行 exe,应该激活失败
```

---

## 💡 常见问题

### Q1: 如何在离线环境使用?
A: 可以在激活后允许一段时间的离线使用,修改心跳间隔和重试次数即可。

### Q2: 如何更换硬件?
A: 提供"重置许可证"功能,管理员可以在后台清除 hwid,允许重新激活。

### Q3: 如何防止破解?
A:
- 使用代码混淆
- 添加反调试代码
- 重要逻辑放在服务器端
- 定期更新验证算法

### Q4: 支持多台设备吗?
A: 支持! 在生成许可证时设置 `max_devices` 即可。

---

## 📚 相关资源

- [完整 API 文档](../README.md)
- [服务器部署指南](../DEPLOYMENT.md)
- [集成快速参考](../集成快速参考.md)

---

开始将许可证验证集成到你的 exe 程序吧! 🚀
