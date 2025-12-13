# Deploy 部署文件说明

本目录包含服务器部署所需的所有配置文件和脚本。

## 📁 文件说明

### 🚀 核心部署文件

| 文件名 | 说明 | 用途 |
|--------|------|------|
| `quick-deploy.sh` | 一键部署脚本 | 在 Linux 服务器上自动安装和配置 |
| `DEPLOY.md` | 完整部署文档 | 三种部署方案的详细步骤 |
| `QUICKSTART.md` | 快速参考卡片 | 常用命令速查表 |

### ⚙️ 配置文件

| 文件名 | 说明 | 用途 |
|--------|------|------|
| `license-server.service` | systemd 服务配置 | Linux 系统服务配置 |
| `nginx.conf` | Nginx 反向代理配置 | HTTPS 和反向代理设置 |

### 🐳 Docker 文件（在根目录）

| 文件名 | 说明 | 位置 |
|--------|------|------|
| `Dockerfile` | Docker 镜像构建文件 | 项目根目录 |
| `docker-compose.yml` | Docker Compose 配置 | 项目根目录 |

---

## 🚀 快速开始

### 方式一：一键脚本（最简单）

```bash
# 在服务器上执行
cd /path/to/project
chmod +x deploy/quick-deploy.sh
sudo bash deploy/quick-deploy.sh
```

### 方式二：Docker 部署

```bash
# 在服务器上执行
cd /path/to/project
echo "JWT_SECRET=$(openssl rand -hex 32)" > .env
docker-compose up -d
```

### 方式三：手动部署

参考 [DEPLOY.md](DEPLOY.md) 中的详细步骤。

---

## 📖 文档导航

- **新手入门**: 阅读 [QUICKSTART.md](QUICKSTART.md)
- **完整部署**: 阅读 [DEPLOY.md](DEPLOY.md)
- **项目说明**: 阅读主目录 [README.md](../README.md)

---

## 🔧 配置修改

### 修改服务端口

编辑 `license-server.service`:
```ini
Environment="PORT=8080"  # 改为你需要的端口
```

### 修改数据库路径

编辑 `license-server.service`:
```ini
Environment="DB_PATH=/var/lib/license-server/licenses.db"
```

### 修改域名

编辑 `nginx.conf`:
```nginx
server_name license.yourdomain.com;  # 改为你的域名
```

---

## ⚠️ 安全提醒

1. **修改 JWT_SECRET**: 务必生成随机密钥
2. **配置 HTTPS**: 生产环境必须启用 SSL
3. **限制管理接口**: 建议配置 IP 白名单
4. **定期备份**: 设置自动备份数据库

---

## 🆘 故障排查

### 服务启动失败
```bash
sudo journalctl -u license-server -n 50
```

### 端口被占用
```bash
sudo lsof -i :8080
```

### 权限问题
```bash
sudo chown -R license-server:license-server /var/lib/license-server
```

更多问题参考 [DEPLOY.md](DEPLOY.md) 的"常见问题排查"部分。

---

## 📞 获取帮助

- GitHub Issues: https://github.com/Lazywords2006/web/issues
- 项目文档: ../README.md
