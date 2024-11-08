package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetConnectedDevices() []string {
	cmd := exec.Command("adb", "devices")
	fmt.Printf("执行命令: %s\n", cmd)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("拉取设备失败", err)
		return []string{}
	}
	lines := strings.Split(string(output), "\n")
	var devices []string
	for _, line := range lines {
		if strings.Contains(line, "device") && !strings.Contains(line, "List of devices") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				devices = append(devices, fields[0])
			}
		}
	}
	return devices
}

func CopyFilesToLocal(deviceID, remotePath, localPath string) error {
	cmd := exec.Command("adb", "-s", deviceID, "pull", remotePath, localPath)
	println("执行命令：", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("错误: %s\n", output)
		return fmt.Errorf("%w", err)
	}
	return nil
}
