package pages

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 个人订单页面
func CreateOrdersPage() fyne.CanvasObject {
	// 订单列表
	orderList := container.NewVBox(
		widget.NewLabel("订单列表"),
		widget.NewButton("订单1 - ￥99.99", func() {
			fmt.Println("查看订单1")
		}),
		widget.NewButton("订单2 - ￥199.99", func() {
			fmt.Println("查看订单2")
		}),
	)

	// 布局
	return container.NewVBox(
		widget.NewLabelWithStyle("订单", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		orderList,
	)
}
