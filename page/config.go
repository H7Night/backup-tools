package page

import (
	"backup-tools/tools"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func InitTab2(srcDirEntry, destDirEntry *widget.Entry) *fyne.Container {
	// 读取配置
	config, err := tools.LoadConfig()
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
			err := tools.SaveConfig(config)
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
