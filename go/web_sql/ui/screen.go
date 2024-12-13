package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func welcomeScreen(win fyne.Window) fyne.CanvasObject {

	text := widget.NewRichText(
		&widget.TextSegment{
			Text: "欢迎",
			Style: widget.RichTextStyle{
				SizeName: theme.SizeNameHeadingText,
			},
		},
	)
	content := container.NewBorder(text, nil, nil, nil)
	return container.NewCenter(content)
}
