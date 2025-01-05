package pages

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 订单页面
func CreateOrdersPage() fyne.CanvasObject {
	// 初始化订单内容
	orderItems := container.NewVBox()
	Label := widget.NewLabel("您尚未登录")

	// 获取订单信息
	fetchOrders(orderItems)

	// 刷新按钮
	refreshButton := widget.NewButtonWithIcon("刷新", theme.ViewRefreshIcon(), func() {
		// 清空订单内容
		orderItems.Objects = nil
		// 重新获取订单信息
		fetchOrders(orderItems)
		if currentUser != nil {
			Label.Hide()
		}
	})

	// 可滚动的订单信息
	scrollContainer := container.NewScroll(orderItems)

	// 整体布局
	return container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("订单", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			Label,
			refreshButton,
		),
		nil,
		nil,
		nil,
		scrollContainer,
	)
}

// 获取订单信息
func fetchOrders(orderItems *fyne.Container) {
	// 发送 GET 请求获取订单信息
	if currentUser == nil {
		fmt.Println("未登录")
		return
	}
	resp, err := http.Get("http://localhost:8080/orders/" + fmt.Sprintf("%d", currentUser.ID))
	if err != nil {
		fmt.Println("获取订单信息失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败:", err)
		return
	}
	fmt.Println("订单信息:", string(body))
	// 解析订单信息
	var orderResponse struct {
		Orders []Order `json:"orders"`
	}
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		fmt.Println("解析订单信息失败:", err)
		return
	}

	// 展示订单信息
	for _, order := range orderResponse.Orders {
		orderInfo := fmt.Sprintf(
			"订单ID: %d\n商品: %s\n数量: %d\n总价: ￥%.2f\n状态: %s\n创建时间: %s",
			order.ID,
			order.Product.Name,
			order.Quantity,
			order.Total,
			order.Status,
			order.CreatedTime.Format("2006-01-02 15:04:05"),
		)
		orderCard := widget.NewCard(
			fmt.Sprintf("订单 %d", order.ID),
			orderInfo,
			nil,
		)
		orderItems.Add(orderCard)
	}
}
