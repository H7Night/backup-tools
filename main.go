package main

import (
	tabPage "backup-tools/page"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("backup-tools")
	t1, srcDirEntry, destDirEntry := tabPage.InitTab1(w)
	t2 := tabPage.InitTab2(srcDirEntry, destDirEntry)
	tabs := container.NewAppTabs(
		container.NewTabItem("操作", t1),
		container.NewTabItem("设置", t2),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	w.SetContent(container.NewVBox(
		tabs,
	))
	w.Resize(fyne.NewSize(720, 480))
	w.ShowAndRun()
}
