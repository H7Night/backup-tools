package main

import (
	a2 "backup-tools/a"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("backup-tools")
	t1 := initTab1()
	t2 := initTab2()
	tabs := container.NewAppTabs(
		container.NewTabItem("Tab 1", t1),
		container.NewTabItem("Tab 2", t2),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	w.SetContent(container.NewVBox(
		tabs,
	))
	w.ShowAndRun()
}

// /sdcard/Download/test.txt
// /Users/jhonhe/Downloads

func initTab1() *fyne.Container {
	deviceList := a2.GetConnectedDevices()
	deviceSelect := widget.NewSelect(deviceList, func(value string) {
		fmt.Println("Selected device: ", value)
	})

	srcDir := widget.NewEntry()
	srcDir.SetPlaceHolder("Enter source on device")
	destDir := widget.NewEntry()
	destDir.SetPlaceHolder("Enter destination on local")

	// 刷新设备按钮
	getDevicesBtn := widget.NewButton("Refresh", func() {
		deviceList = a2.GetConnectedDevices()
		deviceSelect.Options = deviceList
		deviceSelect.Refresh()
		if deviceList != nil {
			fmt.Println("已重新获取设备:")
			for i, item := range deviceList {
				fmt.Printf("%d: %s\n", i, item)
			}
		} else {
			deviceSelect.Selected = ""
			fmt.Println("no devices founded!")
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
		err := a2.CopyFilesToLocal(deviceID, srcPath, destPath)
		if err != nil {
			fmt.Println("Error copying:", err)
		} else {
			fmt.Println("Copy successfully")
		}
	})
	t1 := container.NewVBox(
		container.NewHBox(deviceSelect, getDevicesBtn),
		srcDir,
		destDir,
		copyBtn)
	return t1
}

func initTab2() *fyne.Container {
	t2 := container.NewVBox(
		widget.NewLabel("Hello"),
	)
	return t2
}
