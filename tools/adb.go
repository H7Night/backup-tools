package tools

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
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

func CopyFiles(deviceID, srcPath, destPath string, onProgress func(float64)) error {
	cmd := exec.Command("adb", "-s", deviceID, "pull", srcPath, destPath)
	fmt.Printf("执行命令：%s\n", cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	progressRegex := regexp.MustCompile(`$begin:math:text$(\\d+)%$end:math:text$`) // 匹配类似 "(75%)" 的进度格式
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("输出: ", line)
		match := progressRegex.FindStringSubmatch(line)
		if len(match) > 1 {
			percent, _ := strconv.Atoi(match[1])
			onProgress(float64(percent) / 100) // 将进度百分比传递给回调
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return cmd.Wait()
}
