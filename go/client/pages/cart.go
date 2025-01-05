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
	var (
		cartItems []Cart
		cartList  *widget.List
	)
	label := widget.NewLabel("您尚未登录")
	// 创建可滑动的列表
	cartList = widget.NewList(
		func() int {
			return len(cartItems)
		},
		func() fyne.CanvasObject {
			// 创建列表项的模板
			nameLabel := widget.NewLabel("商品名称") // 商品名称
			nameLabel.TextStyle.Bold = true

			priceLabel := widget.NewLabel("价格") // 商品价格
			priceLabel.Alignment = fyne.TextAlignTrailing

			quantityLabel := widget.NewLabel("数量") // 商品数量
			quantityLabel.Alignment = fyne.TextAlignTrailing

			totalLabel := widget.NewLabel("总价") // 商品总价
			totalLabel.Alignment = fyne.TextAlignTrailing
			totalLabel.TextStyle.Bold = true

			deleteButton := widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {}) // 删除按钮
			deleteButton.Importance = widget.LowImportance

			buyButton := widget.NewButtonWithIcon("购买", theme.ConfirmIcon(), func() {}) // 购买按钮
			buyButton.Importance = widget.HighImportance

			// 商品信息布局
			infoContainer := container.NewGridWithColumns(
				3,
				nameLabel,
				priceLabel,
				quantityLabel,
			)

			// 操作按钮布局
			buttonContainer := container.NewGridWithColumns(
				2,
				deleteButton,
				buyButton,
			)

			// 整体布局
			return widget.NewCard(
				"", // 卡片标题（留空）
				"", // 卡片副标题（留空）
				container.NewVBox(
					container.NewHBox(
						widget.NewIcon(theme.DocumentIcon()), // 图标
						infoContainer,
					),
					totalLabel,
					widget.NewSeparator(),
					buttonContainer,
				),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			cartItem := cartItems[i]
			card := o.(*widget.Card)
			cardContent := card.Content.(*fyne.Container)

			// 获取列表项中的组件
			infoContainer := cardContent.Objects[0].(*fyne.Container).Objects[1].(*fyne.Container)
			totalLabel := cardContent.Objects[1].(*widget.Label)
			buttonContainer := cardContent.Objects[3].(*fyne.Container)

			nameLabel := infoContainer.Objects[0].(*widget.Label)
			priceLabel := infoContainer.Objects[1].(*widget.Label)
			quantityLabel := infoContainer.Objects[2].(*widget.Label)

			deleteButton := buttonContainer.Objects[0].(*widget.Button)
			buyButton := buttonContainer.Objects[1].(*widget.Button)

			// 设置商品信息
			nameLabel.SetText(cartItem.Product.Name)
			priceLabel.SetText(fmt.Sprintf("￥%.2f", cartItem.Product.Price))
			quantityLabel.SetText(fmt.Sprintf("x %d", cartItem.Quantity))
			totalLabel.SetText(fmt.Sprintf("总价: ￥%.2f", cartItem.Product.Price*float64(cartItem.Quantity)))

			// 设置按钮事件
			deleteButton.OnTapped = func() {
				fmt.Println("删除商品:", cartItem.Product.Name)
				deleteCartItem(cartItem.ID, &cartItems, cartList)
			}

			buyButton.OnTapped = func() {
				fmt.Println("购买商品:", cartItem.Product.Name)
				showCouponSelectionDialogCart(currentUser.ID, cartItem)
				deleteCartItem(cartItem.ID, &cartItems, cartList)
			}
		},
	)

	// 获取购物车信息
	fetchCartItems(&cartItems, cartList)

	// 刷新按钮
	refreshButton := widget.NewButtonWithIcon("刷新", theme.ViewRefreshIcon(), func() {
		// 清空购物车内容
		cartItems = nil
		// 重新获取购物车信息
		fetchCartItems(&cartItems, cartList)
		if currentUser != nil {
			label.Hide()
		}
	})

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
		cartList,
	)
}

// 获取购物车信息
func fetchCartItems(cartItems *[]Cart, cartList *widget.List) {
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

	// 更新购物车列表
	*cartItems = cartResponse.CartItems
	cartList.Refresh()
}

// 删除购物车项
func deleteCartItem(cartItemID int, cartItems *[]Cart, cartList *widget.List) {
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
		// 从本地列表中移除该项
		for i, item := range *cartItems {
			if item.ID == cartItemID {
				*cartItems = append((*cartItems)[:i], (*cartItems)[i+1:]...)
				break
			}
		}
		// 刷新列表
		cartList.Refresh()
	} else {
		fmt.Println("删除失败")
	}
}
