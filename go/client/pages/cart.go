package pages

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 购物车页面
func CreateCartPage() fyne.CanvasObject {
	// 初始化购物车内容
	cartItems := container.NewVBox()
	Label := widget.NewLabel("您尚未登录")
	// 获取购物车信息
	fetchCartItems(cartItems)

	// 刷新按钮
	refreshButton := widget.NewButton("刷新", func() {
		// 清空购物车内容
		cartItems.Objects = nil
		// 重新获取购物车信息
		fetchCartItems(cartItems)
		if currentUser != nil {
			Label.Hide()
		}
	})

	// 删除按钮
	deleteButton := widget.NewButton("删除", func() {
		// 获取选中的商品
		var selectedItems []Cart
		for _, item := range cartItems.Objects {
			if checkbox, ok := item.(*widget.Check); ok && checkbox.Checked {
				// 获取对应的购物车项
				selectedItems = append(selectedItems, getCartItemFromCheckbox(checkbox))
			}
		}
		// 删除选中的商品
		if len(selectedItems) > 0 {
			fmt.Println("删除选中的商品:", selectedItems)
			// 这里可以调用 API 删除选中的商品
		} else {
			fmt.Println("未选中任何商品")
		}
	})

	// 结算按钮
	checkoutButton := widget.NewButton("结算", func() {
		// 获取选中的商品
		var selectedItems []Cart
		for _, item := range cartItems.Objects {
			if checkbox, ok := item.(*widget.Check); ok && checkbox.Checked {
				// 获取对应的购物车项
				selectedItems = append(selectedItems, getCartItemFromCheckbox(checkbox))
			}
		}
		// 结算选中的商品
		if len(selectedItems) > 0 {
			fmt.Println("结算选中的商品:", selectedItems)
		} else {
			fmt.Println("未选中任何商品")
		}
	})

	buttonContainer := container.NewGridWithColumns(3, refreshButton, deleteButton, checkoutButton)

	// 可滚动的购物车商品信息
	scrollContainer := container.NewScroll(cartItems)
	// scrollContainer.SetMinSize(fyne.NewSize(400, 300)) // 设置滚动区域的最小大小

	// 整体布局
	return container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("购物车", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			Label,
			buttonContainer,
		),
		nil,
		nil,
		nil,
		scrollContainer,
	)
}

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
	// fmt.Println("响应体内容:", string(body))

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
		checkbox := widget.NewCheck(
			fmt.Sprintf("%s - ￥%.2f x %d", cartItem.Product.Name, cartItem.Product.Price, cartItem.Quantity),
			nil,
		)
		checkbox.SetChecked(false)
		checkbox.OnChanged = func(checked bool) {
			if checked {
				fmt.Println("选中商品:", cartItem.Product.Name)
			} else {
				fmt.Println("取消选中商品:", cartItem.Product.Name)
			}
		}
		cartItems.Add(checkbox)
	}
}

// 从勾选框获取对应的购物车项
func getCartItemFromCheckbox(checkbox *widget.Check) Cart {
	// 这里假设勾选框的标签格式为 "商品名 - ￥价格 x 数量"
	label := checkbox.Text
	// 解析商品名、价格和数量
	var productName string
	var price float64
	var quantity int
	fmt.Sscanf(label, "%s - ￥%f x %d", &productName, &price, &quantity)
	// 返回对应的购物车项
	return Cart{
		Product: Product{
			Name:  productName,
			Price: price,
		},
		Quantity: quantity,
	}
}
