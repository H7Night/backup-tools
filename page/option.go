package page

import (
	"backup-tools/tools"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func InitTab1() (*fyne.Container, *widget.Entry, *widget.Entry) {

	// 创建只读的日志标签
	logContent := ""
	logLabel := widget.NewLabel(logContent)
	logLabel.Wrapping = fyne.TextWrapWord // 自动换行
	// 定义一个滚动容器用于日志标签
	logScroll := container.NewScroll(logLabel)
	logScroll.SetMinSize(fyne.NewSize(0, 100)) // 设置日志区域高度

	deviceList := tools.GetConnectedDevices()
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
		deviceList = tools.GetConnectedDevices()
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
	// 拷贝按钮
	copyBtn := widget.NewButton("拷贝", func() {
		deviceID := deviceSelect.Selected
		if deviceID == "" {
			fmt.Println("没有选择设备")
			logLabel.SetText(logLabel.Text + "\n" + "没有选择设备")
			return
		}
		srcPath := srcDir.Text
		destPath := destDir.Text
		err := tools.CopyFilesToLocal(deviceID, srcPath, destPath)
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
