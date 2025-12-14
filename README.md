# 许可证管理系统 (License Management System)

一个基于 Go 语言实现的完整许可证管理系统,包含服务器端 API、Web 管理界面和客户端 SDK。

## ✨ 核心特性

### 服务器端
- ✅ **许可证生成与管理**: 支持单个/批量生成许可证
- ✅ **激活时计算过期**: 许可证在首次激活时才计算过期时间
- ✅ **硬件绑定**: 防止许可证在多台设备上使用
- ✅ **心跳验证**: 实时监控许可证状态
- ✅ **Web 管理界面**: 可视化管理所有许可证
- ✅ **批量操作**: 一次生成多个许可证密钥
- ✅ **SQLite 数据库**: 轻量级、无需额外配置

### 客户端
- ✅ **跨平台硬件ID**: 支持 Windows、Linux、macOS
- ✅ **自动激活**: 一键完成许可证激活
- ✅ **后台心跳**: 自动维持许可证验证状态
- ✅ **强制退出**: 许可证失效时自动终止应用

### 安全特性
- 🔒 基于 JWT 的认证机制
- 🔒 硬件绑定防止密钥共享
- 🔒 过期时间自动验证
- 🔒 封禁功能支持

---

## 📁 项目结构

```
网络验证/
├── server/                    # 服务器端
│   ├── main.go               # 服务器主程序
│   ├── handlers/             # API 处理器
│   │   ├── license.go        # 许可证激活/心跳
│   │   └── admin.go          # 管理 API
│   ├── database/             # 数据库操作
│   │   └── db.go             # SQLite 初始化
│   ├── models/               # 数据模型
│   │   └── models.go
│   ├── utils/                # 工具函数
│   │   └── utils.go
│   ├── frontend/             # Web 管理界面
│   │   ├── login.html        # 登录页面
│   │   ├── index.html        # 管理后台
│   │   └── test.html         # API 测试页面
│   └── licenses.db           # SQLite 数据库
│
├── client/                    # 客户端SDK
│   ├── main.go               # 客户端主程序
│   ├── auth/                 # 认证模块
│   │   └── auth.go
│   ├── hwid/                 # 硬件ID生成
│   │   └── hwid.go
│   └── heartbeat/            # 心跳监控
│       └── heartbeat.go
│
├── 快速集成.sh                # 一键集成脚本
├── 集成指南.md                # 详细集成文档
├── 集成快速参考.md             # 快速参考手册
└── 示例项目.md                # 代码示例
```

---

## 🚀 快速开始

### 方式 1: 独立部署服务器

```bash
# 1. 进入服务器目录
cd server

# 2. 安装依赖
go mod tidy

# 3. 编译
go build -o server main.go

# 4. 运行
./server
```

### 方式 2: 一键集成到你的项目

```bash
# 集成到你的 Go 项目
./快速集成.sh /path/to/your/project
```

### 方式 3: 仅使用 API (跨语言)

服务器独立运行,任何语言通过 HTTP API 调用:

```python
# Python 示例
import requests
response = requests.post(
    "http://localhost:8080/api/activate",
    json={"key": "LICENSE-2025-XXX", "hwid": "device-id"}
)
```

详细集成方式请参考: [集成到EXE的完整指南.md](./docs/集成到EXE的完整指南.md)

---

## 🎯 核心功能

### 1. 许可证生成 (新版逻辑)

#### 特点:
- **按有效期天数设置**: 生成时只设置 `validity_days` (如 365天)
- **激活时计算过期**: 首次激活时计算 `expires_at = 激活时间 + validity_days`
- **灵活管理**: 未激活的许可证没有固定过期日期

#### API: 生成单个许可证

```bash
POST /api/admin/license
Content-Type: application/json

{
  "key": "LICENSE-2025-XXX",
  "max_devices": 3,
  "validity_days": 365,
  "note": "客户备注"
}
```

**响应:**
```json
{
  "license_key": "LICENSE-2025-XXX",
  "max_devices": 3,
  "validity_days": 365,
  "note": "客户备注",
  "status": "unused"
}
```

### 2. 批量生成许可证 (新功能)

```bash
POST /api/admin/licenses/batch
Content-Type: application/json

{
  "count": 10,
  "prefix": "BATCH",
  "max_devices": 2,
  "validity_days": 180,
  "note": "批量测试"
}
```

**响应:**
```json
{
  "success": 10,
  "failed": 0,
  "total": 10,
  "licenses": [
    {"license_key": "BATCH-C-1927-2BAC-0876"},
    {"license_key": "BATCH-6-21FE-6F8C-BBCF"},
    ...
  ],
  "max_devices": 2,
  "validity_days": 180
}
```

### 3. 许可证激活

```bash
POST /api/activate
Content-Type: application/json

{
  "key": "LICENSE-2025-XXX",
  "hwid": "device-hardware-id"
}
```

**首次激活时:**
1. 验证许可证是否有效
2. 计算过期时间: `expires_at = now() + validity_days`
3. 绑定硬件ID
4. 返回 JWT token

**响应:**
```json
{
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 4. 心跳验证

```bash
POST /api/heartbeat
Authorization: Bearer <token>

{
  "key": "LICENSE-2025-XXX",
  "hwid": "device-hardware-id"
}
```

**响应:**
```json
{
  "status": "alive"
}
```

---

## 📊 完整 API 文档

### 客户端 API (公开)

| 端点 | 方法 | 说明 | 请求体 |
|------|------|------|--------|
| `/api/activate` | POST | 激活许可证 | `{key, hwid}` |
| `/api/heartbeat` | POST | 心跳验证 | `{key, hwid}` (需要 token) |

### 管理 API (需要认证)

| 端点 | 方法 | 说明 | 请求体 |
|------|------|------|--------|
| `/api/admin/license` | POST | 生成许可证 | `{key, max_devices, validity_days, note}` |
| `/api/admin/license` | GET | 获取许可证详情 | query: `?key=xxx` |
| `/api/admin/license` | PUT | 更新许可证 | `{key, max_devices?, status?}` |
| `/api/admin/license` | DELETE | 删除许可证 | query: `?key=xxx` |
| `/api/admin/licenses` | GET | 获取许可证列表 | query: `?status=xxx&user_id=xxx` |
| `/api/admin/licenses/batch` | POST | 批量生成 | `{count, prefix, max_devices, validity_days, note}` |
| `/api/admin/stats` | GET | 统计数据 | - |

### Web 管理界面

| 路径 | 说明 |
|------|------|
| `/login.html` | 登录页面 |
| `/index.html` | 管理后台 |
| `/test.html` | API 测试页面 |

**默认登录信息:**
```
用户名: lazywords
密码: w7168855
```

⚠️ **生产环境请务必修改密码!** 编辑 [server/frontend/login.html](server/frontend/login.html) 中的 `validUsers` 对象。

---

## 🛠️ 数据库结构

### licenses 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键 |
| license_key | TEXT | 许可证密钥 (唯一) |
| product_name | TEXT | 产品名称 |
| hwid | TEXT | 绑定的硬件ID |
| status | TEXT | 状态: unused/active/expired/banned |
| max_devices | INTEGER | 最大设备数 |
| validity_days | INTEGER | **有效期天数** (新) |
| expires_at | DATETIME | **过期时间** (激活时设置) |
| activated_at | DATETIME | 激活时间 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |
| user_id | INTEGER | 用户ID (可选) |
| order_id | TEXT | 订单ID (可选) |
| last_heartbeat | DATETIME | 最后心跳时间 |
| note | TEXT | **备注** (新) |

---

## 💡 工作流程

### 旧流程 (已弃用)
```
1. 管理员生成许可证 → 设置绝对过期日期 (如 2026-01-01)
2. 客户激活许可证 → 验证是否过期 (根据绝对日期)
```

### 新流程 (当前版本)
```
1. 管理员生成许可证
   ↓
   设置 validity_days = 365 天
   expires_at = NULL

2. 客户首次激活
   ↓
   计算 expires_at = 当前时间 + 365 天
   保存 activated_at = 当前时间
   绑定 hwid

3. 后续验证
   ↓
   检查 hwid 是否匹配
   检查 expires_at 是否过期
```

**优势:**
- ✅ 许可证可以提前生成,不用担心过期
- ✅ 激活时间更准确反映实际使用时间
- ✅ 灵活的有效期管理

---

## 📱 客户端集成示例

### Go 客户端

```go
package main

import (
    "yourproject/auth"
    "yourproject/hwid"
    "yourproject/heartbeat"
)

func main() {
    // 1. 生成硬件ID
    hwidStr, _ := hwid.GetHardwareID()

    // 2. 创建认证客户端
    client := auth.NewClient("http://localhost:8080")

    // 3. 激活许可证
    err := client.Activate("LICENSE-2025-XXX", hwidStr)
    if err != nil {
        log.Fatal("激活失败:", err)
    }

    // 4. 启动心跳监控
    monitor := heartbeat.NewMonitor(client, 30*time.Second, 3, 2*time.Second)
    monitor.Start()

    // 5. 运行业务逻辑
    RunMainApp()
}
```

### C# 客户端

```csharp
var client = new HttpClient();
var data = new {
    key = "LICENSE-2025-XXX",
    hwid = GetHardwareID()
};
var json = JsonSerializer.Serialize(data);
var response = await client.PostAsync(
    "http://localhost:8080/api/activate",
    new StringContent(json, Encoding.UTF8, "application/json")
);
```

### Python 客户端

```python
import requests

response = requests.post(
    "http://localhost:8080/api/activate",
    json={
        "key": "LICENSE-2025-XXX",
        "hwid": get_hardware_id()
    }
)

if response.json()["status"] == "success":
    token = response.json()["token"]
    # 保存 token 用于后续心跳验证
```

---

## 🔧 配置说明

### 服务器环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | 8080 | 监听端口 |
| `DB_PATH` | ./licenses.db | 数据库文件路径 |
| `JWT_SECRET` | (自动生成) | JWT 签名密钥 |

### 客户端配置文件 (config.json)

```json
{
  "server_url": "http://your-server.com",
  "license_key": "LICENSE-2025-XXX",
  "heartbeat_interval_seconds": 30,
  "max_retries": 3,
  "retry_delay_seconds": 2
}
```

---

## 🚀 部署指南

### 开发环境

```bash
# 启动服务器
cd server
go run main.go

# 访问管理界面
open http://localhost:8080/login.html
```

### 生产环境 (Docker)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY server/ ./
RUN go mod download
RUN go build -ldflags="-s -w" -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/frontend ./frontend
EXPOSE 8080
CMD ["./server"]
```

**构建并运行:**
```bash
docker build -t license-server .
docker run -d -p 8080:8080 \
  -e JWT_SECRET=your-secret \
  -v $(pwd)/data:/root/data \
  --name license-server \
  license-server
```

### 数据库迁移 (如果从旧版本升级)

```bash
sqlite3 server/licenses.db << 'EOF'
ALTER TABLE licenses ADD COLUMN validity_days INTEGER DEFAULT 365;
ALTER TABLE licenses ADD COLUMN note TEXT;
EOF
```

---

## 📚 相关文档

- [集成到EXE的完整指南.md](./docs/集成到EXE的完整指南.md) - 详细的 EXE 集成步骤
- [Python_GUI_使用说明.md](./docs/Python_GUI_使用说明.md) - Python GUI 示例程序使用说明
- [Python GUI 示例代码](./examples/python_gui_example.py) - 完整的示例代码
- [项目结构说明.md](./项目结构说明.md) - 项目文件夹结构说明

---

## 🧪 测试

### 测试服务器是否正常运行

```bash
# 检查服务器状态
curl http://localhost:8080/api/admin/stats

# 期望输出
# {"licenses":{"total":0,"active":0,"unused":0,"expired":0,"banned":0},"today_activations":0,"users":1}
```

### 测试生成许可证

```bash
curl -X POST http://localhost:8080/api/admin/license \
  -H "Content-Type: application/json" \
  -d '{
    "key": "TEST-2025-001",
    "max_devices": 3,
    "validity_days": 365,
    "note": "测试许可证"
  }'
```

### 测试激活许可证

```bash
curl -X POST http://localhost:8080/api/activate \
  -H "Content-Type: application/json" \
  -d '{
    "key": "TEST-2025-001",
    "hwid": "test-device-001"
  }'
```

---

## 🔒 安全建议

### 开发环境
- ✅ 使用 HTTP 即可
- ✅ 使用默认配置快速开发

### 生产环境
- ⚠️ **必须启用 HTTPS** - 使用 Let's Encrypt 或云服务商证书
- ⚠️ **修改默认密码** - 编辑 frontend/login.html
- ⚠️ **使用强 JWT 密钥** - 设置 `JWT_SECRET` 环境变量
- ⚠️ **添加访问频率限制** - 防止暴力破解
- ⚠️ **定期备份数据库** - `licenses.db` 文件
- ⚠️ **使用防火墙** - 仅开放必要端口
- ⚠️ **启用日志监控** - 监控异常激活行为

---

## 🐛 故障排查

### 问题 1: 端口被占用

```bash
# 查找占用进程
lsof -i :8080

# 停止进程
kill -9 <PID>

# 或修改端口
export PORT=9000
```

### 问题 2: 数据库权限错误

```bash
# 确保数据库文件有读写权限
chmod 644 server/licenses.db
```

### 问题 3: 许可证列表为空

检查服务器日志:
```bash
tail -f server/server.log
```

确认数据库是否有数据:
```bash
sqlite3 server/licenses.db "SELECT COUNT(*) FROM licenses;"
```

### 问题 4: 激活失败

常见原因:
- 许可证密钥不存在
- 许可证已被封禁 (status='banned')
- 许可证已在其他设备激活 (hwid 不匹配)
- 服务器地址配置错误

查看详细错误日志:
```bash
tail -50 server/server.log | grep Activate
```

---

## 📊 常见有效期设置

| 套餐类型 | validity_days | 说明 |
|----------|---------------|------|
| 试用版 | 7 | 7天试用 |
| 月卡 | 30 | 1个月 |
| 季卡 | 90 | 3个月 |
| 半年卡 | 180 | 6个月 |
| 年卡 | 365 | 1年 |
| 两年卡 | 730 | 2年 |
| 终身版 | 36500 | 100年 (相当于终身) |

---

## 🎯 使用场景

### 场景 1: 软件销售
```
1. 客户购买后,管理员生成许可证 (validity_days=365)
2. 将许可证密钥发送给客户
3. 客户在软件中输入密钥激活
4. 软件开始计时,365天后到期
```

### 场景 2: 代理商批发
```
1. 代理商购买100个许可证
2. 管理员批量生成 (count=100, validity_days=180)
3. 导出许可证列表给代理商
4. 代理商分发给终端用户
5. 用户激活时才开始计时
```

### 场景 3: 促销活动
```
1. 活动期间批量生成优惠许可证
2. 设置短期有效期 (validity_days=30)
3. 发放给活动参与者
4. 激活后30天到期
```

---

## 📞 技术支持

遇到问题请检查:
1. **服务器日志**: `server/server.log`
2. **数据库数据**: `sqlite3 server/licenses.db`
3. **浏览器控制台**: F12 查看前端错误
4. **网络连接**: 确认服务器可访问

---

## 📝 更新日志

### v2.0.0 (当前版本) - 2025-12-14

**新增功能:**
- ✅ 激活时计算过期时间 (validity_days 模式)
- ✅ 批量生成许可证功能
- ✅ 许可证备注字段
- ✅ 前端界面优化

**改进:**
- 🔧 完善 NULL 值处理
- 🔧 优化日志输出
- 🔧 前端显示逻辑改进 (未激活显示有效期,已激活显示过期时间)

**API 变更:**
- ⚠️ `POST /api/admin/license` 请求参数从 `expiry_date` 改为 `validity_days`
- ⚠️ 新增 `POST /api/admin/licenses/batch` 批量生成接口

**数据库变更:**
- 📊 新增 `validity_days` 字段
- 📊 新增 `note` 字段
- 📊 `expires_at` 在未激活时为 NULL

**兼容性:**
- ✅ 旧许可证正常工作
- ✅ 自动添加默认 `validity_days = 365`
- ✅ 客户端 API 无变更

---

## 📄 许可证

MIT License

---

## 🎉 总结

本项目提供了一套完整的许可证管理解决方案:

1. **灵活的有效期管理** - 激活时计算过期时间
2. **批量操作支持** - 提高许可证发放效率
3. **完善的 Web 界面** - 可视化管理所有许可证
4. **跨平台客户端** - 支持 Windows/Linux/macOS
5. **跨语言 API** - 任何语言都可以集成
6. **开箱即用** - 无需复杂配置

**快速上手:**
```bash
cd server && go run main.go
open http://localhost:8080/login.html
```

开始构建你的许可证管理系统吧! 🚀
