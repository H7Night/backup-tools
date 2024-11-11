package page

import (
	"backup-tools/tools"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func InitTab1(w fyne.Window) (*fyne.Container, *widget.Entry, *widget.Entry) {

	// 创建只读的日志标签
	logContent := "日志"
	logLabel := widget.NewLabel(logContent)
	logLabel.Wrapping = fyne.TextWrapWord // 自动换行
	// 定义一个滚动容器用于日志标签
	logScroll := container.NewScroll(logLabel)
	logScroll.SetMinSize(fyne.NewSize(0, 300)) // 设置日志区域高度

	deviceList := tools.GetConnectedDevices()
	deviceSelect := widget.NewSelect(deviceList, func(value string) {
		fmt.Println("选择设备: ", value)
		logLabel.SetText(logLabel.Text + "\n" + "选择设备:" + value)
	})
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
			logLabel.SetText(logLabel.Text + "\n" + "没有找到可用设备")
		}
		logScroll.ScrollToBottom()
	})
	// 选择deivce和刷新按钮
	deviceContainer := container.NewBorder(
		nil, nil, nil, getDevicesBtn, deviceSelect)

	srcDir := widget.NewEntry()
	srcDir.SetPlaceHolder("输入源目录/文件")
	destDir := widget.NewEntry()
	destDir.SetPlaceHolder("输入保存目录")

	selectSrcBtn := widget.NewButton(">", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println("打开文件夹出错:", err)
				logLabel.SetText(logLabel.Text + "\n" + "打开文件夹出错:" + err.Error())
				return
			}
			if uri != nil {
				srcDir.SetText(filepath.Join(uri.Path()))
				logLabel.SetText(logLabel.Text + "\n" + "选中目录：" + uri.Path())
			}
		}, w)
	})
	selectDestBtn := widget.NewButton(">", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println("打开文件夹出错:", err)
				logLabel.SetText(logLabel.Text + "\n" + "打开文件夹出错:" + err.Error())
				return
			}
			if uri != nil {
				destDir.SetText(filepath.Join(uri.Path()))
				logLabel.SetText(logLabel.Text + "\n" + "选中目录：" + uri.Path())
			}
		}, w)
	})
	// src和dest 容器,使用 BorderLayout 让按钮和输入框按比例分布
	srcDirContainer := container.NewBorder(
		nil, nil, nil, selectSrcBtn, srcDir)
	destDirContainer := container.NewBorder(
		nil, nil, nil, selectDestBtn, destDir)

	progressBar := widget.NewProgressBar()
	progressBar.Hide()
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
		logLabel.SetText(logLabel.Text + "\n" + "拷贝ing...")
		err := tools.CopyFiles(deviceID, srcPath, destPath)
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
		deviceContainer,
		srcDirContainer,
		destDirContainer,
		copyBtn,
		logScroll,
	)
	return t1, srcDir, destDir
}
