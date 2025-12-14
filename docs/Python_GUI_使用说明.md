# Python GUI 程序 - 使用说明

## 📦 程序说明

这是一个带图形界面的许可证验证程序示例,包含以下功能:
- ✅ 输入许可证密钥
- ✅ 激活按钮
- ✅ 测试连接按钮
- ✅ 实时状态显示
- ✅ 自动心跳监控

## 🚀 快速开始

### 1. 安装依赖

```bash
# 安装 requests 库
pip install requests

# 或者使用 requirements.txt
pip install -r requirements.txt
```

### 2. 运行程序

```bash
python python_gui_example.py
```

### 3. 使用流程

1. **点击"测试连接"** - 确认服务器正常运行
2. **输入许可证密钥** - 例如: `LICENSE-2025-XXX`
3. **点击"激活许可证"** - 开始激活
4. **查看状态日志** - 观察激活进度
5. **激活成功后** - 程序自动启动心跳监控

## 📋 界面说明

```
┌──────────────────────────────────────┐
│     🔐 许可证验证系统                │
├──────────────────────────────────────┤
│  硬件信息                            │
│  硬件ID: abc123...                   │
├──────────────────────────────────────┤
│  许可证密钥                          │
│  [LICENSE-2025-_____________]        │
├──────────────────────────────────────┤
│  [激活许可证]  [测试连接]           │
├──────────────────────────────────────┤
│  状态信息                            │
│  [10:30:15] 系统已启动               │
│  [10:30:20] ✅ 服务器连接成功!       │
│  [10:30:25] 正在激活许可证...        │
│  [10:30:26] ✅ 许可证激活成功!       │
│  [10:30:26] 💓 心跳监控已启动        │
├──────────────────────────────────────┤
│  服务器地址: http://localhost:8080   │
└──────────────────────────────────────┘
```

## 🛠️ 自定义配置

### 修改服务器地址

编辑 `python_gui_example.py`,找到第 20 行:

```python
self.server_url = "http://localhost:8080"
```

改为你的服务器地址:

```python
self.server_url = "https://your-server.com"
```

### 修改心跳间隔

找到第 150 行:

```python
time.sleep(30)  # 30秒心跳间隔
```

改为你需要的间隔:

```python
time.sleep(60)  # 60秒心跳间隔
```

### 修改窗口大小

找到第 17 行:

```python
self.root.geometry("500x400")
```

改为你需要的大小:

```python
self.root.geometry("600x500")
```

## 📦 打包成 exe

### 方法 1: 使用 PyInstaller

```bash
# 安装 PyInstaller
pip install pyinstaller

# 打包成单个 exe 文件
pyinstaller --onefile --noconsole --name="许可证验证程序" python_gui_example.py

# 打包后的 exe 在 dist/ 目录下
```

### 方法 2: 添加图标

```bash
# 准备一个 .ico 图标文件
# 然后打包
pyinstaller --onefile --noconsole --icon=app.ico --name="许可证验证程序" python_gui_example.py
```

### 方法 3: 使用配置文件

创建 `license_app.spec` 文件:

```python
# -*- mode: python ; coding: utf-8 -*-

block_cipher = None

a = Analysis(
    ['python_gui_example.py'],
    pathex=[],
    binaries=[],
    datas=[],
    hiddenimports=['requests'],
    hookspath=[],
    hooksconfig={},
    runtime_hooks=[],
    excludes=[],
    win_no_prefer_redirects=False,
    win_private_assemblies=False,
    cipher=block_cipher,
    noarchive=False,
)

pyz = PYZ(a.pure, a.zipped_data, cipher=block_cipher)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.zipfiles,
    a.datas,
    [],
    name='许可证验证程序',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    upx_exclude=[],
    runtime_tmpdir=None,
    console=False,
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
    icon='app.ico',  # 图标文件
)
```

使用配置文件打包:

```bash
pyinstaller license_app.spec
```

## 🎨 界面增强 (可选)

### 添加主题

```bash
# 安装 ttkbootstrap (现代化主题)
pip install ttkbootstrap
```

修改代码:

```python
import ttkbootstrap as ttk
from ttkbootstrap.constants import *

# 修改第 14 行
def __init__(self, root):
    self.root = root
    self.root.title("许可证验证系统")
    self.root.geometry("500x400")

    # 设置主题
    style = ttk.Style("flatly")  # 可选: flatly, darkly, cosmo 等
```

### 添加进度条

在激活过程中显示进度条:

```python
# 添加进度条
self.progress = ttk.Progressbar(main_frame, mode='indeterminate')
self.progress.grid(row=6, column=0, columnspan=2, sticky=(tk.W, tk.E), pady=10)
self.progress.grid_remove()  # 默认隐藏

# 激活时显示
def activate_license(self):
    self.progress.grid()
    self.progress.start()
    # ... 激活逻辑

# 激活完成后隐藏
def _do_activate(self):
    # ... 激活逻辑
    self.root.after(0, lambda: self.progress.stop())
    self.root.after(0, lambda: self.progress.grid_remove())
```

## 🔧 高级功能

### 1. 添加配置文件保存

```python
import json
import os

def save_config(self):
    """保存配置"""
    config = {
        "server_url": self.server_url,
        "license_key": self.license_key,
        "last_activated": time.strftime("%Y-%m-%d %H:%M:%S")
    }
    with open("config.json", "w") as f:
        json.dump(config, f, indent=2)

def load_config(self):
    """加载配置"""
    if os.path.exists("config.json"):
        with open("config.json", "r") as f:
            config = json.load(f)
            self.server_url = config.get("server_url", "http://localhost:8080")
            saved_key = config.get("license_key", "")
            if saved_key:
                self.license_entry.insert(0, saved_key)
```

### 2. 添加托盘图标

```bash
pip install pystray Pillow
```

```python
from pystray import Icon, Menu, MenuItem
from PIL import Image

def create_tray_icon(self):
    """创建系统托盘图标"""
    # 创建一个简单的图标
    image = Image.new('RGB', (64, 64), color='blue')

    def on_quit(icon, item):
        icon.stop()
        self.root.quit()

    menu = Menu(
        MenuItem('显示', lambda: self.root.deiconify()),
        MenuItem('退出', on_quit)
    )

    icon = Icon("license_app", image, "许可证验证", menu)
    threading.Thread(target=icon.run, daemon=True).start()
```

### 3. 添加日志文件

```python
import logging

# 配置日志
logging.basicConfig(
    filename='app.log',
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def log_status(self, message):
    """记录状态日志"""
    timestamp = time.strftime("%H:%M:%S")
    self.status_text.config(state='normal')
    self.status_text.insert(tk.END, f"[{timestamp}] {message}\n")
    self.status_text.see(tk.END)
    self.status_text.config(state='disabled')

    # 同时写入日志文件
    logging.info(message)
```

## 🐛 故障排查

### 问题 1: 无法连接到服务器

**症状**: 点击"测试连接"显示连接失败

**解决方法**:
1. 确认许可证服务器正在运行:
   ```bash
   curl http://localhost:8080/api/admin/stats
   ```
2. 检查防火墙设置
3. 确认服务器地址正确

### 问题 2: 激活失败

**症状**: 输入密钥后激活失败

**可能原因**:
- 许可证密钥不存在
- 许可证已过期
- 许可证已在其他设备激活
- 许可证已被封禁

**解决方法**:
1. 在管理后台检查许可证状态
2. 确认许可证密钥拼写正确
3. 查看服务器日志了解详情

### 问题 3: 打包后无法运行

**症状**: exe 文件双击后闪退

**解决方法**:
1. 不使用 `--noconsole` 参数重新打包,查看错误信息:
   ```bash
   pyinstaller --onefile python_gui_example.py
   ```
2. 确认所有依赖都已包含
3. 在命令行运行 exe 查看错误:
   ```bash
   cd dist
   .\许可证验证程序.exe
   ```

### 问题 4: 界面显示异常

**症状**: 界面布局混乱或字体太小

**解决方法**:
1. 调整 DPI 设置:
   ```python
   from ctypes import windll
   windll.shcore.SetProcessDpiAwareness(1)  # 在 Windows 上设置 DPI
   ```
2. 调整字体大小
3. 修改窗口尺寸

## 📊 完整流程图

```
用户运行 exe
    ↓
显示 GUI 界面
    ↓
用户输入许可证密钥
    ↓
点击"激活许可证"
    ↓
发送 POST /api/activate
    ↓
服务器验证并激活
    ↓
返回 JWT token
    ↓
GUI 显示"激活成功"
    ↓
启动心跳监控线程
    ↓
每 30 秒发送 POST /api/heartbeat
    ↓
验证成功 → 继续运行
验证失败 → 3 次重试 → 退出程序
```

## 🎯 测试清单

- [ ] 测试连接功能正常
- [ ] 激活功能正常
- [ ] 心跳监控正常
- [ ] 激活失败时显示正确错误
- [ ] 网络断开时能正确重试
- [ ] 打包成 exe 后正常运行
- [ ] 界面在不同分辨率下正常显示
- [ ] 日志记录正常

## 📞 技术支持

遇到问题请检查:
1. **程序日志**: 查看状态信息框的日志
2. **服务器日志**: `server/server.log`
3. **网络连接**: 确认能访问服务器
4. **许可证状态**: 在管理后台查看

---

**祝你使用愉快! 🎉**

如有问题,请查看 [完整文档](./集成到EXE的完整指南.md)
