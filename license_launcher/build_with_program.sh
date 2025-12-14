#!/bin/bash
# 增强版打包脚本 - 将启动器和程序打包到一起

clear
echo "========================================"
echo "  许可证启动器 - 程序打包工具"
echo "========================================"
echo ""

# 检测 Python
if command -v python3.11 &> /dev/null; then
    PYTHON="python3.11"
elif command -v python3 &> /dev/null; then
    PYTHON="python3"
else
    echo "❌ 未找到 Python"
    exit 1
fi

echo "✓ Python: $PYTHON ($($PYTHON --version))"
echo ""

# 检查依赖
echo "📦 检查依赖..."
$PYTHON -c "import requests" 2>/dev/null || $PYTHON -m pip install requests
$PYTHON -c "import PyInstaller" 2>/dev/null || $PYTHON -m pip install pyinstaller
echo "✓ 依赖完成"
echo ""

# 列出可打包的程序
echo "📁 扫描可打包的程序..."
programs=()
for file in *.exe *.app; do
    if [ -f "$file" ]; then
        programs+=("$file")
    fi
done

if [ ${#programs[@]} -eq 0 ]; then
    echo "⚠️  未找到可执行程序 (.exe 或 .app)"
    echo ""
    echo "请将您的程序复制到当前目录，然后重新运行此脚本"
    echo ""
    read -p "按回车键退出..."
    exit 1
fi

echo "找到以下程序:"
for i in "${!programs[@]}"; do
    echo "  $((i+1))) ${programs[$i]}"
done
echo "  0) 不包含程序（仅启动器）"
echo ""

# 选择程序
while true; do
    read -p "请选择要打包的程序 [0-${#programs[@]}]: " choice
    if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 0 ] && [ "$choice" -le "${#programs[@]}" ]; then
        break
    fi
    echo "无效选择，请重试"
done

if [ "$choice" -eq 0 ]; then
    TARGET_EXE=""
    PACK_MODE="仅启动器"
else
    TARGET_EXE="${programs[$((choice-1))]}"
    PACK_MODE="启动器 + $TARGET_EXE"
fi

echo ""
echo "=========================================="
echo "打包配置:"
echo "  模式: $PACK_MODE"
if [ -n "$TARGET_EXE" ]; then
    echo "  程序: $TARGET_EXE"
    echo "  大小: $(ls -lh "$TARGET_EXE" | awk '{print $5}')"
fi
echo "=========================================="
echo ""

# 确认
read -p "确认打包? [y/N]: " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    echo "已取消"
    exit 0
fi

echo ""
echo "🔧 准备打包配置..."

# 更新配置文件
if [ -n "$TARGET_EXE" ]; then
    cat > launcher_config.json << EOF
{
  "server_url": "http://localhost:8080",
  "target_exe": "$TARGET_EXE",
  "license_file": ".license",
  "use_gui": "auto"
}
EOF
else
    cat > launcher_config.json << EOF
{
  "server_url": "http://localhost:8080",
  "target_exe": "",
  "license_file": "license.dat",
  "use_gui": "auto"
}
EOF
fi

echo "✓ 配置已更新"
echo ""

# 构建 PyInstaller 命令
echo "📦 开始打包..."
echo ""

CMD="$PYTHON -m PyInstaller"
CMD="$CMD --onefile"
CMD="$CMD --name=许可证验证"
CMD="$CMD --add-data launcher_config.json:."

# 如果包含程序，添加到打包
if [ -n "$TARGET_EXE" ]; then
    CMD="$CMD --add-data $TARGET_EXE:."
fi

CMD="$CMD --hidden-import=requests"
CMD="$CMD --hidden-import=tkinter"
CMD="$CMD --clean"
CMD="$CMD launcher.py"

# 执行打包
eval $CMD

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================"
    echo "✅ 打包成功！"
    echo "========================================"
    echo ""

    OUTPUT_FILE="dist/许可证验证"
    FILE_SIZE=$(ls -lh "$OUTPUT_FILE" 2>/dev/null | awk '{print $5}')

    echo "生成的文件:"
    echo "  位置: $OUTPUT_FILE"
    echo "  大小: $FILE_SIZE"
    echo ""

    if [ -n "$TARGET_EXE" ]; then
        echo "✨ 特别说明:"
        echo "  - 您的程序已内嵌到启动器中"
        echo "  - 用户只需要这一个文件"
        echo "  - 首次运行会自动释放程序"
        echo ""
    fi

    echo "分发清单:"
    echo "  必需: $OUTPUT_FILE"
    echo "  可选: 使用说明.txt"
    echo ""

    echo "测试运行:"
    echo "  $OUTPUT_FILE"
    echo ""

    # 询问是否测试
    read -p "现在测试运行? [y/N]: " test_run
    if [[ "$test_run" =~ ^[Yy]$ ]]; then
        echo ""
        echo "正在启动..."
        "$OUTPUT_FILE"
    fi

else
    echo ""
    echo "❌ 打包失败"
    echo ""
    echo "常见问题:"
    echo "  1. 检查是否安装了所有依赖"
    echo "  2. 确保程序文件存在且可读"
    echo "  3. 查看上方的错误信息"
    exit 1
fi
