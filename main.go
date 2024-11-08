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
	t1, srcDirEntry, destDirEntry := initTab1()
	t2 := initTab2(srcDirEntry, destDirEntry)
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

func initTab1() (*fyne.Container, *widget.Entry, *widget.Entry) {

	// 创建只读的日志标签
	logContent := ""
	logLabel := widget.NewLabel(logContent)
	logLabel.Wrapping = fyne.TextWrapWord // 自动换行
	// 定义一个滚动容器用于日志标签
	logScroll := container.NewScroll(logLabel)
	logScroll.SetMinSize(fyne.NewSize(0, 100)) // 设置日志区域高度

	deviceList := a2.GetConnectedDevices()
	deviceSelect := widget.NewSelect(deviceList, func(value string) {
		fmt.Println("选择设备: ", value)
		logLabel.SetText(logLabel.Text + "\n" + "选择设备:" + value)
	})

	srcDir := widget.NewEntry()
	srcDir.SetPlaceHolder("输入源目录/文件")
	destDir := widget.NewEntry()
	destDir.SetPlaceHolder("输入保存目录")

	// 刷新设备按钮
	getDevicesBtn := widget.NewButton("刷新", func() {
		deviceList = a2.GetConnectedDevices()
		deviceSelect.Options = deviceList
		deviceSelect.Refresh()
		if deviceList != nil {
			fmt.Println("已重新获取设备:")
			logLabel.SetText(logLabel.Text + "\n" + "已重新获取设备:")
			for i, item := range deviceList {
				fmt.Printf("%d: %s\n", i, item)
				logLabel.SetText(logLabel.Text + "\n" + item)
			}
		} else {
			deviceSelect.Selected = ""
			fmt.Println("没有找到可用设备")
			logLabel.SetText(logLabel.Text + "\n" + "没有选择设备")
		}
		logScroll.ScrollToBottom()
	})

	copyBtn := widget.NewButton("拷贝", func() {
		deviceID := deviceSelect.Selected
		if deviceID == "" {
			fmt.Println("没有选择设备")
			logLabel.SetText(logLabel.Text + "\n" + "没有选择设备")
			return
		}
		srcPath := srcDir.Text
		destPath := destDir.Text
		err := a2.CopyFilesToLocal(deviceID, srcPath, destPath)
		if err != nil {
			fmt.Println("拷贝失败:", err)
			logLabel.SetText(logLabel.Text + "\n" + "拷贝失败:" + err.Error())
		} else {
			fmt.Println("成功")
			logLabel.SetText(logLabel.Text + "\n" + "成功")
		}
		logScroll.ScrollToBottom()
	})
	t1 := container.NewVBox(
		container.NewHBox(deviceSelect, getDevicesBtn),
		srcDir,
		destDir,
		copyBtn,
		logScroll,
	)
	return t1, srcDir, destDir
}

func initTab2(srcDirEntry, destDirEntry *widget.Entry) *fyne.Container {
	// 读取配置
	config, err := a2.LoadConfig()
	if err != nil {
		fmt.Println("加载配置失败:", err)
		return container.NewVBox(widget.NewLabel("加载配置失败"))
	}
	// 配置选择项
	profileOptions := make([]string, 0, len(config.Profiles))
	for p := range config.Profiles {
		profileOptions = append(profileOptions, p)
	}
	// 定义绑定对象，存储选中配置的目录
	srcBind := binding.NewString()
	destBind := binding.NewString()
	// 当前选择项
	selectedProfile := binding.NewString()
	profileSelect := widget.NewSelect(profileOptions, func(value string) {
		println("选择配置：", value)
		selectedProfile.Set(value)
		// 获取并设置对应配置
		profileConfig := config.Profiles[value]
		srcBind.Set(profileConfig.SrcDir)
		destBind.Set(profileConfig.DestDir)
	})
	// tab2中的配置
	confSrcDir := widget.NewEntryWithData(srcBind)
	confDestDir := widget.NewEntryWithData(destBind)

	// 默认选第一个
	if len(profileOptions) > 0 {
		profileSelect.SetSelected(profileOptions[0])
	}
	// 应用按钮，修改tab1配置用
	applyBtn := widget.NewButton("应用", func() {
		selectedProfile := profileSelect.Selected
		if selectedProfile == "" {
			println("未选择配置")
			return
		}
		selectedConfig := config.Profiles[selectedProfile]
		// 应用到tab1中的src和dest
		srcDirEntry.SetText(selectedConfig.SrcDir)
		destDirEntry.SetText(selectedConfig.DestDir)
		println("应用配置：", selectedProfile)
	})
	// 保存按钮，修改配置文件用
	saveBtn := widget.NewButton("保存", func() {
		selected, _ := selectedProfile.Get()
		if selected == "" {
			println("未选择配置")
			return
		}
		// 从Tab2的输入框中获取新设置的值
		newSrcDir, _ := srcBind.Get()
		newDestDir, _ := destBind.Get()
		// 更新配置文件内容
		if selectedConfig, exists := config.Profiles[selected]; exists {
			selectedConfig.SrcDir = newSrcDir
			selectedConfig.DestDir = newDestDir
			config.Profiles[selected] = selectedConfig
			// 保存到文件
			err := a2.SaveConfig(config)
			if err != nil {
				fmt.Println("保存配置失败:", err)
			} else {
				fmt.Println("保存配置成功")
			}
		} else {
			println("配置不存在")
		}
	})

	t2 := container.NewVBox(
		profileSelect,
		confSrcDir,
		confDestDir,
		applyBtn,
		saveBtn,
	)
	return t2
}
