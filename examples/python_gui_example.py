"""
许可证验证 GUI 程序
功能: 输入许可证密钥进行验证和激活
"""

import sys
import platform

# macOS 上需要使用系统自带的 Python 来运行 Tkinter
if platform.system() == 'Darwin':
    import os
    # 检查是否使用系统 Python
    if '/Library/Developer/CommandLineTools' in sys.executable:
        print("⚠️  检测到您在使用 Xcode Command Line Tools 的 Python")
        print("在 macOS 上运行 Tkinter 程序,请使用系统自带的 Python:")
        print()
        print("解决方案:")
        print("1. 使用系统 Python: /usr/bin/python3 python_gui_example.py")
        print("2. 或者安装独立的 Python (from python.org)")
        print()
        sys.exit(1)

import tkinter as tk
from tkinter import ttk, messagebox
import requests
import hashlib
import uuid
import threading
import time

class LicenseApp:
    def __init__(self, root):
        self.root = root
        self.root.title("许可证验证系统")
        self.root.geometry("500x400")
        self.root.resizable(False, False)

        # 服务器配置
        self.server_url = "http://106.14.255.49:8080"
        self.token = None
        self.license_key = None
        self.hwid = self.get_hardware_id()
        self.heartbeat_running = False

        # 创建界面
        self.create_widgets()

    def get_hardware_id(self):
        """获取硬件ID"""
        mac = ':'.join(['{:02x}'.format((uuid.getnode() >> elements) & 0xff)
                       for elements in range(0, 48, 8)][::-1])
        return hashlib.sha256(mac.encode()).hexdigest()

    def create_widgets(self):
        """创建界面组件"""
        # 主框架
        main_frame = ttk.Frame(self.root, padding="20")
        main_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))

        # 标题
        title_label = ttk.Label(main_frame, text="🔐 许可证验证系统",
                               font=("Arial", 20, "bold"))
        title_label.grid(row=0, column=0, columnspan=2, pady=20)

        # 硬件ID显示
        hwid_frame = ttk.LabelFrame(main_frame, text="硬件信息", padding="10")
        hwid_frame.grid(row=1, column=0, columnspan=2, sticky=(tk.W, tk.E), pady=10)

        hwid_label = ttk.Label(hwid_frame, text=f"硬件ID: {self.hwid[:32]}...",
                              font=("Courier", 9))
        hwid_label.grid(row=0, column=0, sticky=tk.W)

        # 许可证输入框
        license_frame = ttk.LabelFrame(main_frame, text="许可证密钥", padding="10")
        license_frame.grid(row=2, column=0, columnspan=2, sticky=(tk.W, tk.E), pady=10)

        self.license_entry = ttk.Entry(license_frame, width=40, font=("Arial", 11))
        self.license_entry.grid(row=0, column=0, padx=5, pady=5)
        self.license_entry.insert(0, "LICENSE-2025-")

        # 按钮
        button_frame = ttk.Frame(main_frame)
        button_frame.grid(row=3, column=0, columnspan=2, pady=20)

        self.activate_btn = ttk.Button(button_frame, text="激活许可证",
                                       command=self.activate_license,
                                       width=20)
        self.activate_btn.grid(row=0, column=0, padx=5)

        self.test_btn = ttk.Button(button_frame, text="测试连接",
                                   command=self.test_connection,
                                   width=20)
        self.test_btn.grid(row=0, column=1, padx=5)

        # 状态显示
        status_frame = ttk.LabelFrame(main_frame, text="状态信息", padding="10")
        status_frame.grid(row=4, column=0, columnspan=2, sticky=(tk.W, tk.E), pady=10)

        self.status_text = tk.Text(status_frame, height=8, width=55,
                                   font=("Courier", 9), state='disabled')
        self.status_text.grid(row=0, column=0)

        # 滚动条
        scrollbar = ttk.Scrollbar(status_frame, orient="vertical",
                                 command=self.status_text.yview)
        scrollbar.grid(row=0, column=1, sticky=(tk.N, tk.S))
        self.status_text['yscrollcommand'] = scrollbar.set

        # 底部信息
        footer_label = ttk.Label(main_frame,
                                text="服务器地址: " + self.server_url,
                                font=("Arial", 8))
        footer_label.grid(row=5, column=0, columnspan=2, pady=10)

        # 初始日志
        self.log_status("系统已启动")
        self.log_status(f"硬件ID: {self.hwid[:32]}...")

    def log_status(self, message):
        """记录状态日志"""
        timestamp = time.strftime("%H:%M:%S")
        self.status_text.config(state='normal')
        self.status_text.insert(tk.END, f"[{timestamp}] {message}\n")
        self.status_text.see(tk.END)
        self.status_text.config(state='disabled')

    def test_connection(self):
        """测试服务器连接"""
        self.log_status("正在测试服务器连接...")
        try:
            response = requests.get(f"{self.server_url}/api/admin/stats", timeout=5)
            if response.status_code == 200:
                self.log_status("✅ 服务器连接成功!")
                messagebox.showinfo("成功", "服务器连接正常!")
            else:
                self.log_status(f"❌ 服务器返回错误: {response.status_code}")
                messagebox.showerror("错误", f"服务器返回错误: {response.status_code}")
        except requests.exceptions.ConnectionError:
            self.log_status("❌ 无法连接到服务器")
            messagebox.showerror("错误",
                               "无法连接到服务器!\n\n请确认:\n" +
                               "1. 许可证服务器正在运行\n" +
                               f"2. 服务器地址正确: {self.server_url}\n" +
                               "3. 网络连接正常")
        except Exception as e:
            self.log_status(f"❌ 测试失败: {str(e)}")
            messagebox.showerror("错误", f"测试失败: {str(e)}")

    def activate_license(self):
        """激活许可证"""
        license_key = self.license_entry.get().strip()

        if not license_key:
            messagebox.showwarning("警告", "请输入许可证密钥!")
            return

        self.license_key = license_key
        self.log_status(f"正在激活许可证: {license_key}")
        self.activate_btn.config(state='disabled', text="激活中...")

        # 在新线程中执行激活
        threading.Thread(target=self._do_activate, daemon=True).start()

    def _do_activate(self):
        """执行激活请求"""
        try:
            response = requests.post(
                f"{self.server_url}/api/activate",
                json={
                    "key": self.license_key,
                    "hwid": self.hwid
                },
                timeout=10
            )

            self.root.after(0, lambda: self.activate_btn.config(state='normal', text="激活许可证"))

            if response.status_code == 200:
                data = response.json()
                if data.get("status") == "success":
                    self.token = data.get("token")
                    self.root.after(0, lambda: self.log_status("✅ 许可证激活成功!"))
                    self.root.after(0, lambda: messagebox.showinfo("成功",
                        "许可证激活成功!\n\n系统已开始心跳监控\n应用程序现在可以正常使用"))

                    # 启动心跳监控
                    self.start_heartbeat()

                    # 启动主应用
                    self.root.after(0, self.start_main_app)
                else:
                    error_msg = data.get("error", "未知错误")
                    self.root.after(0, lambda: self.log_status(f"❌ 激活失败: {error_msg}"))
                    self.root.after(0, lambda: messagebox.showerror("失败",
                        f"许可证激活失败!\n\n错误信息: {error_msg}"))
            else:
                error_msg = f"HTTP {response.status_code}"
                try:
                    data = response.json()
                    error_msg = data.get("error", error_msg)
                except:
                    pass

                self.root.after(0, lambda: self.log_status(f"❌ 激活失败: {error_msg}"))
                self.root.after(0, lambda: messagebox.showerror("失败",
                    f"许可证激活失败!\n\n错误信息: {error_msg}\n\n可能原因:\n" +
                    "• 许可证密钥不存在\n" +
                    "• 许可证已过期\n" +
                    "• 许可证已在其他设备激活\n" +
                    "• 许可证已被封禁"))

        except requests.exceptions.ConnectionError:
            self.root.after(0, lambda: self.activate_btn.config(state='normal', text="激活许可证"))
            self.root.after(0, lambda: self.log_status("❌ 无法连接到服务器"))
            self.root.after(0, lambda: messagebox.showerror("错误",
                "无法连接到服务器!\n\n请检查网络连接和服务器状态"))
        except Exception as e:
            self.root.after(0, lambda: self.activate_btn.config(state='normal', text="激活许可证"))
            self.root.after(0, lambda: self.log_status(f"❌ 激活异常: {str(e)}"))
            self.root.after(0, lambda: messagebox.showerror("错误",
                f"激活过程发生异常:\n\n{str(e)}"))

    def heartbeat(self):
        """心跳验证"""
        try:
            headers = {}
            if self.token:
                headers["Authorization"] = f"Bearer {self.token}"

            response = requests.post(
                f"{self.server_url}/api/heartbeat",
                json={
                    "key": self.license_key,
                    "hwid": self.hwid
                },
                headers=headers,
                timeout=10
            )

            if response.status_code == 200:
                data = response.json()
                if data.get("status") == "alive":
                    return True

            return False
        except:
            return False

    def start_heartbeat(self):
        """启动心跳监控"""
        if self.heartbeat_running:
            return

        self.heartbeat_running = True
        self.log_status("💓 心跳监控已启动")

        def heartbeat_loop():
            retry_count = 0
            max_retries = 3

            while self.heartbeat_running:
                time.sleep(30)  # 30秒心跳间隔

                if self.heartbeat():
                    retry_count = 0
                    self.root.after(0, lambda: self.log_status("💓 心跳验证成功"))
                else:
                    retry_count += 1
                    self.root.after(0, lambda: self.log_status(
                        f"⚠️ 心跳验证失败 ({retry_count}/{max_retries})"))

                    if retry_count >= max_retries:
                        self.root.after(0, lambda: self.log_status("❌ 许可证验证失败,程序即将退出"))
                        self.root.after(0, lambda: messagebox.showerror("错误",
                            "许可证验证失败!\n\n可能原因:\n" +
                            "• 许可证已过期\n" +
                            "• 许可证已被封禁\n" +
                            "• 网络连接中断\n\n" +
                            "程序将自动退出"))
                        self.root.after(0, self.root.quit)
                        break

        threading.Thread(target=heartbeat_loop, daemon=True).start()

    def start_main_app(self):
        """启动主应用程序"""
        self.log_status("🎉 应用程序已启动!")
        self.log_status("现在可以正常使用所有功能")

        # 禁用激活按钮和输入框
        self.activate_btn.config(state='disabled')
        self.license_entry.config(state='disabled')

        # 这里添加你的应用主逻辑
        # 例如: 打开主窗口,启动功能等

def main():
    root = tk.Tk()
    app = LicenseApp(root)
    root.mainloop()

if __name__ == "__main__":
    main()
