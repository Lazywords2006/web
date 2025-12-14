#!/bin/bash
# 启动器 - 命令行版本

cd "$(dirname "$0")"

echo ""
echo "🚀 许可证验证系统启动器"
echo ""

# 检查 Python 3.11
if command -v python3.11 &> /dev/null; then
    PYTHON="python3.11"
elif command -v python3 &> /dev/null; then
    PYTHON="python3"
else
    echo "❌ 未找到 Python，请先安装 Python"
    exit 1
fi

echo "✓ 使用 Python: $PYTHON ($($PYTHON --version))"
echo ""

# 检查依赖
$PYTHON -c "import requests" 2>/dev/null
if [ $? -ne 0 ]; then
    echo "📦 正在安装依赖 requests..."
    $PYTHON -m pip install requests
    echo ""
fi

# 运行启动器
$PYTHON launcher_cli.py

echo ""
echo "按回车键退出..."
read
