#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
程序启动器 - 许可证验证包装器
在启动主程序前验证许可证
"""

import sys
import os
import subprocess
import uuid
import platform
import threading
import time
import json

# ⚠️ 重要：macOS Tkinter 兼容性检查必须在导入 tkinter 之前
if platform.system() == 'Darwin':
    if '/Library/Developer/CommandLineTools' in sys.executable:
        print("\n" + "="*60)
        print("❌ 检测到您在使用 Xcode Command Line Tools 的 Python")
        print("="*60)
        print("\n该 Python 版本的 Tkinter 在 macOS 上存在兼容性问题。")
        print("\n解决方案:")
        print("  1. 使用系统 Python:")
        print("     /usr/bin/python3", ' '.join(sys.argv))
        print("\n  2. 或者安装独立的 Python:")
        print("     brew install python@3.11")
        print("     或从 https://www.python.org 下载安装")
        print("\n" + "="*60 + "\n")
        sys.exit(1)

# 通过检查后才导入 tkinter
import tkinter as tk
from tkinter import messagebox
import requests

class LicenseLauncher:
    def __init__(self):
        # 配置
        self.server_url = "http://localhost:8080"
        self.target_exe = "lzy_zte_12.10.exe"  # 要启动的目标程序
        self.license_file = "license.dat"  # 保存许可证信息的文件

        # 窗口设置
        self.root = tk.Tk()
        self.root.title("许可证验证")
        self.root.geometry("450x300")
        self.root.resizable(False, False)

        # 获取硬件ID
        self.hwid = self.get_hardware_id()

        # 创建界面
        self.create_widgets()

        # 检查保存的许可证
        self.check_saved_license()

    def get_hardware_id(self):
        """生成硬件ID"""
        try:
            # macOS
            if platform.system() == 'Darwin':
                import subprocess
                result = subprocess.run(['system_profiler', 'SPHardwareDataType'],
                                      capture_output=True, text=True)
                for line in result.stdout.split('\n'):
                    if 'Serial Number' in line or 'UUID' in line:
                        return line.split(':')[1].strip()

            # Windows
            elif platform.system() == 'Windows':
                import subprocess
                result = subprocess.run(['wmic', 'csproduct', 'get', 'uuid'],
                                      capture_output=True, text=True)
                lines = result.stdout.strip().split('\n')
                if len(lines) >= 2:
                    return lines[1].strip()

            # Linux
            elif platform.system() == 'Linux':
                with open('/etc/machine-id', 'r') as f:
                    return f.read().strip()
        except:
            pass

        # 兜底方案
        return str(uuid.uuid5(uuid.NAMESPACE_DNS, platform.node()))

    def create_widgets(self):
        """创建界面组件"""
        # 标题
        title_label = tk.Label(
            self.root,
            text="🔐 许可证验证系统",
            font=("Arial", 16, "bold")
        )
        title_label.pack(pady=20)

        # 许可证输入框
        input_frame = tk.Frame(self.root)
        input_frame.pack(pady=10, padx=30, fill='x')

        tk.Label(input_frame, text="许可证密钥:", font=("Arial", 11)).pack(anchor='w')

        self.license_entry = tk.Entry(input_frame, font=("Arial", 11), width=40)
        self.license_entry.pack(pady=5, fill='x')
        self.license_entry.bind('<Return>', lambda e: self.activate_license())

        # 激活按钮
        self.activate_btn = tk.Button(
            self.root,
            text="激活并启动程序",
            font=("Arial", 12, "bold"),
            bg="#4CAF50",
            fg="white",
            command=self.activate_license,
            cursor="hand2",
            height=2
        )
        self.activate_btn.pack(pady=15, padx=30, fill='x')

        # 状态标签
        self.status_label = tk.Label(
            self.root,
            text="",
            font=("Arial", 10),
            fg="gray"
        )
        self.status_label.pack(pady=5)

        # 硬件ID显示
        hwid_label = tk.Label(
            self.root,
            text=f"设备ID: {self.hwid[:32]}...",
            font=("Arial", 8),
            fg="gray"
        )
        hwid_label.pack(side='bottom', pady=10)

    def check_saved_license(self):
        """检查保存的许可证"""
        if os.path.exists(self.license_file):
            try:
                with open(self.license_file, 'r') as f:
                    data = json.load(f)
                    license_key = data.get('key')
                    hwid = data.get('hwid')

                    # 验证硬件ID是否匹配
                    if hwid == self.hwid:
                        self.status_label.config(text="正在验证已保存的许可证...", fg="blue")
                        self.root.update()

                        # 验证许可证是否仍然有效
                        if self.verify_license(license_key):
                            self.launch_program()
                            return
            except:
                pass

        self.status_label.config(text="请输入许可证密钥", fg="gray")

    def verify_license(self, license_key):
        """验证许可证"""
        try:
            response = requests.post(
                f"{self.server_url}/api/activate",
                json={"key": license_key, "hwid": self.hwid},
                timeout=5
            )

            if response.status_code == 200:
                return True

        except Exception as e:
            print(f"验证失败: {e}")

        return False

    def activate_license(self):
        """激活许可证"""
        license_key = self.license_entry.get().strip()

        if not license_key:
            messagebox.showwarning("提示", "请输入许可证密钥")
            return

        self.activate_btn.config(state='disabled')
        self.status_label.config(text="正在激活许可证...", fg="blue")
        self.root.update()

        # 在后台线程中执行激活
        threading.Thread(target=self._do_activate, args=(license_key,), daemon=True).start()

    def _do_activate(self, license_key):
        """执行激活（后台线程）"""
        try:
            response = requests.post(
                f"{self.server_url}/api/activate",
                json={"key": license_key, "hwid": self.hwid},
                timeout=10
            )

            if response.status_code == 200:
                data = response.json()

                # 保存许可证信息
                with open(self.license_file, 'w') as f:
                    json.dump({
                        'key': license_key,
                        'hwid': self.hwid
                    }, f)

                self.root.after(0, lambda: self.status_label.config(
                    text=f"✅ 激活成功！正在启动程序...",
                    fg="green"
                ))

                time.sleep(1)
                self.root.after(0, self.launch_program)

            else:
                error_msg = response.json().get('error', '未知错误')
                self.root.after(0, lambda: messagebox.showerror("激活失败", error_msg))
                self.root.after(0, lambda: self.activate_btn.config(state='normal'))
                self.root.after(0, lambda: self.status_label.config(text="", fg="gray"))

        except requests.exceptions.ConnectionError:
            self.root.after(0, lambda: messagebox.showerror(
                "连接错误",
                "无法连接到许可证服务器\n请检查服务器是否运行"
            ))
            self.root.after(0, lambda: self.activate_btn.config(state='normal'))
            self.root.after(0, lambda: self.status_label.config(text="", fg="gray"))

        except Exception as e:
            self.root.after(0, lambda: messagebox.showerror("错误", f"激活失败: {str(e)}"))
            self.root.after(0, lambda: self.activate_btn.config(state='normal'))
            self.root.after(0, lambda: self.status_label.config(text="", fg="gray"))

    def launch_program(self):
        """启动目标程序"""
        if not os.path.exists(self.target_exe):
            messagebox.showerror("错误", f"找不到程序文件: {self.target_exe}")
            self.root.quit()
            return

        try:
            # 隐藏验证窗口
            self.root.withdraw()

            # 启动目标程序
            if platform.system() == 'Windows':
                subprocess.Popen([self.target_exe])
            else:
                # 如果是在 macOS/Linux 上用 Wine 运行
                subprocess.Popen(['wine', self.target_exe])

            # 等待一下确保程序启动
            time.sleep(2)

            # 关闭启动器
            self.root.quit()

        except Exception as e:
            messagebox.showerror("启动失败", f"无法启动程序: {str(e)}")
            self.root.deiconify()

    def run(self):
        """运行启动器"""
        self.root.mainloop()

if __name__ == "__main__":
    # 兼容性检查已在文件开头完成
    app = LicenseLauncher()
    app.run()
