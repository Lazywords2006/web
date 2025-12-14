# 服务器部署指南

## 📋 前提条件

### 服务器要求
- **操作系统**: Linux (Ubuntu 20.04+ / CentOS 7+)
- **内存**: 最少 512MB,推荐 1GB
- **磁盘**: 最少 1GB 可用空间
- **网络**: 可访问的公网 IP 或域名

### 需要安装的软件
- Go 1.19+ (编译用)
- Git
- systemd (服务管理)

---

## 🚀 快速部署

### 1. 克隆项目

```bash
# 在服务器上克隆项目
cd /opt
git clone https://github.com/Lazywords2006/web.git
cd web
```

### 2. 编译服务器

```bash
cd server
go build -o license-server main.go
```

### 3. 配置服务器

服务器使用内置的默认配置,主要参数:
- **端口**: 8080
- **数据库**: SQLite (licenses.db)
- **管理员账号**: lazywords / w7168855

如需修改,编辑 [server/main.go](server/main.go)

### 4. 启动服务器

#### 方法 A: 直接运行(测试用)
```bash
cd server
./license-server
```

#### 方法 B: 后台运行
```bash
nohup ./license-server > server.log 2>&1 &
```

#### 方法 C: Systemd 服务(推荐生产环境)
创建服务文件 `/etc/systemd/system/license-server.service`:

```ini
[Unit]
Description=License Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/web/server
ExecStart=/opt/web/server/license-server
Restart=always
RestartSec=5

# 日志
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl daemon-reload
sudo systemctl enable license-server
sudo systemctl start license-server
sudo systemctl status license-server
```

### 5. 配置防火墙

```bash
# Ubuntu (ufw)
sudo ufw allow 8080/tcp

# CentOS (firewalld)
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

### 6. 验证服务

```bash
# 检查服务状态
curl http://localhost:8080/health

# 访问管理后台
# http://your-server-ip:8080/login.html
```

---

## 🌐 域名和 HTTPS 配置(可选)

### 使用 Nginx 反向代理

#### 1. 安装 Nginx
```bash
# Ubuntu
sudo apt install nginx

# CentOS
sudo yum install nginx
```

#### 2. 配置 Nginx

创建配置文件 `/etc/nginx/sites-available/license`:

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

启用配置:
```bash
sudo ln -s /etc/nginx/sites-available/license /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

#### 3. 配置 HTTPS (Let's Encrypt)

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo systemctl enable certbot.timer
```

---

## 📊 数据库管理

### 数据库位置
```
server/licenses.db
```

### 备份数据库
```bash
# 手动备份
cp licenses.db licenses.db.backup-$(date +%Y%m%d)

# 定时备份(crontab)
0 2 * * * cp /opt/web/server/licenses.db /backup/licenses.db.$(date +\%Y\%m\%d)
```

### 查看数据库
```bash
cd server
sqlite3 licenses.db

# 查看所有许可证
SELECT * FROM licenses;

# 查看活跃许可证
SELECT * FROM licenses WHERE status='active';
```

---

## 🔒 安全建议

### 1. 修改默认管理员密码

在首次部署后,立即修改管理员密码:

```bash
cd server
sqlite3 licenses.db
```

```sql
-- 查看当前用户
SELECT * FROM users;

-- 更新密码(需要自己生成新的哈希)
UPDATE users SET password_hash = 'new_hash' WHERE username = 'lazywords';
```

### 2. 使用环境变量

创建 `.env` 文件(不要提交到 git):
```bash
ADMIN_USERNAME=your_username
ADMIN_PASSWORD=your_secure_password
DB_PATH=/opt/web/server/licenses.db
```

### 3. 限制文件权限

```bash
chmod 600 licenses.db
chmod 700 server
chown www-data:www-data server licenses.db
```

### 4. 配置防火墙规则

只允许必要的端口:
```bash
# 只允许 80, 443, 22
sudo ufw default deny incoming
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

---

## 📝 日志管理

### 查看日志

```bash
# Systemd 日志
sudo journalctl -u license-server -f

# 直接运行的日志
tail -f server.log
```

### 日志轮转

创建 `/etc/logrotate.d/license-server`:

```
/opt/web/server/server.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
    create 0640 www-data www-data
    postrotate
        systemctl reload license-server
    endscript
}
```

---

## 🔧 运维命令

### 服务管理

```bash
# 启动服务
sudo systemctl start license-server

# 停止服务
sudo systemctl stop license-server

# 重启服务
sudo systemctl restart license-server

# 查看状态
sudo systemctl status license-server

# 查看日志
sudo journalctl -u license-server -f
```

### 更新部署

```bash
# 1. 拉取最新代码
cd /opt/web
git pull

# 2. 备份数据库
cp server/licenses.db server/licenses.db.backup

# 3. 重新编译
cd server
go build -o license-server main.go

# 4. 重启服务
sudo systemctl restart license-server

# 5. 验证
sudo systemctl status license-server
```

---

## 🐛 故障排查

### 服务无法启动

```bash
# 检查日志
sudo journalctl -u license-server -n 50

# 检查端口占用
sudo lsof -i :8080

# 检查文件权限
ls -la /opt/web/server/
```

### 数据库错误

```bash
# 检查数据库完整性
cd server
sqlite3 licenses.db "PRAGMA integrity_check;"

# 恢复备份
cp licenses.db.backup licenses.db
sudo systemctl restart license-server
```

### 网络无法访问

```bash
# 检查防火墙
sudo ufw status
sudo firewall-cmd --list-all

# 检查 Nginx
sudo nginx -t
sudo systemctl status nginx

# 检查服务监听
sudo netstat -tlnp | grep 8080
```

---

## 📈 监控建议

### 1. 服务监控

使用 `systemd` 的自动重启功能已经配置。

### 2. 性能监控

```bash
# 安装监控工具
sudo apt install htop iotop

# 查看资源使用
htop
```

### 3. 告警通知(可选)

可以集成:
- Prometheus + Grafana
- Zabbix
- Uptime Robot

---

## 🔄 完整部署脚本

创建 `deploy.sh`:

```bash
#!/bin/bash
set -e

echo "======================================"
echo "  许可证服务器部署脚本"
echo "======================================"
echo ""

# 检查 Go
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装"
    echo "安装方法: https://golang.org/doc/install"
    exit 1
fi

echo "✓ Go 已安装: $(go version)"

# 编译服务器
echo ""
echo "📦 编译服务器..."
cd server
go build -o license-server main.go
chmod +x license-server

echo "✓ 编译完成"

# 创建 systemd 服务
echo ""
echo "🔧 配置 systemd 服务..."
sudo tee /etc/systemd/system/license-server.service > /dev/null <<EOF
[Unit]
Description=License Server
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/license-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable license-server
sudo systemctl start license-server

echo "✓ 服务已配置并启动"

# 显示状态
echo ""
echo "======================================"
echo "  部署完成!"
echo "======================================"
echo ""
echo "服务状态:"
sudo systemctl status license-server --no-pager
echo ""
echo "管理后台: http://$(hostname -I | awk '{print $1}'):8080/login.html"
echo "用户名: lazywords"
echo "密码: w7168855"
echo ""
echo "常用命令:"
echo "  查看日志: sudo journalctl -u license-server -f"
echo "  重启服务: sudo systemctl restart license-server"
echo "  停止服务: sudo systemctl stop license-server"
echo ""
```

使用方法:
```bash
chmod +x deploy.sh
./deploy.sh
```

---

## 📞 技术支持

### 常见问题

**Q: 如何更改端口?**
A: 编辑 `server/main.go` 中的 `PORT` 常量,重新编译并重启

**Q: 如何重置管理员密码?**
A: 直接修改数据库或重新初始化数据库

**Q: 如何备份数据?**
A: 定期备份 `licenses.db` 文件

**Q: 支持集群部署吗?**
A: 当前版本使用 SQLite,不支持集群。如需高可用,建议迁移到 PostgreSQL/MySQL

---

## ✅ 部署检查清单

部署前:
- [ ] 服务器满足最低要求
- [ ] 已安装 Go
- [ ] 已克隆项目

部署后:
- [ ] 服务正常启动
- [ ] 可以访问管理后台
- [ ] 已修改默认密码
- [ ] 已配置防火墙
- [ ] 已配置数据库备份
- [ ] 已配置 HTTPS (如需要)

---

**部署文档版本**: 1.0
**最后更新**: 2025-12-14
**GitHub**: https://github.com/Lazywords2006/web
