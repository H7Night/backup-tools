package tools

import (
	"fmt"
	"os/exec"
)

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
