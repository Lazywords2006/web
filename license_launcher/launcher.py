#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
许可证验证启动器 - 通用版本
支持 GUI 和命令行两种模式，自动选择最佳模式
"""

import sys
import os
import subprocess
import uuid
import platform
import time
import json

# 尝试导入 requests
try:
    import requests
except ImportError:
    print("\n❌ 缺少 requests 库")
    print("\n安装方法:")
    print("  python -m pip install requests")
    print("\n或者:")
    print("  python3 -m pip install requests")
    sys.exit(1)

# 全局配置
CONFIG = {
    'server_url': 'http://106.14.255.49:8080',
    'target_exe': '',  # 留空表示仅验证许可证，不启动程序
    'license_file': 'license.dat',
    'use_gui': 'auto',  # 'auto', 'force_gui', 'force_cli'
}

class LicenseManager:
    """许可证管理器核心类"""

    def __init__(self, config=None):
        self.config = config or CONFIG
        self.server_url = self.config['server_url']
        self.target_exe = self.config['target_exe']
        self.license_file = self.config['license_file']
        self.hwid = self.get_hardware_id()

    def get_hardware_id(self):
        """生成硬件ID"""
        try:
            if platform.system() == 'Darwin':
                result = subprocess.run(['system_profiler', 'SPHardwareDataType'],
                                      capture_output=True, text=True)
                for line in result.stdout.split('\n'):
                    if 'Serial Number' in line or 'UUID' in line:
                        return line.split(':')[1].strip()

            elif platform.system() == 'Windows':
                result = subprocess.run(['wmic', 'csproduct', 'get', 'uuid'],
                                      capture_output=True, text=True)
                lines = result.stdout.strip().split('\n')
                if len(lines) >= 2:
                    return lines[1].strip()

            elif platform.system() == 'Linux':
                with open('/etc/machine-id', 'r') as f:
                    return f.read().strip()
        except:
            pass

        return str(uuid.uuid5(uuid.NAMESPACE_DNS, platform.node()))

    def load_saved_license(self):
        """加载保存的许可证"""
        if os.path.exists(self.license_file):
            try:
                with open(self.license_file, 'r') as f:
                    data = json.load(f)
                    license_key = data.get('key')
                    hwid = data.get('hwid')

                    if hwid == self.hwid:
                        return license_key
            except:
                pass
        return None

    def save_license(self, license_key):
        """保存许可证"""
        try:
            with open(self.license_file, 'w') as f:
                json.dump({
                    'key': license_key,
                    'hwid': self.hwid
                }, f)
            return True
        except Exception as e:
            print(f"保存许可证失败: {e}")
            return False

    def verify_license(self, license_key):
        """验证许可证"""
        try:
            response = requests.post(
                f"{self.server_url}/api/activate",
                json={"key": license_key, "hwid": self.hwid},
                timeout=10
            )

            if response.status_code == 200:
                return True, response.json()
            else:
                error = response.json().get('error', '未知错误')
                return False, error

        except requests.exceptions.ConnectionError:
            return False, "无法连接到服务器"
        except Exception as e:
            return False, str(e)

    def extract_embedded_program(self):
        """释放内嵌的程序文件（如果有）"""
        import shutil

        if not self.target_exe:
            return None

        # 检查是否打包环境
        if not getattr(sys, 'frozen', False):
            # 不是打包环境，直接返回原路径
            return self.target_exe if os.path.exists(self.target_exe) else None

        # 打包环境，查找内嵌的程序
        if hasattr(sys, '_MEIPASS'):
            # PyInstaller 临时目录
            embedded_path = os.path.join(sys._MEIPASS, self.target_exe)

            if os.path.exists(embedded_path):
                # 释放到可执行文件所在目录
                if hasattr(sys, 'executable'):
                    # 获取可执行文件所在目录
                    exe_dir = os.path.dirname(os.path.abspath(sys.executable))
                else:
                    # 降级到当前工作目录
                    exe_dir = os.getcwd()

                extract_path = os.path.join(exe_dir, self.target_exe)

                # 如果已存在且是最新的，直接使用
                if os.path.exists(extract_path):
                    # 比较文件大小
                    if os.path.getsize(embedded_path) == os.path.getsize(extract_path):
                        return extract_path

                # 复制文件
                try:
                    shutil.copy2(embedded_path, extract_path)

                    # macOS/Linux 需要添加执行权限
                    if platform.system() != 'Windows':
                        os.chmod(extract_path, 0o755)

                    return extract_path
                except Exception as e:
                    print(f"释放文件失败: {e}")
                    return None

        # 未找到内嵌程序，尝试当前目录
        if os.path.exists(self.target_exe):
            return self.target_exe

        return None

    def launch_program(self, program_path=None):
        """启动目标程序"""
        program = program_path or self.target_exe

        if not program:
            return True, "无需启动程序"

        # 尝试释放内嵌程序
        actual_program = self.extract_embedded_program()

        if not actual_program:
            return False, f"找不到程序: {program}"

        try:
            if platform.system() == 'Windows':
                subprocess.Popen([actual_program])
            else:
                subprocess.Popen(['wine', actual_program])

            return True, "程序已启动"
        except Exception as e:
            return False, f"启动失败: {e}"


class CLIInterface:
    """命令行界面"""

    def __init__(self, manager):
        self.manager = manager

    def print_header(self):
        print("\n" + "="*60)
        print("  🔐 许可证验证系统")
        print("="*60 + "\n")

    def run(self):
        self.print_header()

        # 检查保存的许可证
        saved_key = self.manager.load_saved_license()

        if saved_key:
            print("📋 找到已保存的许可证，正在验证...")
            success, result = self.manager.verify_license(saved_key)

            if success:
                print("✅ 许可证验证通过！")
                self.show_license_info(result)
                self.launch_if_needed()
                return True
            else:
                print(f"❌ 验证失败: {result}")
                os.remove(self.manager.license_file)
                print("已删除无效许可证，请重新输入\n")

        # 需要输入许可证
        print("请输入许可证密钥:")
        license_key = input("密钥: ").strip()

        if not license_key:
            print("\n❌ 许可证密钥不能为空")
            return False

        print(f"\n正在激活许可证...")
        print(f"设备ID: {self.manager.hwid[:32]}...")

        success, result = self.manager.verify_license(license_key)

        if success:
            print("\n" + "="*60)
            print("✅ 许可证激活成功！")
            print("="*60)

            self.show_license_info(result)
            self.manager.save_license(license_key)
            self.launch_if_needed()
            return True
        else:
            print("\n" + "="*60)
            print(f"❌ 激活失败: {result}")
            print("="*60)
            return False

    def show_license_info(self, info):
        """显示许可证信息"""
        if 'expires_at' in info:
            print(f"📅 过期时间: {info['expires_at']}")
        if 'product_name' in info:
            print(f"📦 产品名称: {info['product_name']}")

    def launch_if_needed(self):
        """启动程序（如果配置了）"""
        if self.manager.target_exe:
            print("\n" + "-"*60)
            print(f"🚀 正在启动程序: {self.manager.target_exe}")
            success, msg = self.manager.launch_program()
            if success:
                print(f"✅ {msg}")
            else:
                print(f"❌ {msg}")
            print("-"*60 + "\n")


class GUIInterface:
    """图形界面"""

    def __init__(self, manager):
        self.manager = manager
        self.root = None
        self.license_entry = None
        self.status_label = None
        self.activate_btn = None

    def run(self):
        import tkinter as tk
        from tkinter import messagebox
        import threading

        self.root = tk.Tk()
        self.root.title("许可证验证系统")
        self.root.geometry("500x350")
        self.root.resizable(False, False)

        # 创建界面
        self.create_widgets()

        # 检查保存的许可证
        self.check_saved_license()

        self.root.mainloop()

    def create_widgets(self):
        import tkinter as tk

        # 标题
        title_label = tk.Label(
            self.root,
            text="🔐 许可证验证系统",
            font=("Arial", 18, "bold")
        )
        title_label.pack(pady=25)

        # 输入框
        input_frame = tk.Frame(self.root)
        input_frame.pack(pady=15, padx=40, fill='x')

        tk.Label(input_frame, text="许可证密钥:", font=("Arial", 12)).pack(anchor='w')

        self.license_entry = tk.Entry(input_frame, font=("Arial", 11), width=45)
        self.license_entry.pack(pady=8, fill='x')
        self.license_entry.bind('<Return>', lambda e: self.activate_license())

        # 激活按钮
        self.activate_btn = tk.Button(
            self.root,
            text="激活许可证",
            font=("Arial", 13, "bold"),
            bg="#4CAF50",
            fg="white",
            command=self.activate_license,
            cursor="hand2",
            height=2
        )
        self.activate_btn.pack(pady=20, padx=40, fill='x')

        # 状态标签
        self.status_label = tk.Label(
            self.root,
            text="",
            font=("Arial", 10),
            fg="gray"
        )
        self.status_label.pack(pady=10)

        # 设备ID
        hwid_label = tk.Label(
            self.root,
            text=f"设备ID: {self.manager.hwid[:40]}...",
            font=("Arial", 8),
            fg="gray"
        )
        hwid_label.pack(side='bottom', pady=15)

    def check_saved_license(self):
        """检查保存的许可证"""
        import tkinter as tk
        from tkinter import messagebox

        saved_key = self.manager.load_saved_license()

        if saved_key:
            self.status_label.config(text="正在验证已保存的许可证...", fg="blue")
            self.root.update()

            success, result = self.manager.verify_license(saved_key)

            if success:
                self.status_label.config(text="✅ 许可证验证通过", fg="green")
                messagebox.showinfo("验证成功", "许可证已验证通过！")

                if self.manager.target_exe:
                    self.launch_program()
                else:
                    self.root.after(2000, self.root.quit)
            else:
                os.remove(self.manager.license_file)
                self.status_label.config(text="请输入新的许可证密钥", fg="gray")

    def activate_license(self):
        """激活许可证"""
        import tkinter as tk
        from tkinter import messagebox
        import threading

        license_key = self.license_entry.get().strip()

        if not license_key:
            messagebox.showwarning("提示", "请输入许可证密钥")
            return

        self.activate_btn.config(state='disabled')
        self.status_label.config(text="正在激活许可证...", fg="blue")
        self.root.update()

        threading.Thread(target=self._do_activate, args=(license_key,), daemon=True).start()

    def _do_activate(self, license_key):
        """执行激活（后台线程）"""
        from tkinter import messagebox

        success, result = self.manager.verify_license(license_key)

        if success:
            self.manager.save_license(license_key)

            self.root.after(0, lambda: self.status_label.config(
                text="✅ 激活成功！",
                fg="green"
            ))

            self.root.after(0, lambda: messagebox.showinfo(
                "激活成功",
                "许可证已成功激活！"
            ))

            time.sleep(1)

            if self.manager.target_exe:
                self.root.after(0, self.launch_program)
            else:
                self.root.after(0, self.root.quit)
        else:
            self.root.after(0, lambda: messagebox.showerror("激活失败", result))
            self.root.after(0, lambda: self.activate_btn.config(state='normal'))
            self.root.after(0, lambda: self.status_label.config(text="", fg="gray"))

    def launch_program(self):
        """启动程序"""
        from tkinter import messagebox

        if not self.manager.target_exe:
            self.root.quit()
            return

        success, msg = self.manager.launch_program()

        if success:
            self.root.withdraw()
            time.sleep(2)
            self.root.quit()
        else:
            messagebox.showerror("启动失败", msg)


def check_gui_available():
    """检查 GUI 是否可用"""
    # macOS Tkinter 兼容性检查
    if platform.system() == 'Darwin':
        if '/Library/Developer/CommandLineTools' in sys.executable:
            return False, "Xcode Command Line Tools Python 不支持 Tkinter"

    try:
        import tkinter
        return True, "Tkinter 可用"
    except ImportError:
        return False, "Tkinter 未安装"


def main():
    """主函数"""

    # 从配置文件加载配置（如果存在）
    config_file = 'launcher_config.json'

    # 在打包环境中,配置文件在 _MEIPASS 目录
    if getattr(sys, 'frozen', False) and hasattr(sys, '_MEIPASS'):
        config_file = os.path.join(sys._MEIPASS, 'launcher_config.json')

    if os.path.exists(config_file):
        try:
            with open(config_file, 'r', encoding='utf-8') as f:
                user_config = json.load(f)
                CONFIG.update(user_config)
        except Exception as e:
            print(f"警告: 加载配置文件失败: {e}")

    # 创建许可证管理器
    manager = LicenseManager(CONFIG)

    # 决定使用哪种界面
    use_gui = CONFIG.get('use_gui', 'auto')

    if use_gui == 'force_cli':
        # 强制使用命令行
        interface = CLIInterface(manager)
        return interface.run()

    elif use_gui == 'force_gui':
        # 强制使用 GUI
        gui_available, _ = check_gui_available()
        if gui_available:
            interface = GUIInterface(manager)
            return interface.run()
        else:
            print("❌ GUI 不可用，请使用命令行模式")
            return False

    else:  # auto
        # 自动选择
        gui_available, reason = check_gui_available()

        if gui_available:
            try:
                interface = GUIInterface(manager)
                return interface.run()
            except Exception as e:
                print(f"GUI 启动失败，切换到命令行模式: {e}\n")
                interface = CLIInterface(manager)
                return interface.run()
        else:
            # 使用命令行
            interface = CLIInterface(manager)
            return interface.run()


if __name__ == "__main__":
    try:
        success = main()
        sys.exit(0 if success else 1)
    except KeyboardInterrupt:
        print("\n\n已取消")
        sys.exit(0)
    except Exception as e:
        print(f"\n❌ 发生错误: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
