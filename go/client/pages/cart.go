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

// 购物车页面
func CreateCartPage() fyne.CanvasObject {
	// 初始化购物车内容
	cartItems := container.NewVBox()
	label := widget.NewLabel("您尚未登录")

	// 获取购物车信息
	fetchCartItems(cartItems)

	// 刷新按钮
	refreshButton := widget.NewButtonWithIcon("刷新", theme.ViewRefreshIcon(), func() {
		// 清空购物车内容
		cartItems.Objects = nil
		// 重新获取购物车信息
		fetchCartItems(cartItems)
		if currentUser != nil {
			label.Hide()
		}
	})

	// 可滚动的购物车商品信息
	scrollContainer := container.NewScroll(cartItems)

	// 整体布局
	return container.NewBorder(
		container.NewVBox(
			widget.NewRichText(
				&widget.TextSegment{
					Text: "购物车",
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

// 获取购物车信息
func fetchCartItems(cartItems *fyne.Container) {
	// 发送 GET 请求获取购物车信息
	if currentUser == nil {
		fmt.Println("未登录")
		return
	}
	resp, err := http.Get("http://localhost:8080/cart/items?user_id=" + fmt.Sprintf("%d", currentUser.ID))
	if err != nil {
		fmt.Println("获取购物车信息失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败:", err)
		return
	}

	// 解析购物车信息
	var cartResponse struct {
		CartItems []Cart `json:"cart_items"`
	}
	if err := json.Unmarshal(body, &cartResponse); err != nil {
		fmt.Println("解析购物车信息失败:", err)
		return
	}

	// 展示购物车信息
	for _, cartItem := range cartResponse.CartItems {
		cartCard := createCartCard(cartItem)
		cartItems.Add(cartCard)
	}
}

// 创建购物车卡片
func createCartCard(cartItem Cart) fyne.CanvasObject {
	// 商品信息
	productInfo := widget.NewLabel(fmt.Sprintf(
		"价格: ￥%.2f\n数量: %d\n总价: ￥%.2f",
		cartItem.Product.Price,
		cartItem.Quantity,
		cartItem.Product.Price*float64(cartItem.Quantity),
	))
	productInfo.Wrapping = fyne.TextWrapWord

	// 删除按钮
	deleteButton := widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {
		fmt.Println("删除商品:", cartItem.Product.Name)
		// 调用 API 删除商品
		deleteCartItem(cartItem.ID)
	})

	// 购买按钮
	buyButton := widget.NewButtonWithIcon("购买", theme.ConfirmIcon(), func() {
		fmt.Println("购买商品:", cartItem.Product.Name)
		// 调用 API 购买商品
		buyCartItem(cartItem.ID)
	})

	// 按钮容器
	buttonContainer := container.NewGridWithColumns(
		2,
		deleteButton,
		buyButton,
	)

	// 卡片内容
	cardContent := container.NewVBox(
		container.NewHBox(widget.NewIcon(theme.FileIcon()), widget.NewLabel(cartItem.Product.Name)),
		productInfo,
		widget.NewSeparator(),
		buttonContainer,
	)

	// 创建卡片
	return widget.NewCard(
		"",
		"",
		cardContent,
	)
}

// 删除购物车项
func deleteCartItem(cartItemID int) {
	// 调用 API 删除购物车项
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/cart/items/%d", cartItemID), nil)
	if err != nil {
		fmt.Println("创建删除请求失败:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("删除购物车项失败:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("删除成功")
	} else {
		fmt.Println("删除失败")
	}
}

// 购买购物车项
func buyCartItem(cartItemID int) {
	// 调用 API 购买购物车项
	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/cart/items/%d/buy", cartItemID), nil)
	if err != nil {
		fmt.Println("创建购买请求失败:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("购买购物车项失败:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("购买成功")
	} else {
		fmt.Println("购买失败")
	}
}
