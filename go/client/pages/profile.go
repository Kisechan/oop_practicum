package pages

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 个人主页
func CreateProfilePage() fyne.CanvasObject {
	// 用户信息
	profileInfo := container.NewVBox(
		widget.NewLabel("个人主页"),
		widget.NewLabel("用户名: Alice"),
		widget.NewLabel("邮箱: alice@example.com"),
	)

	// 修改信息按钮
	editButton := widget.NewButton("修改信息", func() {
		fmt.Println("修改个人信息")
	})

	// 布局
	return container.NewVBox(
		widget.NewLabelWithStyle("个人主页", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		profileInfo,
		editButton,
	)
}
