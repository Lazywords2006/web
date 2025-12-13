#!/bin/bash

# 许可证管理系统 - 快速集成脚本
# 使用方法: ./快速集成.sh /path/to/your/project

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目源目录
SOURCE_DIR="/Users/lazywords/Documents/网络验证"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查参数
if [ $# -eq 0 ]; then
    print_error "请提供目标项目路径"
    echo "用法: $0 /path/to/your/project"
    echo ""
    echo "示例:"
    echo "  $0 ~/MyProject              # 集成到你的项目"
    echo "  $0 ~/Desktop/LicenseServer  # 创建独立服务器"
    exit 1
fi

TARGET_DIR="$1"

# 显示欢迎信息
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}  许可证管理系统 - 快速集成工具${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
print_info "源目录: $SOURCE_DIR"
print_info "目标目录: $TARGET_DIR"
echo ""

# 询问集成方式
echo "请选择集成方式:"
echo "  1) 完整独立部署 (推荐)"
echo "  2) 集成到现有 Go 项目"
echo "  3) 仅复制核心文件"
echo ""
read -p "请输入选项 (1-3): " choice

case $choice in
    1)
        print_info "开始完整独立部署..."

        # 创建目标目录
        mkdir -p "$TARGET_DIR"
        cd "$TARGET_DIR"

        # 复制所有文件
        print_info "复制服务器文件..."
        cp -r "$SOURCE_DIR/server" .

        print_info "复制文档..."
        cp "$SOURCE_DIR/README.md" .
        cp "$SOURCE_DIR/集成指南.md" .

        # 初始化 Go 模块
        print_info "初始化 Go 模块..."
        cd server
        if [ ! -f "go.mod" ]; then
            go mod init license-server
            go get github.com/mattn/go-sqlite3
            go mod tidy
        fi

        # 编译
        print_info "编译服务器..."
        go build -o server main.go

        print_success "部署完成!"
        echo ""
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo -e "${GREEN}启动服务器:${NC}"
        echo "  cd $TARGET_DIR/server"
        echo "  ./server"
        echo ""
        echo -e "${GREEN}访问管理界面:${NC}"
        echo "  http://localhost:8080/login.html"
        echo "  用户名: lazywords"
        echo "  密码: w7168855"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ;;

    2)
        print_info "集成到现有 Go 项目..."

        # 检查目标目录是否存在
        if [ ! -d "$TARGET_DIR" ]; then
            print_error "目标目录不存在: $TARGET_DIR"
            exit 1
        fi

        cd "$TARGET_DIR"

        # 创建 license 目录结构
        print_info "创建目录结构..."
        mkdir -p license/{handlers,models,database,utils,frontend}

        # 复制后端代码
        print_info "复制后端代码..."
        cp "$SOURCE_DIR/server/handlers/admin.go" license/handlers/
        cp "$SOURCE_DIR/server/handlers/license.go" license/handlers/
        cp "$SOURCE_DIR/server/models/models.go" license/models/
        cp "$SOURCE_DIR/server/database/db.go" license/database/
        cp "$SOURCE_DIR/server/utils/utils.go" license/utils/

        # 复制前端
        print_info "复制前端界面..."
        cp -r "$SOURCE_DIR/server/frontend/"* license/frontend/

        # 复制集成文档
        cp "$SOURCE_DIR/集成指南.md" .

        print_success "文件复制完成!"
        echo ""
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo -e "${YELLOW}接下来你需要:${NC}"
        echo ""
        echo "1. 修改你的 main.go 添加路由 (参考 集成指南.md)"
        echo ""
        echo "2. 安装依赖:"
        echo "   go get github.com/mattn/go-sqlite3"
        echo "   go mod tidy"
        echo ""
        echo "3. 导入许可证模块:"
        echo "   import \"yourproject/license/database\""
        echo "   import \"yourproject/license/handlers\""
        echo ""
        echo "详细步骤请查看: 集成指南.md"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ;;

    3)
        print_info "仅复制核心文件..."

        mkdir -p "$TARGET_DIR/license"
        cd "$TARGET_DIR"

        # 复制核心代码
        print_info "复制核心代码..."
        cp -r "$SOURCE_DIR/server/handlers" license/
        cp -r "$SOURCE_DIR/server/models" license/
        cp -r "$SOURCE_DIR/server/database" license/
        cp -r "$SOURCE_DIR/server/utils" license/

        print_success "核心文件复制完成!"
        echo ""
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo -e "${YELLOW}文件位置:${NC}"
        echo "  $TARGET_DIR/license/"
        echo ""
        echo -e "${YELLOW}你需要自己实现:${NC}"
        echo "  - HTTP 服务器"
        echo "  - 路由注册"
        echo "  - 前端界面 (可选)"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ;;

    *)
        print_error "无效选项"
        exit 1
        ;;
esac

echo ""
print_success "集成完成! 🎉"
echo ""
