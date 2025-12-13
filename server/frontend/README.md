# Web 管理界面

## 🌐 访问地址

```
http://YOUR-SERVER-IP:8080/
```

## 🔐 默认账号

| 用户名 | 密码 | 权限 |
|--------|------|------|
| admin  | admin123 | 管理员 |
| root   | root123  | 管理员 |

**⚠️ 重要：首次登录后请立即修改密码！**

## 🔒 修改密码

编辑 `login.html` 文件，修改 `validUsers` 对象：

```javascript
const validUsers = {
    'admin': 'your-new-password',  // 修改密码
    'your-username': 'your-password'  // 添加新用户
};
```

### 使用密码哈希（推荐）

为了更安全，可以使用密码哈希：

```javascript
// 在浏览器控制台生成 SHA-256 哈希
async function hashPassword(password) {
    const encoder = new TextEncoder();
    const data = encoder.encode(password);
    const hash = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hash));
    const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    return hashHex;
}

// 使用
await hashPassword('your-password');
```

然后在 `login.html` 中使用哈希值：

```javascript
const validUsers = {
    'admin': '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918'  // "admin" 的 SHA-256
};

// 验证时也需要哈希
async function hashPassword(password) {
    const encoder = new TextEncoder();
    const data = encoder.encode(password);
    const hash = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hash));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

async function login(event) {
    event.preventDefault();
    const password = document.getElementById('password').value;
    const hashedPassword = await hashPassword(password);
    // 比较 hashedPassword 和存储的哈希值
}
```

## 🔧 功能说明

### 📊 统计仪表盘
实时显示许可证统计数据

### 📝 生成许可证
- 许可证密钥：唯一标识
- 最大设备数：允许激活的设备数量
- 过期时间：许可证有效期
- 备注：客户信息等

### 🔍 查询许可证
根据密钥查询详细信息

### ✏️ 更新许可证
修改设备数、过期时间、状态

### 🗑️ 删除许可证
永久删除许可证（不可恢复）

### 📋 许可证列表
- 查看所有许可证
- 搜索过滤
- 状态展示
- 快速操作

## 🔐 安全建议

### 1. 使用 Nginx 基础认证（推荐）

```bash
# 安装工具
sudo apt install apache2-utils

# 创建密码文件
sudo htpasswd -c /etc/nginx/.htpasswd admin

# Nginx 配置
location / {
    auth_basic "License Admin";
    auth_basic_user_file /etc/nginx/.htpasswd;
    proxy_pass http://127.0.0.1:8080;
}
```

### 2. 使用 HTTPS

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d license.yourdomain.com
```

### 3. 限制 IP 访问

在 Nginx 配置中添加：

```nginx
location /admin/ {
    allow 192.168.1.0/24;  # 允许内网
    deny all;              # 拒绝其他
    proxy_pass http://127.0.0.1:8080;
}
```

### 4. 启用防火墙

```bash
# 只允许特定 IP 访问
sudo ufw allow from 192.168.1.0/24 to any port 8080
```

## 📱 响应式设计

界面支持桌面和移动设备访问。

## 🔄 会话管理

- 登录有效期：24小时
- 超时后自动跳转到登录页面
- 支持手动退出登录

## 🎨 自定义

### 修改主题颜色

编辑 `index.html` 和 `login.html`，修改 CSS 渐变：

```css
background: linear-gradient(135deg, #YOUR-COLOR-1 0%, #YOUR-COLOR-2 100%);
```

### 修改标题

编辑 HTML 中的：

```html
<h1>🔐 你的系统名称</h1>
<p>Your System Description</p>
```

## 🐛 故障排查

### 无法访问管理界面

1. 检查服务器是否运行：
```bash
systemctl status license-server
```

2. 检查防火墙：
```bash
sudo ufw status
```

3. 检查端口监听：
```bash
sudo lsof -i :8080
```

### 登录后跳转失败

清除浏览器缓存或使用无痕模式。

### API 请求失败

检查浏览器控制台（F12）查看具体错误。

## 📞 技术支持

如有问题，请查看：
- 服务器日志：`journalctl -u license-server -f`
- 浏览器控制台（F12）
- GitHub Issues
