# 快速部署指令卡

## 🚀 最快部署（一键脚本）

```bash
# 1. 上传代码到服务器
git clone https://github.com/Lazywords2006/web.git
cd web

# 2. 运行一键部署
chmod +x deploy/quick-deploy.sh
sudo bash deploy/quick-deploy.sh

# 3. 完成！
```

---

## 🐳 Docker 快速部署

```bash
# 1. 克隆代码
git clone https://github.com/Lazywords2006/web.git
cd web

# 2. 配置密钥
echo "JWT_SECRET=$(openssl rand -hex 32)" > .env

# 3. 启动
docker-compose up -d

# 4. 查看状态
docker-compose ps
docker-compose logs -f
```

---

## 📋 常用管理命令

```bash
# 服务管理
systemctl status license-server    # 查看状态
systemctl start license-server     # 启动
systemctl stop license-server      # 停止
systemctl restart license-server   # 重启

# 日志查看
journalctl -u license-server -f    # 实时日志
journalctl -u license-server -n 100  # 最近100行

# 测试 API
curl http://localhost:8080/api/admin/stats
```

---

## 🔑 生成许可证

```bash
curl -X POST http://YOUR-SERVER:8080/api/admin/license \
  -H "Content-Type: application/json" \
  -d '{
    "key": "YOUR-KEY-001",
    "max_devices": 5,
    "expiry_date": "2025-12-31T23:59:59Z",
    "note": "客户名称"
  }'
```

---

## 🔍 查询许可证

```bash
# 查询指定许可证
curl "http://YOUR-SERVER:8080/api/admin/license?key=YOUR-KEY-001"

# 列出所有许可证
curl "http://YOUR-SERVER:8080/api/admin/licenses"

# 获取统计信息
curl "http://YOUR-SERVER:8080/api/admin/stats"
```

---

## 🔐 SSL 配置（Let's Encrypt）

```bash
# 1. 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 2. 获取证书
sudo certbot --nginx -d license.yourdomain.com

# 3. 自动续期
echo "0 3 * * * certbot renew --quiet" | sudo crontab -
```

---

## 💾 数据库备份

```bash
# 手动备份
sudo cp /var/lib/license-server/licenses.db \
  /backup/licenses-$(date +%Y%m%d).db

# 自动备份脚本
cat > /root/backup.sh <<'EOF'
#!/bin/bash
cp /var/lib/license-server/licenses.db \
  /backup/licenses-$(date +%Y%m%d).db
EOF

chmod +x /root/backup.sh
echo "0 2 * * * /root/backup.sh" | crontab -
```

---

## 🐛 故障排查

```bash
# 检查服务状态
systemctl status license-server

# 查看错误日志
journalctl -u license-server -n 50

# 检查端口
sudo lsof -i :8080

# 检查防火墙
sudo ufw status
sudo firewall-cmd --list-all

# 测试本地连接
curl http://localhost:8080/api/admin/stats
```

---

## 🌐 客户端配置

客户端 `config.json`:
```json
{
  "server_url": "https://license.yourdomain.com",
  "license_key": "",
  "heartbeat_interval_seconds": 300,
  "max_retries": 3,
  "retry_delay_seconds": 2
}
```

编译客户端:
```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o validator.exe
```

---

## 📞 重要端口

- `8080` - 服务器监听端口（HTTP）
- `80` - Nginx HTTP（可选）
- `443` - Nginx HTTPS（推荐）

---

## ⚠️ 重要提醒

1. **保存 JWT_SECRET** - 系统安全关键
2. **配置 HTTPS** - 生产环境必须
3. **定期备份** - 每天自动备份数据库
4. **监控日志** - 及时发现异常
5. **更新系统** - 保持系统安全补丁

---

完整文档：[deploy/DEPLOY.md](deploy/DEPLOY.md)
