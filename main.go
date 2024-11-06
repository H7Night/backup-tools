package main

import (
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("backup-tools")

	deviceList := getConnectedDevices()
	deviceSelect := widget.NewSelect(deviceList, func(value string) {
		fmt.Println("Selected device: ", value)
	})

	srcDir := widget.NewEntry()
	srcDir.SetPlaceHolder("Enter source directory on device")
	destDir := widget.NewEntry()
	destDir.SetPlaceHolder("Enter destination directory on local machine")

	// 刷新设备按钮
	getDevicesBtn := widget.NewButton("Refresh Devices", func() {
		deviceList = getConnectedDevices() // 获取最新的设备列表
		deviceSelect.Options = deviceList  // 更新下拉框选项
		deviceSelect.Refresh()             // 刷新下拉框以显示新选项
		fmt.Println("已重新获取设备:")
		for i, item := range deviceList {
			fmt.Printf("%d: %s\n", i, item)
		}
	})

	copyBtn := widget.NewButton("Copy", func() {
		deviceID := deviceSelect.Selected
		if deviceID == "" {
			fmt.Println("No device selected")
			return
		}
		srcPath := srcDir.Text
		destPath := destDir.Text
		err := copyFilesToLocal(deviceID, srcPath, destPath)
		if err != nil {
			fmt.Println("Error copying:", err)
		} else {
			fmt.Println("Copy successfully")
		}
	})
	w.SetContent(container.NewVBox(
		container.NewHBox(deviceSelect, getDevicesBtn),
		srcDir,
		destDir,
		copyBtn,
	))
	w.ShowAndRun()
}

func getConnectedDevices() []string {
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

func copyFilesToLocal(deviceID, remotePath, localPath string) error {
	cmd := exec.Command("adb", "-s", deviceID, "pull", remotePath, localPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", output)
		return fmt.Errorf("%w", err)
	}
	fmt.Printf("Command output: %s\n", string(output)) // 输出成功信息
	return nil
}

// /sdcard/Download/test.txt
// /Users/jhonhe/Downloads
