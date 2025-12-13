#!/bin/bash
# 许可证服务器一键部署脚本
# 适用于 Ubuntu/Debian/CentOS Linux 服务器

set -e

echo "================================"
echo "  许可证服务器一键部署脚本"
echo "================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置变量
PORT="${PORT:-8080}"
DB_PATH="${DB_PATH:-/var/lib/license-server/licenses.db}"
JWT_SECRET="${JWT_SECRET:-$(openssl rand -hex 32)}"
INSTALL_DIR="/opt/license-server"
SERVICE_USER="license-server"

# 检测操作系统
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        VER=$VERSION_ID
    else
        echo -e "${RED}无法检测操作系统${NC}"
        exit 1
    fi
    echo -e "${GREEN}检测到操作系统: $OS $VER${NC}"
}

# 检查是否为 root 用户
check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}请使用 root 权限运行此脚本${NC}"
        echo "使用命令: sudo bash $0"
        exit 1
    fi
}

# 安装依赖
install_dependencies() {
    echo -e "${YELLOW}正在安装依赖...${NC}"

    if [[ "$OS" == "ubuntu" ]] || [[ "$OS" == "debian" ]]; then
        apt-get update
        apt-get install -y wget curl git gcc make
    elif [[ "$OS" == "centos" ]] || [[ "$OS" == "rhel" ]]; then
        yum install -y wget curl git gcc make
    fi

    # 检查 Go 是否已安装
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}Go 未安装，正在安装 Go 1.21...${NC}"
        wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
        rm -rf /usr/local/go
        tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
        echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
        export PATH=$PATH:/usr/local/go/bin
        rm go1.21.0.linux-amd64.tar.gz
    fi

    echo -e "${GREEN}依赖安装完成${NC}"
}

# 创建用户
create_user() {
    if ! id "$SERVICE_USER" &>/dev/null; then
        echo -e "${YELLOW}创建服务用户: $SERVICE_USER${NC}"
        useradd -r -s /bin/false $SERVICE_USER
    fi
}

# 编译服务器
build_server() {
    echo -e "${YELLOW}正在编译服务器...${NC}"

    # 进入项目目录
    cd "$(dirname "$0")/.."

    # 编译服务器端
    cd server
    go mod download
    go build -ldflags="-s -w" -o license-server

    echo -e "${GREEN}编译完成${NC}"
}

# 安装服务器
install_server() {
    echo -e "${YELLOW}正在安装服务器到 $INSTALL_DIR${NC}"

    # 创建目录
    mkdir -p $INSTALL_DIR
    mkdir -p /var/lib/license-server
    mkdir -p /var/log/license-server

    # 复制文件
    cp server/license-server $INSTALL_DIR/

    # 设置权限
    chown -R $SERVICE_USER:$SERVICE_USER /var/lib/license-server
    chown -R $SERVICE_USER:$SERVICE_USER /var/log/license-server
    chmod +x $INSTALL_DIR/license-server

    echo -e "${GREEN}安装完成${NC}"
}

# 创建 systemd 服务
create_systemd_service() {
    echo -e "${YELLOW}创建 systemd 服务...${NC}"

    cat > /etc/systemd/system/license-server.service <<EOF
[Unit]
Description=License Server
After=network.target

[Service]
Type=simple
User=$SERVICE_USER
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/license-server
Environment="PORT=$PORT"
Environment="DB_PATH=$DB_PATH"
Environment="JWT_SECRET=$JWT_SECRET"
Restart=always
RestartSec=10
StandardOutput=append:/var/log/license-server/server.log
StandardError=append:/var/log/license-server/error.log

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable license-server

    echo -e "${GREEN}systemd 服务创建完成${NC}"
}

# 配置防火墙
configure_firewall() {
    echo -e "${YELLOW}配置防火墙...${NC}"

    if command -v ufw &> /dev/null; then
        ufw allow $PORT/tcp
        echo -e "${GREEN}UFW 防火墙规则已添加${NC}"
    elif command -v firewall-cmd &> /dev/null; then
        firewall-cmd --permanent --add-port=$PORT/tcp
        firewall-cmd --reload
        echo -e "${GREEN}FirewallD 规则已添加${NC}"
    else
        echo -e "${YELLOW}未检测到防火墙，请手动配置${NC}"
    fi
}

# 启动服务
start_service() {
    echo -e "${YELLOW}启动服务...${NC}"
    systemctl start license-server
    sleep 2

    if systemctl is-active --quiet license-server; then
        echo -e "${GREEN}服务启动成功！${NC}"
    else
        echo -e "${RED}服务启动失败，请检查日志${NC}"
        journalctl -u license-server -n 50
        exit 1
    fi
}

# 显示信息
show_info() {
    echo ""
    echo "================================"
    echo -e "${GREEN}部署完成！${NC}"
    echo "================================"
    echo ""
    echo "服务器信息："
    echo "  - 监听端口: $PORT"
    echo "  - 数据库路径: $DB_PATH"
    echo "  - JWT密钥: $JWT_SECRET"
    echo "  - 安装目录: $INSTALL_DIR"
    echo ""
    echo "管理命令："
    echo "  - 查看状态: systemctl status license-server"
    echo "  - 启动服务: systemctl start license-server"
    echo "  - 停止服务: systemctl stop license-server"
    echo "  - 重启服务: systemctl restart license-server"
    echo "  - 查看日志: journalctl -u license-server -f"
    echo ""
    echo "API 端点："
    echo "  - 激活: http://YOUR-SERVER-IP:$PORT/api/activate"
    echo "  - 心跳: http://YOUR-SERVER-IP:$PORT/api/heartbeat"
    echo "  - 管理: http://YOUR-SERVER-IP:$PORT/api/admin/license"
    echo ""
    echo -e "${YELLOW}重要：请保存 JWT_SECRET，客户端需要使用！${NC}"
    echo -e "${YELLOW}JWT_SECRET: $JWT_SECRET${NC}"
    echo ""
    echo "下一步："
    echo "1. 配置域名和 SSL 证书（推荐使用 Nginx 反向代理）"
    echo "2. 生成第一个许可证："
    echo "   curl -X POST http://localhost:$PORT/api/admin/license \\"
    echo "     -H 'Content-Type: application/json' \\"
    echo "     -d '{\"key\":\"TEST-KEY-001\",\"max_devices\":5,\"expiry_date\":\"2025-12-31T23:59:59Z\"}'"
    echo ""
}

# 主函数
main() {
    check_root
    detect_os
    install_dependencies
    create_user
    build_server
    install_server
    create_systemd_service
    configure_firewall
    start_service
    show_info
}

# 运行主函数
main
