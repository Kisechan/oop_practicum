package ui

import (
	"web_sql/rep"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {

	content := container.NewBorder(nil, nil, nil, nil)
	return container.NewCenter(content)
}

func tableScreen[T any](_ fyne.Window) fyne.CanvasObject {
	buttons := container.NewGridWithColumns(3,
		widget.NewButton("增加", func() {
		}),
		widget.NewButton("删除", func() {
		}),
		widget.NewButton("修改", func() {
		}),
	)
	entry := widget.NewEntry()
	entry.SetPlaceHolder("查询")
	search := widget.NewButton("查询", func() {
		rep.GetField[T](rep.DB, "ID", entry.Text)
	})
	searchLine := container.NewGridWithColumns(2, entry, search)
	curd := container.NewVBox(widget.NewSeparator(), buttons, searchLine, widget.NewSeparator())
	// 增删改查按钮

	return container.NewBorder(curd, nil, nil, nil, nil)
}
