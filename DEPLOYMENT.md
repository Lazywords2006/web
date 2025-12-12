# 服务器部署指南 - Deployment Guide

本指南将帮助您将许可证管理系统部署到生产服务器。

## 📋 目录

1. [环境要求](#环境要求)
2. [快速部署（Docker）](#快速部署docker)
3. [手动部署](#手动部署)
4. [域名和SSL配置](#域名和ssl配置)
5. [数据库备份](#数据库备份)
6. [监控和维护](#监控和维护)
7. [故障排查](#故障排查)

---

## 环境要求

### 最低配置
- **CPU**: 1核
- **内存**: 512MB
- **磁盘**: 10GB
- **操作系统**: Linux (Ubuntu 20.04+ / CentOS 7+ / Debian 10+)

### 推荐配置
- **CPU**: 2核+
- **内存**: 2GB+
- **磁盘**: 20GB+ SSD
- **操作系统**: Ubuntu 22.04 LTS

### 软件依赖
- Docker 20.10+
- Docker Compose 2.0+
- Git（用于克隆代码）

---

## 快速部署（Docker）

### 1. 安装Docker和Docker Compose

#### Ubuntu/Debian
```bash
# 更新包索引
sudo apt update

# 安装Docker
curl -fsSL https://get.docker.com | sh

# 启动Docker服务
sudo systemctl start docker
sudo systemctl enable docker

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker-compose --version
```

#### CentOS/RHEL
```bash
# 安装Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io

# 启动Docker
sudo systemctl start docker
sudo systemctl enable docker

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. 克隆项目

```bash
cd /opt
git clone https://your-repo/license-system.git
cd license-system
```

### 3. 配置环境

编辑 `deployment/docker-compose.yml`：

```bash
cd deployment
nano docker-compose.yml
```

修改以下配置：
- 端口映射（如果80/443端口被占用）
- 时区设置
- 数据库路径

### 4. 启动服务

```bash
# 构建并启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 查看运行状态
docker-compose ps
```

### 5. 验证部署

```bash
# 检查服务健康状态
curl http://localhost:8080/

# 访问管理后台
curl http://localhost:8080/admin/

# 测试API
curl -X POST http://localhost:8080/api/activate \
  -H "Content-Type: application/json" \
  -d '{"key":"test","hwid":"test"}'
```

### 6. 访问界面

- 管理后台: `http://your-server-ip:8080/admin/`
- 公开销售页: `http://your-server-ip:8080/public/`
- API文档: 查看 `/api/` 端点

---

## 手动部署

### 1. 安装Go环境

```bash
# 下载Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# 解压
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证
go version
```

### 2. 安装依赖

```bash
# 安装SQLite
sudo apt install sqlite3 libsqlite3-dev

# 安装gcc（CGO需要）
sudo apt install build-essential
```

### 3. 编译服务器

```bash
cd server
go mod tidy
go build -o license-server main.go
```

### 4. 创建Systemd服务

创建 `/etc/systemd/system/license-server.service`：

```ini
[Unit]
Description=License Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/license-system/server
ExecStart=/opt/license-system/server/license-server
Restart=always
RestartSec=10
Environment="DB_PATH=/var/lib/license-server/licenses.db"
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 重新加载systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start license-server

# 设置开机自启
sudo systemctl enable license-server

# 查看状态
sudo systemctl status license-server

# 查看日志
sudo journalctl -u license-server -f
```

### 5. 配置Nginx（推荐）

安装Nginx：

```bash
sudo apt install nginx
```

创建配置 `/etc/nginx/sites-available/license-server`：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/license-server /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

---

## 域名和SSL配置

### 1. 配置域名

在您的DNS提供商处添加A记录：

```
类型: A
主机: @
值: your-server-ip
TTL: 3600
```

### 2. 安装SSL证书（Let's Encrypt）

#### 使用Certbot自动配置

```bash
# 安装Certbot
sudo apt install certbot python3-certbot-nginx

# 自动获取证书并配置Nginx
sudo certbot --nginx -d your-domain.com

# 测试自动续期
sudo certbot renew --dry-run
```

#### 使用Docker中的Certbot

```bash
# 创建证书目录
mkdir -p deployment/ssl

# 运行Certbot容器
docker run -it --rm \
  -v $(pwd)/deployment/ssl:/etc/letsencrypt \
  certbot/certbot certonly \
  --standalone \
  -d your-domain.com \
  --email your-email@example.com \
  --agree-tos
```

### 3. 更新Nginx配置

编辑 `deployment/nginx.conf`，将 `your-domain.com` 替换为您的实际域名。

重启Nginx：

```bash
docker-compose restart nginx
```

---

## 数据库备份

### 自动备份脚本

创建 `/opt/backup-license-db.sh`：

```bash
#!/bin/bash
BACKUP_DIR="/backup/license-db"
DATE=$(date +%Y%m%d_%H%M%S)
DB_PATH="/data/licenses.db"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
sqlite3 $DB_PATH ".backup $BACKUP_DIR/licenses_$DATE.db"

# 压缩备份
gzip $BACKUP_DIR/licenses_$DATE.db

# 删除7天前的备份
find $BACKUP_DIR -name "*.gz" -mtime +7 -delete

echo "Backup completed: licenses_$DATE.db.gz"
```

设置执行权限：

```bash
chmod +x /opt/backup-license-db.sh
```

添加到crontab（每天凌晨3点备份）：

```bash
crontab -e

# 添加以下行
0 3 * * * /opt/backup-license-db.sh >> /var/log/license-backup.log 2>&1
```

### Docker环境备份

```bash
# 备份数据卷
docker run --rm \
  -v license-system_license-data:/data \
  -v $(pwd)/backup:/backup \
  alpine tar czf /backup/license-data-$(date +%Y%m%d).tar.gz /data

# 恢复数据卷
docker run --rm \
  -v license-system_license-data:/data \
  -v $(pwd)/backup:/backup \
  alpine tar xzf /backup/license-data-20241213.tar.gz -C /
```

---

## 监控和维护

### 1. 日志监控

```bash
# Docker环境
docker-compose logs -f license-server

# Systemd环境
sudo journalctl -u license-server -f

# Nginx日志
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### 2. 性能监控

使用 `htop` 监控资源使用：

```bash
sudo apt install htop
htop
```

### 3. 数据库维护

```bash
# 进入数据库
sqlite3 /data/licenses.db

# 查看表结构
.schema

# 查询许可证统计
SELECT status, COUNT(*) FROM licenses GROUP BY status;

# 查看最近激活
SELECT * FROM activation_logs ORDER BY created_at DESC LIMIT 10;

# 优化数据库
VACUUM;
```

### 4. 更新服务

```bash
# Docker环境
cd deployment
git pull
docker-compose down
docker-compose up -d --build

# 手动部署
cd server
git pull
go build -o license-server main.go
sudo systemctl restart license-server
```

---

## 故障排查

### 问题1：服务无法启动

**检查步骤：**

```bash
# 查看日志
docker-compose logs license-server

# 检查端口占用
sudo netstat -tlnp | grep 8080

# 检查文件权限
ls -la /data/
```

### 问题2：数据库锁定

**解决方案：**

```bash
# 关闭服务
docker-compose down

# 检查数据库
sqlite3 /data/licenses.db "PRAGMA integrity_check;"

# 重新启动
docker-compose up -d
```

### 问题3：SSL证书错误

**解决方案：**

```bash
# 检查证书有效期
sudo certbot certificates

# 手动续期
sudo certbot renew

# 重启Nginx
docker-compose restart nginx
```

### 问题4：高内存占用

**解决方案：**

```bash
# 限制Docker容器内存
# 编辑 docker-compose.yml
services:
  license-server:
    mem_limit: 512m
    mem_reservation: 256m
```

---

## 安全加固

### 1. 防火墙配置

```bash
# 安装UFW
sudo apt install ufw

# 允许SSH
sudo ufw allow 22/tcp

# 允许HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 启用防火墙
sudo ufw enable

# 查看状态
sudo ufw status
```

### 2. 更改默认密码

登录管理后台后，立即修改默认管理员密码（默认：admin@example.com / admin123）。

### 3. 限制API访问速率

在Nginx配置中添加：

```nginx
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

location /api/ {
    limit_req zone=api_limit burst=20;
    # ... 其他配置
}
```

### 4. 启用HTTPS Only

在服务器代码中强制HTTPS：

```go
// server/main.go
if os.Getenv("FORCE_HTTPS") == "true" {
    http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
} else {
    http.ListenAndServe(":8080", nil)
}
```

---

## 生产环境检查清单

部署前请确认：

- [ ] 已更改默认管理员密码
- [ ] 已配置SSL证书
- [ ] 已设置防火墙规则
- [ ] 已配置自动备份
- [ ] 已设置监控告警
- [ ] 已限制API访问频率
- [ ] 已更新JWT密钥（`server/utils/utils.go`）
- [ ] 已配置域名DNS
- [ ] 已测试所有API端点
- [ ] 已阅读安全建议

---

## 获取帮助

- 📖 文档: [README.md](../README.md)
- 🔧 集成指南: [INTEGRATION.md](../INTEGRATION.md)
- 📧 技术支持: support@yourcompany.com
- 🐛 问题反馈: GitHub Issues

---

**祝您部署顺利！**
