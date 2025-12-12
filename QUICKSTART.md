# 快速开始指南 (Quick Start)

## 1. 启动测试服务器

在一个终端窗口中运行：

```bash
# 使用Makefile
make server

# 或直接运行
go run test-server.go
```

您应该看到：
```
=== Mock License Server ===
Starting on http://localhost:8080

Valid test license keys:
  - TEST-LICENSE-KEY-12345
  - DEMO-KEY-ABCDEF
  - VALID-KEY-XYZ789

API Endpoints:
  POST /api/activate   - License activation
  POST /api/heartbeat  - Heartbeat validation
========================================
```

## 2. 配置客户端

复制配置文件：

```bash
# 使用Makefile
make config

# 或手动复制
cp config.json.example config.json
```

## 3. 编译客户端

```bash
# 使用Makefile
make build

# 或直接编译
go build -o secure-client main.go
```

## 4. 运行客户端

```bash
./secure-client
```

当提示输入许可证密钥时，输入以下任一测试密钥：
- `TEST-LICENSE-KEY-12345`
- `DEMO-KEY-ABCDEF`
- `VALID-KEY-XYZ789`

## 预期输出

```
=== Secure Always-Online Client v1.0.0 ===
[Config] Loaded from config.json
[Config] Server: http://localhost:8080
[Config] Heartbeat: 30s | Retries: 3 | Delay: 2s
[HWID] Generating hardware identifier...
[HWID] Generated: 1a2b3c4d5e6f7g8h...
Enter your license key: TEST-LICENSE-KEY-12345
[Auth] Activating license...
[Auth] ✓ License activated successfully
[Auth] Token received: mock-jwt-token-TEST-...
[Heartbeat] Starting background monitor...
[Heartbeat] Monitor started

[App] All security checks passed. Starting main application...
[App] ========================================
[App] Main business logic is running...
[App] (This is a placeholder - replace with your actual application code)
[Heartbeat] OK
[App] Business logic tick #1 - Application is running normally
[App] Processing task #1...
[App] Task #1 completed
[Heartbeat] OK
```

## 测试Kill Switch

### 方法1：停止测试服务器

1. 在运行客户端时，按 `Ctrl+C` 停止测试服务器
2. 等待30秒（心跳间隔）
3. 客户端将尝试3次重试（每次间隔2秒）
4. 所有重试失败后，客户端将自动终止

预期输出：
```
[Heartbeat] Attempt 1/3 failed: network error...
[Heartbeat] Attempt 2/3 failed: network error...
[Heartbeat] Attempt 3/3 failed: network error...
[Heartbeat] CRITICAL: All retries failed...
[KILL SWITCH] Force exit triggered: Heartbeat validation failed...
[KILL SWITCH] Application will terminate immediately
```

### 方法2：使用无效的许可证密钥

1. 输入无效的许可证密钥（例如：`INVALID-KEY`）
2. 激活将失败

预期输出：
```
[Auth] Activating license...
License activation failed: activation failed: Invalid license key
```

## Makefile命令参考

```bash
make build          # 编译主程序
make build-server   # 编译测试服务器
make build-all      # 编译所有
make run            # 运行客户端
make server         # 运行测试服务器
make test           # 运行测试
make clean          # 清理编译产物
make build-cross    # 跨平台编译
make build-release  # 优化编译
make config         # 创建配置文件
make help           # 显示帮助
```

## 下一步

现在您已经验证了系统工作正常，可以：

1. **替换业务逻辑**：修改 `main.go` 中的 `RunMainApp()` 函数
2. **部署服务器**：将 `config.json` 中的 `server_url` 改为您的实际服务器地址
3. **启用安全功能**：参考代码中的TODO注释启用SSL Pinning和请求签名
4. **自定义配置**：调整心跳间隔、重试次数等参数

详细信息请参阅 [README.md](README.md)。
