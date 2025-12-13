# 服务器部署完整指南

本指南提供三种完整的服务器部署方案，你可以根据实际情况选择。

---

## 📋 前置要求

### 服务器要求
- **操作系统**: Ubuntu 20.04+, Debian 10+, CentOS 7+
- **内存**: 最少 512MB（推荐 1GB+）
- **磁盘**: 最少 2GB 可用空间
- **网络**: 公网 IP 或域名

### 需要准备
- SSH 访问权限（root 或 sudo）
- 域名（可选，推荐用于生产环境）
- SSL 证书（推荐使用 Let's Encrypt 免费证书）

---

## 🚀 方案一：一键部署脚本（最简单，推荐）

### 1. 上传代码到服务器

```bash
# 在本地执行
scp -r /Users/lazywords/Documents/网络验证 root@your-server-ip:/root/

# 或使用 git
ssh root@your-server-ip
git clone https://github.com/Lazywords2006/web.git
cd web
```

### 2. 运行一键部署脚本

```bash
cd /root/网络验证
chmod +x deploy/quick-deploy.sh
sudo bash deploy/quick-deploy.sh
```

脚本会自动完成：
- ✅ 安装 Go 和依赖
- ✅ 编译服务器程序
- ✅ 创建系统用户
- ✅ 配置 systemd 服务
- ✅ 配置防火墙
- ✅ 启动服务

### 3. 验证部署

```bash
# 检查服务状态
systemctl status license-server

# 测试 API
curl http://localhost:8080/api/admin/stats

# 查看日志
journalctl -u license-server -f
```

### 4. 生成第一个许可证

```bash
curl -X POST http://localhost:8080/api/admin/license \
  -H "Content-Type: application/json" \
  -d '{
    "key": "TEST-KEY-001",
    "max_devices": 5,
    "expiry_date": "2025-12-31T23:59:59Z",
    "note": "测试许可证"
  }'
```

---

## 🐳 方案二：Docker 部署（推荐生产环境）

### 1. 安装 Docker

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com | bash

# 启动 Docker
systemctl start docker
systemctl enable docker
```

### 2. 上传项目文件

```bash
scp -r /Users/lazywords/Documents/网络验证 root@your-server-ip:/opt/
cd /opt/网络验证
```

### 3. 配置环境变量

创建 `.env` 文件：
```bash
cat > .env <<EOF
JWT_SECRET=$(openssl rand -hex 32)
EOF
```

### 4. 构建并启动

#### 方式 A: 使用 docker-compose（推荐）

```bash
# 安装 docker-compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f license-server
```

#### 方式 B: 使用 docker 命令

```bash
# 构建镜像
docker build -t license-server:latest .

# 运行容器
docker run -d \
  --name license-server \
  --restart always \
  -p 8080:8080 \
  -e JWT_SECRET=$(openssl rand -hex 32) \
  -v /opt/license-data:/app/data \
  license-server:latest

# 查看日志
docker logs -f license-server
```

### 5. 验证部署

```bash
# 检查容器状态
docker ps

# 测试 API
curl http://localhost:8080/api/admin/stats
```

---

## ⚙️ 方案三：手动编译部署（完全控制）

### 1. 安装 Go 环境

```bash
# 下载 Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# 安装
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

### 2. 上传并编译项目

```bash
# 上传项目
scp -r /Users/lazywords/Documents/网络验证 root@your-server-ip:/opt/

# 编译
cd /opt/网络验证/server
go mod download
go build -ldflags="-s -w" -o license-server
```

### 3. 创建系统用户

```bash
sudo useradd -r -s /bin/false license-server
```

### 4. 创建目录

```bash
sudo mkdir -p /opt/license-server
sudo mkdir -p /var/lib/license-server
sudo mkdir -p /var/log/license-server

sudo cp /opt/网络验证/server/license-server /opt/license-server/
sudo chown -R license-server:license-server /var/lib/license-server
sudo chown -R license-server:license-server /var/log/license-server
sudo chmod +x /opt/license-server/license-server
```

### 5. 安装 systemd 服务

```bash
# 复制服务文件
sudo cp /opt/网络验证/deploy/license-server.service /etc/systemd/system/

# 编辑服务文件，设置 JWT_SECRET
sudo nano /etc/systemd/system/license-server.service
# 修改这一行: Environment="JWT_SECRET=YOUR-SECRET-KEY-HERE"

# 重新加载并启动
sudo systemctl daemon-reload
sudo systemctl enable license-server
sudo systemctl start license-server

# 查看状态
sudo systemctl status license-server
```

### 6. 配置防火墙

```bash
# UFW (Ubuntu/Debian)
sudo ufw allow 8080/tcp
sudo ufw reload

# FirewallD (CentOS/RHEL)
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

---

## 🌐 配置 Nginx 反向代理 + SSL（推荐）

### 1. 安装 Nginx

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y nginx

# CentOS/RHEL
sudo yum install -y nginx
```

### 2. 安装 Certbot（Let's Encrypt）

```bash
# Ubuntu/Debian
sudo apt install -y certbot python3-certbot-nginx

# CentOS/RHEL
sudo yum install -y certbot python3-certbot-nginx
```

### 3. 配置 Nginx

```bash
# 复制配置文件
sudo cp /opt/网络验证/deploy/nginx.conf /etc/nginx/sites-available/license-server

# 修改域名
sudo nano /etc/nginx/sites-available/license-server
# 将 license.yourdomain.com 替换为你的域名

# 创建软链接
sudo ln -s /etc/nginx/sites-available/license-server /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启 Nginx
sudo systemctl restart nginx
```

### 4. 获取 SSL 证书

```bash
# 使用 Certbot 自动配置
sudo certbot --nginx -d license.yourdomain.com

# 或手动获取证书
sudo certbot certonly --nginx -d license.yourdomain.com
```

### 5. 设置自动续期

```bash
# 添加 cron 任务
echo "0 3 * * * certbot renew --quiet && systemctl reload nginx" | sudo crontab -
```

---

## 📊 管理和监控

### 服务管理命令

```bash
# 查看状态
sudo systemctl status license-server

# 启动服务
sudo systemctl start license-server

# 停止服务
sudo systemctl stop license-server

# 重启服务
sudo systemctl restart license-server

# 查看日志
sudo journalctl -u license-server -f

# 查看最近 100 行日志
sudo journalctl -u license-server -n 100
```

### 数据库管理

```bash
# 备份数据库
sudo cp /var/lib/license-server/licenses.db /backup/licenses-$(date +%Y%m%d).db

# 查看数据库
sudo sqlite3 /var/lib/license-server/licenses.db "SELECT * FROM licenses;"
```

### 性能监控

```bash
# 查看资源使用
sudo systemctl status license-server

# 查看进程
ps aux | grep license-server

# 查看网络连接
sudo netstat -tulpn | grep 8080
```

---

## 🔧 常见问题排查

### 1. 服务启动失败

```bash
# 查看详细错误
sudo journalctl -u license-server -n 50

# 检查端口占用
sudo lsof -i :8080

# 检查文件权限
ls -la /opt/license-server/
ls -la /var/lib/license-server/
```

### 2. 数据库权限问题

```bash
sudo chown -R license-server:license-server /var/lib/license-server
sudo chmod 755 /var/lib/license-server
```

### 3. 防火墙问题

```bash
# 检查防火墙状态
sudo ufw status
sudo firewall-cmd --list-all

# 临时关闭防火墙测试
sudo ufw disable  # 测试后记得重新开启
```

### 4. Nginx 502 错误

```bash
# 检查后端服务是否运行
curl http://localhost:8080/api/admin/stats

# 检查 SELinux（CentOS）
sudo setenforce 0  # 临时关闭测试
```

---

## 🔐 安全建议

### 1. 修改默认端口

编辑服务配置，将 8080 改为其他端口：
```bash
sudo nano /etc/systemd/system/license-server.service
# 修改 Environment="PORT=8080" 为其他端口
sudo systemctl daemon-reload
sudo systemctl restart license-server
```

### 2. 配置 IP 白名单（可选）

在 Nginx 配置中添加：
```nginx
location /api/admin/ {
    allow 192.168.1.0/24;  # 允许的 IP 段
    deny all;              # 拒绝其他
    proxy_pass http://127.0.0.1:8080;
}
```

### 3. 启用 fail2ban 防暴力破解

```bash
sudo apt install fail2ban
# 配置规则...
```

### 4. 定期备份

```bash
# 创建备份脚本
cat > /root/backup-license.sh <<'EOF'
#!/bin/bash
BACKUP_DIR="/backup/license-server"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR
cp /var/lib/license-server/licenses.db $BACKUP_DIR/licenses-$DATE.db
find $BACKUP_DIR -mtime +30 -delete  # 删除30天前的备份
EOF

chmod +x /root/backup-license.sh

# 添加到 crontab
echo "0 2 * * * /root/backup-license.sh" | crontab -
```

---

## 📞 测试部署

### 1. 健康检查

```bash
curl http://your-server-ip:8080/api/admin/stats
```

### 2. 生成测试许可证

```bash
curl -X POST http://your-server-ip:8080/api/admin/license \
  -H "Content-Type: application/json" \
  -d '{
    "key": "TEST-2024-001",
    "max_devices": 3,
    "expiry_date": "2025-12-31T23:59:59Z",
    "note": "测试许可证"
  }'
```

### 3. 激活测试

```bash
curl -X POST http://your-server-ip:8080/api/activate \
  -H "Content-Type: application/json" \
  -d '{
    "key": "TEST-2024-001",
    "hwid": "test-hardware-id-123"
  }'
```

---

## 📚 下一步

1. ✅ 配置域名指向服务器 IP
2. ✅ 安装 SSL 证书
3. ✅ 测试客户端连接
4. ✅ 生成生产许可证
5. ✅ 设置监控告警
6. ✅ 配置定期备份

---

## 🆘 获取帮助

如遇到问题：
1. 查看日志: `journalctl -u license-server -f`
2. 检查 GitHub Issues: https://github.com/Lazywords2006/web/issues
3. 参考 README.md 文档

---

**部署完成后，请妥善保管 JWT_SECRET，这是系统安全的关键！**
