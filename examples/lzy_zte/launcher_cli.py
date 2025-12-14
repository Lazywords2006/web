#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
程序启动器 - 命令行版本
无需 GUI，纯命令行界面，兼容所有 Python 版本
"""

import sys
import os
import subprocess
import uuid
import platform
import time
import json
import requests

class LicenseActivator:
    def __init__(self):
        # 配置
        self.server_url = "http://localhost:8080"
        self.target_exe = "lzy_zte_12.10.exe"
        self.license_file = "license.dat"

        # 获取硬件ID
        self.hwid = self.get_hardware_id()

    def get_hardware_id(self):
        """生成硬件ID"""
        try:
            # macOS
            if platform.system() == 'Darwin':
                result = subprocess.run(['system_profiler', 'SPHardwareDataType'],
                                      capture_output=True, text=True)
                for line in result.stdout.split('\n'):
                    if 'Serial Number' in line or 'UUID' in line:
                        return line.split(':')[1].strip()

            # Windows
            elif platform.system() == 'Windows':
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

    def print_header(self):
        """打印标题"""
        print("\n" + "="*60)
        print("  🔐 许可证验证系统")
        print("="*60 + "\n")

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
                        print("📋 找到已保存的许可证，正在验证...")

                        if self.verify_license(license_key):
                            print("✅ 许可证验证通过！")
                            return license_key
                        else:
                            print("❌ 许可证验证失败")
                            os.remove(self.license_file)
            except Exception as e:
                print(f"⚠️  读取许可证失败: {e}")

        return None

    def verify_license(self, license_key):
        """验证许可证"""
        try:
            response = requests.post(
                f"{self.server_url}/api/activate",
                json={"key": license_key, "hwid": self.hwid},
                timeout=5
            )
            return response.status_code == 200
        except Exception as e:
            print(f"⚠️  连接服务器失败: {e}")
            return False

    def activate_license(self, license_key):
        """激活许可证"""
        print(f"\n正在激活许可证...")
        print(f"设备ID: {self.hwid[:32]}...")

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

                print("\n" + "="*60)
                print("✅ 许可证激活成功！")
                print("="*60)

                if 'expires_at' in data:
                    print(f"📅 过期时间: {data['expires_at']}")
                if 'product_name' in data:
                    print(f"📦 产品名称: {data['product_name']}")

                return True

            else:
                error_msg = response.json().get('error', '未知错误')
                print("\n" + "="*60)
                print(f"❌ 激活失败: {error_msg}")
                print("="*60)
                return False

        except requests.exceptions.ConnectionError:
            print("\n" + "="*60)
            print("❌ 无法连接到许可证服务器")
            print("="*60)
            print("\n请确保服务器正在运行:")
            print("  cd server")
            print("  ./server")
            return False

        except Exception as e:
            print(f"\n❌ 激活失败: {str(e)}")
            return False

    def launch_program(self):
        """启动目标程序"""
        if not os.path.exists(self.target_exe):
            print(f"\n❌ 找不到程序文件: {self.target_exe}")
            return False

        try:
            print(f"\n🚀 正在启动程序: {self.target_exe}")

            # 启动目标程序
            if platform.system() == 'Windows':
                subprocess.Popen([self.target_exe])
            else:
                # macOS/Linux 使用 Wine
                subprocess.Popen(['wine', self.target_exe])

            print("✅ 程序已启动！")
            return True

        except Exception as e:
            print(f"\n❌ 启动失败: {str(e)}")

            if platform.system() != 'Windows':
                print("\n提示: 在 macOS/Linux 上运行 Windows 程序需要 Wine")
                print("安装 Wine: brew install wine-stable")

            return False

    def run(self):
        """运行启动器"""
        self.print_header()

        # 检查保存的许可证
        saved_key = self.check_saved_license()

        if not saved_key:
            # 需要输入许可证
            print("请输入许可证密钥:")
            license_key = input("密钥: ").strip()

            if not license_key:
                print("\n❌ 许可证密钥不能为空")
                return

            if not self.activate_license(license_key):
                return

        # 启动程序
        print("\n" + "-"*60)
        self.launch_program()
        print("-"*60 + "\n")

if __name__ == "__main__":
    try:
        activator = LicenseActivator()
        activator.run()
    except KeyboardInterrupt:
        print("\n\n已取消")
        sys.exit(0)
    except Exception as e:
        print(f"\n❌ 发生错误: {e}")
        sys.exit(1)
