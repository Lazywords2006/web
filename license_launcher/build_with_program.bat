@echo off
chcp 65001 > nul
cls

echo ========================================
echo   许可证启动器 - 程序打包工具
echo ========================================
echo.

REM 检查 Python
where python >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ 未找到 Python
    pause
    exit /b 1
)

echo ✓ Python 已安装
python --version
echo.

REM 检查依赖
echo 📦 检查依赖...
python -c "import requests" >nul 2>&1
if %ERRORLEVEL% NEQ 0 python -m pip install requests
python -c "import PyInstaller" >nul 2>&1
if %ERRORLEVEL% NEQ 0 python -m pip install pyinstaller
echo ✓ 依赖完成
echo.

REM 扫描可执行程序
echo 📁 扫描可打包的程序...
set count=0
for %%f in (*.exe) do (
    if not "%%f"=="build_with_program.exe" (
        set /a count+=1
        set "prog!count!=%%f"
        echo   !count!^) %%f
    )
)

if %count%==0 (
    echo ⚠️  未找到可执行程序 ^(.exe^)
    echo.
    echo 请将您的程序复制到当前目录，然后重新运行此脚本
    echo.
    pause
    exit /b 1
)

echo   0^) 不包含程序（仅启动器）
echo.

REM 选择程序
:choose
set /p choice="请选择要打包的程序 [0-%count%]: "
if "%choice%"=="0" (
    set TARGET_EXE=
    set PACK_MODE=仅启动器
    goto confirm
)

if %choice% GTR %count% goto choose
if %choice% LSS 1 goto choose

setlocal enabledelayedexpansion
set TARGET_EXE=!prog%choice%!
set PACK_MODE=启动器 + !TARGET_EXE!
endlocal & set TARGET_EXE=%TARGET_EXE%& set PACK_MODE=%PACK_MODE%

:confirm
echo.
echo ==========================================
echo 打包配置:
echo   模式: %PACK_MODE%
if defined TARGET_EXE (
    echo   程序: %TARGET_EXE%
    for %%f in (%TARGET_EXE%) do echo   大小: %%~zf 字节
)
echo ==========================================
echo.

set /p confirm="确认打包? [y/N]: "
if /i not "%confirm%"=="y" (
    echo 已取消
    pause
    exit /b 0
)

echo.
echo 🔧 准备打包配置...

REM 更新配置文件
if defined TARGET_EXE (
    (
        echo {
        echo   "server_url": "http://localhost:8080",
        echo   "target_exe": "%TARGET_EXE%",
        echo   "license_file": ".license",
        echo   "use_gui": "auto"
        echo }
    ) > launcher_config.json
) else (
    (
        echo {
        echo   "server_url": "http://localhost:8080",
        echo   "target_exe": "",
        echo   "license_file": "license.dat",
        echo   "use_gui": "auto"
        echo }
    ) > launcher_config.json
)

echo ✓ 配置已更新
echo.

REM 构建打包命令
echo 📦 开始打包...
echo.

set CMD=python -m PyInstaller
set CMD=%CMD% --onefile
set CMD=%CMD% --name=许可证验证
set CMD=%CMD% --add-data "launcher_config.json;."

if defined TARGET_EXE (
    set CMD=%CMD% --add-data "%TARGET_EXE%;."
)

set CMD=%CMD% --hidden-import=requests
set CMD=%CMD% --hidden-import=tkinter
set CMD=%CMD% --clean
set CMD=%CMD% launcher.py

REM 执行打包
%CMD%

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo ✅ 打包成功！
    echo ========================================
    echo.

    set OUTPUT_FILE=dist\许可证验证.exe

    echo 生成的文件:
    echo   位置: %OUTPUT_FILE%
    for %%f in (%OUTPUT_FILE%) do echo   大小: %%~zf 字节
    echo.

    if defined TARGET_EXE (
        echo ✨ 特别说明:
        echo   - 您的程序已内嵌到启动器中
        echo   - 用户只需要这一个文件
        echo   - 首次运行会自动释放程序
        echo.
    )

    echo 分发清单:
    echo   必需: %OUTPUT_FILE%
    echo   可选: 使用说明.txt
    echo.

    echo 测试运行:
    echo   %OUTPUT_FILE%
    echo.

    set /p test_run="现在测试运行? [y/N]: "
    if /i "%test_run%"=="y" (
        echo.
        echo 正在启动...
        start "" "%OUTPUT_FILE%"
    )

) else (
    echo.
    echo ❌ 打包失败
    echo.
    echo 常见问题:
    echo   1. 检查是否安装了所有依赖
    echo   2. 确保程序文件存在且可读
    echo   3. 查看上方的错误信息
)

pause
