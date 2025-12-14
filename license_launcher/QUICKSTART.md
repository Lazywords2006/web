# 快速开始

## 🎯 三步完成

### 1. 配置

编辑 `launcher_config.json`:
```json
{
  "server_url": "http://localhost:8080",
  "target_exe": "your_program.exe",
  "use_gui": "auto"
}
```

### 2. 测试

```bash
# 安装依赖
pip install -r requirements.txt

# 运行
python launcher.py
```

### 3. 打包

**macOS/Linux:**
```bash
./build.sh
```

**Windows:**
```bash
build.bat
```

## 📦 打包后

生成的文件在 `dist/` 目录：
- `许可证验证` (macOS/Linux)
- `许可证验证.exe` (Windows)

## 💡 使用

```
dist/
└── 许可证验证.exe        # 分发这个文件
```

用户双击运行，输入许可证即可！

---

详细文档: [README.md](README.md)
