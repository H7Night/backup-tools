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

	// 获取 stdout 和 stderr 的输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("无法获取标准输出: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("无法获取标准错误: %w", err)
	}

	// 启动命令执行
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("命令启动失败: %w", err)
	}

	// 创建两个 scanner 用于分别读取 stdout 和 stderr
	scannerOut := bufio.NewScanner(stdout)
	scannerErr := bufio.NewScanner(stderr)
	progressRegex := regexp.MustCompile(`$begin:math:display$(\\d+)%$end:math:display$`)

	// 定义一个 goroutine 读取 stdout
	go func() {
		for scannerOut.Scan() {
			line := scannerOut.Text()
			fmt.Println("stdout:", line) // 打印 stdout 行

			match := progressRegex.FindStringSubmatch(line)
			if len(match) > 1 {
				percentStr := match[1]
				percent, err := strconv.ParseFloat(percentStr, 64)
				if err == nil {
					onProgress(percent / 100.0) // 更新进度条
				}
			}
		}
	}()

	// 定义一个 goroutine 读取 stderr
	go func() {
		for scannerErr.Scan() {
			line := scannerErr.Text()
			fmt.Println("stderr:", line) // 打印 stderr 行

			match := progressRegex.FindStringSubmatch(line)
			if len(match) > 1 {
				percentStr := match[1]
				percent, err := strconv.ParseFloat(percentStr, 64)
				if err == nil {
					onProgress(percent / 100.0) // 更新进度条
				}
			}
		}
	}()

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("命令执行失败: %w", err)
	}

	return nil
}
