package a

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetConnectedDevices() []string {
	cmd := exec.Command("adb", "devices")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error fetching device:", err)
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
