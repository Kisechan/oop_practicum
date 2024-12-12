package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {

	content := container.NewBorder(
		widget.NewLabelWithStyle("Top", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle("Bottom", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabel("Left"),
		widget.NewLabel("Right"),
		widget.NewLabel("Border Container"))
	return container.NewCenter(content)
}
