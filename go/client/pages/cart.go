package pages

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 购物车页面
func CreateCartPage() fyne.CanvasObject {
	// 购物车内容
	cartItems := container.NewVBox(
		widget.NewLabel("购物车"),
		widget.NewButton("商品1 - ￥99.99", func() {
			fmt.Println("删除商品1")
		}),
		widget.NewButton("商品2 - ￥199.99", func() {
			fmt.Println("删除商品2")
		}),
	)

	// 结算按钮
	checkoutButton := widget.NewButton("结算", func() {
		fmt.Println("结算购物车")
	})

	// 布局
	return container.NewVBox(
		widget.NewLabelWithStyle("购物车", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		cartItems,
		checkoutButton,
	)
}
