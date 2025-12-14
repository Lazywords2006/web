# 许可证验证启动器 - 使用文档

## 📦 简介

这是一个**通用的许可证验证启动器**，可以：
- ✅ 自动选择 GUI 或命令行界面
- ✅ 验证许可证并保存
- ✅ 可选择启动目标程序
- ✅ 支持打包成独立可执行文件
- ✅ 跨平台支持（Windows/macOS/Linux）

## 🚀 快速开始

### 方式 1：直接运行 Python 脚本

```bash
# 安装依赖
pip install -r requirements.txt

# 运行启动器
python launcher.py
```

### 方式 2：打包成可执行文件

**macOS/Linux:**
```bash
chmod +x build.sh
./build.sh
```

**Windows:**
```bash
build.bat
```

打包完成后，在 `dist/` 目录找到可执行文件。

## ⚙️ 配置说明

编辑 `launcher_config.json` 配置文件：

```json
{
  "server_url": "http://localhost:8080",
  "target_exe": "",
  "license_file": "license.dat",
  "use_gui": "auto"
}
```

### 配置项说明

| 配置项 | 说明 | 可选值 |
|--------|------|--------|
| `server_url` | 许可证服务器地址 | 任何有效的 HTTP(S) URL |
| `target_exe` | 启动的目标程序路径 | 文件路径或留空 |
| `license_file` | 许可证保存文件名 | 任意文件名 |
| `use_gui` | 界面模式 | `auto`, `force_gui`, `force_cli` |

### 界面模式说明

- **auto** (推荐): 自动检测环境，优先使用 GUI，不可用时降级到命令行
- **force_gui**: 强制使用 GUI 界面，如果不可用则报错
- **force_cli**: 强制使用命令行界面

## 📝 使用场景

### 场景 1：仅验证许可证（不启动程序）

配置文件:
```json
{
  "server_url": "http://your-server.com",
  "target_exe": "",
  "use_gui": "auto"
}
```

运行后仅验证许可证，不启动任何程序。

### 场景 2：验证后启动程序

配置文件:
```json
{
  "server_url": "http://your-server.com",
  "target_exe": "your_program.exe",
  "use_gui": "auto"
}
```

验证成功后自动启动 `your_program.exe`。

### 场景 3：作为程序包装器

将启动器和目标程序放在同一目录：
```
MyApp/
├── 许可证验证.exe        # 启动器
├── launcher_config.json  # 配置文件
└── MyProgram.exe         # 目标程序
```

配置 `target_exe` 为 `MyProgram.exe`，用户双击启动器即可。

## 🎨 界面预览

### GUI 模式
- 图形化界面
- 输入许可证密钥
- 显示验证状态
- 自动保存许可证

### 命令行模式
```
============================================================
  🔐 许可证验证系统
============================================================

请输入许可证密钥:
密钥: LICENSE-XXXX-XXXX

正在激活许可证...
设备ID: xxx...

============================================================
✅ 许可证激活成功！
============================================================
📅 过期时间: 2026-12-14
📦 产品名称: 标准版
```

## 🔧 高级用法

### 自定义服务器地址

打包时内置服务器地址：

1. 编辑 `launcher.py` 第 21 行：
   ```python
   CONFIG = {
       'server_url': 'https://your-domain.com',
       ...
   }
   ```

2. 重新打包即可。

### 禁用配置文件

如果不希望用户修改配置，可以删除 `launcher_config.json`，程序将使用代码中的默认配置。

### 添加图标

打包时添加自定义图标：

```bash
# macOS/Linux
pyinstaller --onefile --icon=icon.icns ...

# Windows
pyinstaller --onefile --icon=icon.ico ...
```

### 打包时包含目标程序

```bash
pyinstaller \
    --onefile \
    --add-data "YourProgram.exe:." \
    --add-data "launcher_config.json:." \
    launcher.py
```

## 🛠️ 开发说明

### 文件结构

```
license_launcher/
├── launcher.py              # 主程序（核心代码）
├── launcher_config.json     # 配置文件示例
├── build.sh                 # macOS/Linux 打包脚本
├── build.bat                # Windows 打包脚本
├── requirements.txt         # Python 依赖
└── README.md               # 本文档
```

### 核心类说明

**LicenseManager**
- 核心许可证管理逻辑
- 硬件ID生成
- 许可证验证和保存
- 程序启动

**CLIInterface**
- 命令行界面实现
- 输入输出处理

**GUIInterface**
- 图形界面实现
- Tkinter 组件管理

### 扩展功能

可以基于 `LicenseManager` 类添加：
- 心跳监控
- 自动更新检查
- 使用统计
- 离线验证
- 多许可证管理

## 📦 打包选项详解

### 基本打包

```bash
pyinstaller --onefile launcher.py
```

### 完整打包（推荐）

```bash
pyinstaller \
    --onefile \
    --name="许可证验证" \
    --add-data "launcher_config.json:." \
    --hidden-import=requests \
    --hidden-import=tkinter \
    --clean \
    launcher.py
```

### 参数说明

- `--onefile`: 打包成单个可执行文件
- `--name`: 指定输出文件名
- `--add-data`: 包含额外文件（配置等）
- `--hidden-import`: 包含隐式导入的模块
- `--clean`: 清理临时文件
- `--noconsole`: 隐藏控制台窗口（仅 Windows GUI）

## 🐛 故障排查

### 问题 1：GUI 无法启动

**症状**: 提示 "GUI 不可用"

**解决方案**:
1. 检查 Python 是否支持 Tkinter
2. macOS: 不要使用 Xcode Command Line Tools 的 Python
3. 或者设置 `use_gui` 为 `force_cli`

### 问题 2：无法连接服务器

**症状**: "无法连接到服务器"

**检查**:
1. 服务器是否运行？
2. `server_url` 配置是否正确？
3. 网络是否通畅？
4. 防火墙是否阻止？

### 问题 3：打包后无法运行

**症状**: 双击没反应或报错

**解决方案**:
1. 在终端运行查看错误信息
2. 检查是否缺少依赖模块
3. 添加 `--hidden-import` 包含缺失模块
4. 检查 `--add-data` 路径是否正确

### 问题 4：许可证验证失败

**检查**:
1. 许可证密钥是否正确？
2. 许可证是否已被激活？
3. 硬件ID是否匹配？
4. 许可证是否过期？

## 💡 最佳实践

### 1. 生产环境配置

```json
{
  "server_url": "https://license.yourdomain.com",
  "target_exe": "YourApp.exe",
  "license_file": ".license",
  "use_gui": "auto"
}
```

- 使用 HTTPS
- 隐藏许可证文件（.开头）
- 自动选择界面模式

### 2. 分发建议

**给用户的文件包:**
```
YourApp_v1.0/
├── 许可证验证.exe         # 启动器
├── YourApp.exe            # 您的程序
└── 使用说明.txt           # 简单说明
```

**使用说明.txt:**
```
使用方法：
1. 双击运行"许可证验证.exe"
2. 输入您的许可证密钥
3. 完成！程序将自动启动

如遇问题，请联系技术支持
```

### 3. 安全建议

- ✅ 使用 HTTPS 连接服务器
- ✅ 不要在配置文件中存储敏感信息
- ✅ 定期更新依赖库
- ✅ 考虑代码混淆（PyArmor）
- ✅ 添加签名（Windows）

### 4. 用户体验优化

- 提供清晰的错误提示
- 自动保存许可证（避免重复输入）
- 支持复制粘贴许可证
- 添加重试机制
- 显示剩余有效期

## 🔗 相关资源

- [主项目文档](../README.md)
- [API 文档](../README.md#api)
- [服务器部署](../README.md#部署)

## 📞 技术支持

遇到问题？

1. 查看本文档的故障排查部分
2. 检查服务器日志
3. 查看启动器运行日志
4. 联系技术支持

---

**版本**: 1.0.0
**最后更新**: 2025-12-14
**作者**: Claude Code
