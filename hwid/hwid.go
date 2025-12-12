package hwid

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Generator 硬件ID生成器接口
type Generator interface {
	GetHWID() (string, error)
}

// GetHardwareID 获取跨平台的硬件ID（基于CPU、磁盘、主板信息）
// 返回一个SHA256哈希值作为稳定的机器指纹
func GetHardwareID() (string, error) {
	var hwInfo string
	var err error

	switch runtime.GOOS {
	case "windows":
		hwInfo, err = getWindowsHWID()
	case "linux":
		hwInfo, err = getLinuxHWID()
	case "darwin":
		hwInfo, err = getDarwinHWID()
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err != nil {
		return "", fmt.Errorf("failed to get hardware info: %w", err)
	}

	// 使用SHA256生成稳定的硬件ID
	hash := sha256.Sum256([]byte(hwInfo))
	return hex.EncodeToString(hash[:]), nil
}

// getWindowsHWID Windows平台硬件ID获取
func getWindowsHWID() (string, error) {
	var components []string

	// CPU ID (通过WMIC)
	cpuID, err := runCommand("wmic", "cpu", "get", "ProcessorId")
	if err == nil && cpuID != "" {
		components = append(components, strings.TrimSpace(cpuID))
	}

	// 主板序列号
	boardSerial, err := runCommand("wmic", "baseboard", "get", "SerialNumber")
	if err == nil && boardSerial != "" {
		components = append(components, strings.TrimSpace(boardSerial))
	}

	// 磁盘序列号
	diskSerial, err := runCommand("wmic", "diskdrive", "get", "SerialNumber")
	if err == nil && diskSerial != "" {
		components = append(components, strings.TrimSpace(diskSerial))
	}

	if len(components) == 0 {
		return "", fmt.Errorf("no hardware identifiers found on Windows")
	}

	return strings.Join(components, "|"), nil
}

// getLinuxHWID Linux平台硬件ID获取
func getLinuxHWID() (string, error) {
	var components []string

	// CPU信息
	cpuInfo, err := os.ReadFile("/proc/cpuinfo")
	if err == nil {
		// 提取CPU序列号或型号
		lines := strings.Split(string(cpuInfo), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Serial") || strings.Contains(line, "model name") {
				components = append(components, strings.TrimSpace(line))
				break
			}
		}
	}

	// 机器ID（systemd）
	machineID, err := os.ReadFile("/etc/machine-id")
	if err == nil {
		components = append(components, strings.TrimSpace(string(machineID)))
	} else {
		// 备用：尝试 /var/lib/dbus/machine-id
		machineID, err = os.ReadFile("/var/lib/dbus/machine-id")
		if err == nil {
			components = append(components, strings.TrimSpace(string(machineID)))
		}
	}

	// 主板信息（通过dmidecode，需要root权限）
	boardSerial, err := runCommand("dmidecode", "-s", "baseboard-serial-number")
	if err == nil && boardSerial != "" {
		components = append(components, strings.TrimSpace(boardSerial))
	}

	if len(components) == 0 {
		return "", fmt.Errorf("no hardware identifiers found on Linux")
	}

	return strings.Join(components, "|"), nil
}

// getDarwinHWID macOS平台硬件ID获取
func getDarwinHWID() (string, error) {
	var components []string

	// 硬件UUID
	hwUUID, err := runCommand("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	if err == nil && hwUUID != "" {
		// 提取IOPlatformUUID
		for _, line := range strings.Split(hwUUID, "\n") {
			if strings.Contains(line, "IOPlatformUUID") {
				components = append(components, strings.TrimSpace(line))
				break
			}
		}
	}

	// 序列号
	serialNumber, err := runCommand("ioreg", "-l")
	if err == nil && serialNumber != "" {
		for _, line := range strings.Split(serialNumber, "\n") {
			if strings.Contains(line, "IOPlatformSerialNumber") {
				components = append(components, strings.TrimSpace(line))
				break
			}
		}
	}

	// 备用方案：使用system_profiler（较慢）
	if len(components) == 0 {
		hwInfo, err := runCommand("system_profiler", "SPHardwareDataType")
		if err == nil && hwInfo != "" {
			for _, line := range strings.Split(hwInfo, "\n") {
				if strings.Contains(line, "Serial Number") || strings.Contains(line, "Hardware UUID") {
					components = append(components, strings.TrimSpace(line))
				}
			}
		}
	}

	if len(components) == 0 {
		return "", fmt.Errorf("no hardware identifiers found on macOS")
	}

	return strings.Join(components, "|"), nil
}

// runCommand 执行系统命令并返回输出
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
