package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateSettingsPage() fyne.CanvasObject {
	// 创建主题切换的单选按钮组
	themeOptions := []string{"关", "开"}
	themeRadio := widget.NewRadioGroup(themeOptions, func(selected string) {
		// 根据选择的主题切换应用的主题
		switch selected {
		case "关":
			fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
		case "开":
			fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
		}
	})
	themeRadio.Horizontal = true // 水平排列单选按钮
	themeRadio.SetSelected("关")  // 默认选择 Light 模式

	// 创建设置页面的布局
	settingsPage := container.NewVBox(
		widget.NewRichText(
			&widget.TextSegment{
				Text: "设置",
				Style: widget.RichTextStyle{
					SizeName:  theme.SizeNameHeadingText,
					Alignment: fyne.TextAlignCenter,
					TextStyle: fyne.TextStyle{
						Bold: true,
					},
				},
			},
		),
		container.NewGridWithColumns(
			2,
			widget.NewLabel("夜间模式："),
			themeRadio,
		))

	return settingsPage
}
