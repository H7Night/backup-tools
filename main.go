package main

import (
	a2 "backup-tools/a"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	// 加载配置
	conf, err := a2.LoadConfig()
	if err != nil {
		println(err)
	}

	srcDir := widget.NewEntry()
	srcDir.SetPlaceHolder("Enter source on device")
	destDir := widget.NewEntry()
	destDir.SetPlaceHolder("Enter destination on local")

	srcDir.Text = conf.SrcDir
	destDir.Text = conf.DestDir

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
	// 读取配置
	config, err := a2.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return container.NewVBox(widget.NewLabel("Failed to load configuration"))
	}

	cSrcDir := binding.NewString()
	cSrcDir.Set(config.SrcDir)
	cDestDir := binding.NewString()
	cDestDir.Set(config.DestDir)

	confSrcDir := widget.NewEntryWithData(cSrcDir)
	confDestDir := widget.NewEntryWithData(cDestDir)

	saveBtn := widget.NewButton("Save", func() {
		srcDir, _ := cSrcDir.Get()
		destDir, _ := cDestDir.Get()
		config.SrcDir = srcDir
		config.DestDir = destDir

		err := a2.SaveConfig(config)
		if err != nil {
			fmt.Println("Save config error:", err)
		} else {
			fmt.Println("Save successfully")
		}
	})

	t2 := container.NewVBox(
		confSrcDir,
		confDestDir,
		saveBtn,
	)
	return t2
}
