#!/bin/bash
# 许可证验证启动器 - macOS 兼容版本

cd "$(dirname "$0")"

echo "正在启动许可证验证系统..."
echo ""

# 使用系统 Python (避免 Xcode Command Line Tools 的兼容性问题)
/usr/bin/python3 launcher_wrapper.py

# 保持窗口打开以显示错误信息
if [ $? -ne 0 ]; then
    echo ""
    echo "按任意键退出..."
    read -n 1
fi
