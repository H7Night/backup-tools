package a

import (
	"fmt"
	"os/exec"
)

func CopyFilesToLocal(deviceID, remotePath, localPath string) error {
	cmd := exec.Command("adb", "-s", deviceID, "pull", remotePath, localPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", output)
		return fmt.Errorf("%w", err)
	}
	fmt.Printf("Command output: %s\n", string(output)) // 输出成功信息
	return nil
}
