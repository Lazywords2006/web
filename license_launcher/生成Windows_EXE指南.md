# 如何生成 Windows .exe 文件

## ⚠️ 重要说明

**PyInstaller 不支持跨平台编译**,这意味着:
- 在 macOS 上只能生成 macOS 可执行文件
- 在 Windows 上只能生成 Windows .exe 文件
- 在 Linux 上只能生成 Linux 可执行文件

要生成 Windows .exe 文件,您必须在 Windows 系统上运行打包命令。

---

## ✅ 方案 1: 在 Windows 系统上打包(推荐)

### 步骤 1: 准备 Windows 环境

您需要一台 Windows 电脑或虚拟机,安装:
- **Python 3.11** (从 python.org 下载)
- **Git** (可选,用于传输文件)

### 步骤 2: 传输文件

将整个 `license_launcher` 文件夹复制到 Windows:

**选项 A: 使用 U 盘**
```
直接复制整个 license_launcher 文件夹
```

**选项 B: 使用 Git**
```bash
# 在 Windows 上克隆仓库
git clone <您的仓库地址>
cd 网络验证/license_launcher
```

**选项 C: 使用网络共享/OneDrive/Dropbox**
```
通过云盘同步文件夹
```

### 步骤 3: 在 Windows 上运行打包脚本

1. 打开 **命令提示符** (CMD) 或 **PowerShell**
2. 进入 license_launcher 目录:
   ```cmd
   cd path\to\license_launcher
   ```

3. 运行打包脚本:
   ```cmd
   build_with_program.bat
   ```

4. 按提示选择 `lzy_zte_12.10.exe`

### 步骤 4: 获取生成的 .exe 文件

打包完成后,在 `dist\` 目录找到:
```
dist\许可证验证.exe
```

这就是 Windows 可执行文件!

---

## 💻 方案 2: 使用云端 Windows 虚拟机

### Azure Windows VM (推荐)

1. 创建 Windows Server VM
2. 通过远程桌面连接
3. 安装 Python 和依赖
4. 上传文件并打包

### AWS EC2 Windows

1. 启动 Windows Server 实例
2. 使用 RDP 连接
3. 安装 Python 和依赖
4. 上传文件并打包

---

## 🐳 方案 3: 使用 Docker (实验性)

**注意**: 这个方案较复杂,仅适合有 Docker 经验的用户。

### 前提条件

- 安装 Docker Desktop (支持 Windows 容器)
- 切换到 Windows 容器模式

### 创建 Dockerfile

在 `license_launcher` 目录创建 `Dockerfile.windows`:

```dockerfile
# escape=`
FROM python:3.11-windowsservercore

WORKDIR /app

# 复制文件
COPY requirements.txt .
COPY launcher.py .
COPY launcher_config.json .
COPY lzy_zte_12.10.exe .

# 安装依赖
RUN pip install --no-cache-dir -r requirements.txt
RUN pip install pyinstaller

# 打包
RUN pyinstaller --onefile `
    --name="lzy_zte_许可证验证" `
    --add-data "launcher_config.json;." `
    --add-data "lzy_zte_12.10.exe;." `
    --hidden-import=requests `
    --hidden-import=tkinter `
    --clean `
    launcher.py

CMD ["cmd"]
```

### 构建和提取

```bash
# 构建镜像
docker build -f Dockerfile.windows -t license-builder-windows .

# 运行容器
docker run -d --name builder license-builder-windows

# 提取生成的 exe
docker cp builder:/app/dist/lzy_zte_许可证验证.exe ./

# 清理
docker rm -f builder
```

---

## 🔧 方案 4: 手动在 Windows 上打包

如果自动化脚本不工作,可以手动运行命令。

### 在 Windows 命令提示符中运行:

```cmd
# 1. 进入目录
cd path\to\license_launcher

# 2. 安装依赖
python -m pip install pyinstaller requests

# 3. 确认文件存在
dir lzy_zte_12.10.exe
dir launcher.py
dir launcher_config.json

# 4. 运行打包命令
python -m PyInstaller ^
    --onefile ^
    --name="lzy_zte_许可证验证" ^
    --add-data "launcher_config.json;." ^
    --add-data "lzy_zte_12.10.exe;." ^
    --hidden-import=requests ^
    --hidden-import=tkinter ^
    --clean ^
    launcher.py

# 5. 检查结果
dir dist\许可证验证.exe
```

---

## 📋 当前可用的文件

在您的 `license_launcher` 目录中:

```
license_launcher/
├── launcher.py                    # 启动器源代码
├── launcher_config.json           # 配置文件
├── lzy_zte_12.10.exe             # 要打包的程序
├── build_with_program.bat        # Windows 自动打包脚本 ⭐
├── build_with_program.sh         # macOS/Linux 自动打包脚本
└── requirements.txt               # Python 依赖
```

**要生成 Windows .exe,请使用 `build_with_program.bat` 在 Windows 系统上运行。**

---

## ✨ 快速检查清单

在 Windows 系统上打包前,确认:

- [ ] 已安装 Python 3.7+
- [ ] 已进入 license_launcher 目录
- [ ] lzy_zte_12.10.exe 文件存在
- [ ] launcher.py 文件存在
- [ ] launcher_config.json 文件存在
- [ ] 已安装 pip

然后运行:
```cmd
build_with_program.bat
```

---

## 🎯 推荐流程总结

**最简单的方法:**

1. 找一台 Windows 电脑(实体机/虚拟机/云服务器)
2. 复制整个 `license_launcher` 文件夹到 Windows
3. 双击运行 `build_with_program.bat`
4. 选择 `lzy_zte_12.10.exe`
5. 等待打包完成
6. 在 `dist\` 目录获取生成的 .exe 文件

**就这么简单!** 🎉

---

## 💡 常见问题

### Q: 我能在 macOS 上生成 Windows .exe 吗?

**A:** 不能。PyInstaller 不支持跨平台编译。必须在 Windows 上生成 Windows 可执行文件。

### Q: 我没有 Windows 电脑怎么办?

**A:** 可以使用:
1. 云端 Windows VM (Azure/AWS/阿里云)
2. 虚拟机软件 (Parallels/VMware/VirtualBox)
3. 让有 Windows 电脑的朋友帮忙打包

### Q: 虚拟机需要什么配置?

**A:** 最低配置:
- Windows 10/11 或 Windows Server
- 2GB RAM
- 10GB 磁盘空间
- 可以联网(下载 Python 和依赖)

### Q: 打包需要多长时间?

**A:** 通常 1-3 分钟,取决于:
- 电脑性能
- 程序大小
- 是否首次安装依赖

---

**创建时间**: 2025-12-14
**平台要求**: Windows (用于生成 .exe)
