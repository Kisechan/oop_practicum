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
	label := widget.NewLabel("您尚未登录")

	// 获取订单信息
	fetchOrders(orderItems)

	// 刷新按钮
	refreshButton := widget.NewButtonWithIcon("刷新", theme.ViewRefreshIcon(), func() {
		// 清空订单内容
		orderItems.Objects = nil
		// 重新获取订单信息
		fetchOrders(orderItems)
		if currentUser != nil {
			label.Hide()
		}
	})

	// 可滚动的订单信息
	scrollContainer := container.NewScroll(orderItems)

	// 整体布局
	return container.NewBorder(
		container.NewVBox(
			widget.NewRichText(
				&widget.TextSegment{
					Text: "我的订单",
					Style: widget.RichTextStyle{
						SizeName:  theme.SizeNameHeadingText,
						Alignment: fyne.TextAlignCenter,
						TextStyle: fyne.TextStyle{
							Bold: true,
						},
					},
				},
			),
			label,
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
	resp, err := http.Get(ServerAddress + "orders/" + fmt.Sprintf("%d", currentUser.ID))
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
		orderCard := createOrderCard(order)
		orderItems.Add(orderCard)
	}
}

// 创建订单卡片
func createOrderCard(order Order) fyne.CanvasObject {

	// 订单状态标签
	statusLabel := widget.NewLabelWithStyle(
		fmt.Sprintf("状态: %s", OrderStatus[order.Status]),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	// 订单基本信息
	orderInfo := widget.NewLabel(fmt.Sprintf(
		"商品:     %s\n数量:     %d\n总价:     ￥%.2f\n折扣:     ￥%.2f\n实付:     ￥%.2f\n创建时间:  %s",
		order.Product.Name,
		order.Quantity,
		order.Total,
		order.Discount,
		order.Payable,
		order.CreatedTime.Format("2006-01-02 15:04:05"),
	))
	orderInfo.Wrapping = fyne.TextWrapWord

	// 订单卡片布局
	cardContent := container.NewVBox(
		container.NewHBox(
			widget.NewIcon(theme.DocumentIcon()),
			widget.NewLabel("订单-"+order.OrderNumber),
		),
		widget.NewSeparator(),
		statusLabel,
		orderInfo,
		// widget.NewSeparator(),
	)

	// 创建卡片
	return widget.NewCard(
		"",
		"",
		cardContent,
	)
}
