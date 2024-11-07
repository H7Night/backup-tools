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
	w.Resize(fyne.NewSize(480, 360))
	w.ShowAndRun()
}

func initTab1() *fyne.Container {
	deviceList := a2.GetConnectedDevices()
	deviceSelect := widget.NewSelect(deviceList, func(value string) {
		fmt.Println("选择设备: ", value)
	})
	// 加载配置
	conf, err := a2.LoadConfig()
	if err != nil {
		println(err)
	}

	srcDir := widget.NewEntry()
	srcDir.SetPlaceHolder("输入源目录/文件")
	destDir := widget.NewEntry()
	destDir.SetPlaceHolder("输入保存目录")

	srcDir.Text = conf.SrcDir
	destDir.Text = conf.DestDir

	// 刷新设备按钮
	getDevicesBtn := widget.NewButton("刷新", func() {
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
			fmt.Println("没有找到可用设备")
		}
	})

	copyBtn := widget.NewButton("拷贝", func() {
		deviceID := deviceSelect.Selected
		if deviceID == "" {
			fmt.Println("没有选择设备")
			return
		}
		srcPath := srcDir.Text
		destPath := destDir.Text
		err := a2.CopyFilesToLocal(deviceID, srcPath, destPath)
		if err != nil {
			fmt.Println("拷贝失败:", err)
		} else {
			fmt.Println("成功")
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
		fmt.Println("加载配置失败:", err)
		return container.NewVBox(widget.NewLabel("加载配置失败"))
	}

	cSrcDir := binding.NewString()
	cSrcDir.Set(config.SrcDir)
	cDestDir := binding.NewString()
	cDestDir.Set(config.DestDir)

	confSrcDir := widget.NewEntryWithData(cSrcDir)
	confDestDir := widget.NewEntryWithData(cDestDir)

	saveBtn := widget.NewButton("保存", func() {
		srcDir, _ := cSrcDir.Get()
		destDir, _ := cDestDir.Get()
		config.SrcDir = srcDir
		config.DestDir = destDir

		err := a2.SaveConfig(config)
		if err != nil {
			fmt.Println("保存配置失败:", err)
		} else {
			fmt.Println("保存配置成功")
		}
	})

	t2 := container.NewVBox(
		confSrcDir,
		confDestDir,
		saveBtn,
	)
	return t2
}
