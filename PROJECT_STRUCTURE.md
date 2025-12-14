# 网络验证系统 - 项目结构

## 📁 项目目录结构

```
网络验证/
├── README.md                    # 主文档（完整的项目说明和使用指南）
├── .gitignore                   # Git 忽略配置
├── go.mod                       # Go 模块定义
├── Makefile                     # 编译和部署脚本
├── Dockerfile                   # Docker 镜像配置
├── docker-compose.yml           # Docker Compose 配置
│
├── server/                      # 🚀 服务器端（核心）
│   ├── main.go                 # 服务器主程序
│   ├── server                  # 编译后的服务器可执行文件
│   ├── licenses.db             # SQLite 数据库
│   │
│   ├── database/               # 数据库操作
│   │   └── db.go              # 数据库初始化和连接
│   │
│   ├── handlers/               # API 处理器
│   │   ├── admin.go           # 管理员 API
│   │   ├── license.go         # 许可证激活和验证
│   │   └── auth.go            # 用户认证
│   │
│   ├── models/                 # 数据模型
│   │   ├── license.go         # 许可证模型
│   │   └── user.go            # 用户模型
│   │
│   ├── utils/                  # 工具函数
│   │   └── license.go         # 许可证生成工具
│   │
│   ├── middleware/             # 中间件
│   │   └── cors.go            # CORS 处理
│   │
│   └── frontend/               # Web 管理界面
│       ├── index.html         # 管理后台
│       └── login.html         # 登录页面
│
├── docs/                        # 📚 文档目录
│   ├── README.md               # 文档说明
│   ├── 项目结构说明.md         # 本文件
│   ├── 集成到EXE的完整指南.md  # EXE 集成指南
│   ├── Python_GUI_使用说明.md  # Python GUI 文档
│   ├── 程序包装器使用说明.md   # 包装器文档
│   │
│   └── lzy_zte_integration/    # lzy_zte 项目集成文档
│       ├── 集成说明_lzy_zte.md
│       └── 快速开始_lzy_zte.txt
│
├── examples/                    # 💡 示例代码
│   ├── README.md               # 示例说明
│   │
│   └── lzy_zte/                # lzy_zte 项目示例
│       ├── launcher_cli.py     # 命令行启动器
│       ├── launcher_wrapper.py # GUI 启动器
│       ├── 启动器_命令行.command # 命令行快速启动
│       ├── 启动器.command       # GUI 快速启动
│       └── lzy_zte_12.10.exe  # 目标程序
│
├── auth/                        # 🔐 客户端认证模块（可复用）
│   └── auth.go                 # 许可证激活和验证
│
├── hwid/                        # 🖥️  硬件ID生成模块（可复用）
│   └── hwid.go                 # 跨平台硬件ID生成
│
├── heartbeat/                   # 💓 心跳监控模块（可复用）
│   └── heartbeat.go            # 后台心跳监控
│
├── deploy/                      # 🌐 部署配置
│   ├── docker/                 # Docker 部署
│   ├── systemd/                # Systemd 服务配置
│   └── nginx/                  # Nginx 配置示例
│
├── frontend/                    # 🎨 前端静态资源（旧版，已弃用）
│   └── ...
│
└── deployment/                  # 📦 其他部署相关
    └── ...
```

## 📝 主要文件说明

### 核心文件

- **README.md** - 项目主文档，包含完整的功能说明、API 文档、使用指南
- **server/main.go** - 服务器主程序入口
- **server/licenses.db** - SQLite 数据库，存储许可证和用户数据

### 文档文件

- **docs/集成到EXE的完整指南.md** - 如何将许可证验证集成到各种语言的 exe 程序
- **docs/Python_GUI_使用说明.md** - Python GUI 示例的详细使用说明
- **docs/lzy_zte_integration/** - lzy_zte 项目的专属集成文档

### 示例代码

- **examples/lzy_zte/launcher_cli.py** - 命令行版启动器（推荐，兼容性最好）
- **examples/lzy_zte/launcher_wrapper.py** - GUI 版启动器（需要支持 Tkinter 的 Python）

### 客户端模块（可复用）

这些模块可以直接复制到您的 Go 项目中使用：

- **auth/** - 许可证激活和验证逻辑
- **hwid/** - 跨平台硬件ID生成
- **heartbeat/** - 后台心跳监控

## 🚀 快速开始

### 1. 启动服务器

```bash
cd server
./server
```

服务器将在 `http://localhost:8080` 运行

### 2. 访问管理后台

```
http://localhost:8080/login.html
用户名: lazywords
密码: w7168855
```

### 3. 生成许可证

在管理后台点击"生成许可证"

### 4. 客户端集成

根据您的需求选择集成方式：

- **有源代码** → 查看 `docs/集成到EXE的完整指南.md`
- **无源代码** → 使用 `examples/lzy_zte/` 中的启动器包装器

## 📚 文档导航

### 核心文档
- [README.md](../README.md) - 主文档
- [集成到EXE的完整指南](集成到EXE的完整指南.md) - 集成教程

### 示例文档
- [lzy_zte 集成说明](lzy_zte_integration/集成说明_lzy_zte.md)
- [lzy_zte 快速开始](lzy_zte_integration/快速开始_lzy_zte.txt)

### 参考文档
- [Python GUI 使用说明](Python_GUI_使用说明.md)
- [程序包装器使用说明](程序包装器使用说明.md)

## 🔧 开发说明

### 编译服务器

```bash
cd server
go build -o server main.go
```

### 运行测试

```bash
cd server
go test ./...
```

### 部署到生产环境

参考 `deploy/` 目录下的配置文件和 README 主文档的部署章节。

## 📊 目录用途总结

| 目录 | 用途 | 重要性 |
|------|------|--------|
| `server/` | 服务器核心代码 | ⭐⭐⭐ 必需 |
| `docs/` | 项目文档 | ⭐⭐⭐ 推荐 |
| `examples/` | 示例代码 | ⭐⭐ 可选 |
| `auth/`, `hwid/`, `heartbeat/` | 可复用客户端模块 | ⭐⭐ 推荐 |
| `deploy/` | 部署配置 | ⭐ 可选 |
| `frontend/` | 旧版前端（已弃用） | ❌ 可删除 |
| `deployment/` | 旧版部署（已弃用） | ❌ 可删除 |

## 💡 提示

- 服务器代码在 `server/` 目录
- 所有文档在 `docs/` 目录
- 示例代码在 `examples/` 目录
- 主文档是根目录的 `README.md`

---

**最后更新**: 2025-12-14
