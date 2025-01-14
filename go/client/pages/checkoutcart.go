package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func getAvailableCoupons(userID int) ([]Coupon, error) {
	resp, err := http.Get(fmt.Sprintf(ServerAddress+"coupons/user/%d", userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var couponsResponse struct {
		// Status string   `json:"status"`
		Data []Coupon `json:"coupons"`
	}
	if err := json.Unmarshal(body, &couponsResponse); err != nil {
		return nil, err
	}

	return couponsResponse.Data, nil
}

func placeOrder(userID, productID, quantity int, couponCode string, discount, payable, total float64) (string, error) {
	orderNumber := GenerateOrderNumber() // 生成订单号

	// 构造订单请求
	orderRequest := map[string]interface{}{
		"user_id":      userID,
		"product_id":   productID,
		"quantity":     quantity,
		"coupon_code":  couponCode,
		"discount":     discount,
		"payable":      payable,
		"total":        total,
		"order_number": orderNumber,
	}

	// 发送订单请求
	orderJSON, err := json.Marshal(orderRequest)
	if err != nil {
		return "", fmt.Errorf("编码订单请求失败: %v", err)
	}

	resp, err := http.Post(ServerAddress+"orders/checkout", "application/json", bytes.NewBuffer(orderJSON))
	if err != nil {
		return "", fmt.Errorf("发送订单请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下单失败: %s", resp.Status)
	}

	return orderNumber, nil
}

func pollOrderResult(orderNumber string) (string, string, error) {
	for {
		resp, err := http.Get(fmt.Sprintf(ServerAddress+"orders/checkout/result/%s", orderNumber))
		if err != nil {
			return "错误", "订单结果查询失败", fmt.Errorf("轮询订单结果失败: %v", err)
		}
		defer resp.Body.Close()

		// 处理不同的 HTTP 状态码
		switch resp.StatusCode {
		case http.StatusNotFound:
			// 订单结果未找到，继续轮询
			time.Sleep(1 * time.Second)
			continue
		case http.StatusOK:
			// 读取响应体
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "错误", "响应体读取失败", fmt.Errorf("读取响应体失败: %v", err)
			}

			// 打印响应体以便调试
			fmt.Println("响应体:", string(body))

			// 解析外层 JSON
			var outerResult struct {
				OrderNumber string `json:"order_number"`
				Result      string `json:"result"` // result 是一个嵌套的 JSON 字符串
			}
			if err := json.Unmarshal(body, &outerResult); err != nil {
				return "错误", "外层订单解析失败", fmt.Errorf("解析外层订单结果失败: %v", err)
			}

			// 解析嵌套的 JSON（result 字段）
			var innerResult struct {
				OrderNumber string `json:"order_number"`
				Status      string `json:"status"`
				Message     string `json:"message"`
			}
			if err := json.Unmarshal([]byte(outerResult.Result), &innerResult); err != nil {
				return "错误", "嵌套读取失败", fmt.Errorf("解析嵌套订单结果失败: %v", err)
			}

			// 返回订单状态
			return innerResult.Status, innerResult.Message, nil
		default:
			// 其他错误状态码
			return "错误", "其他错误", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
		}
	}
}

func showPurchaseResultCart(status string, message string) {
	// 显示购买结果
	dialog.ShowInformation("购买结果", Message[status]+"\n"+Message[message], fyne.CurrentApp().Driver().AllWindows()[0])
}

func sendCheckoutOnCouponCart(product Product, quantity int, couponCode string, discount float64, payable float64, originalTotal float64) {
	orderNumber, err := placeOrder(currentUser.ID, product.ID, quantity, couponCode, discount, payable, originalTotal)

	if err != nil {
		dialog.ShowError(fmt.Errorf("下单失败: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	fmt.Println("OrderNumber is", orderNumber)
	// 轮询订单结果
	go func() {
		fmt.Println("开始轮询订单结果，订单号:", orderNumber)
		status, message, err := pollOrderResult(orderNumber)
		if err != nil {
			dialog.ShowError(fmt.Errorf("获取订单结果失败: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		fmt.Println("购买结果:", status)
		// 显示购买结果
		showPurchaseResultCart(status, message)
	}()
}

func showCouponSelectionDialogCart(userID int, cart Cart) {
	product := cart.Product
	coupons, err := getAvailableCoupons(userID)
	quantity := cart.Quantity
	if err != nil {
		dialog.ShowError(fmt.Errorf("获取优惠券失败: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	fmt.Println("Get Coupons:", coupons)

	// 生成优惠券选项
	couponOptions := make([]string, len(coupons)+1)
	couponOptions[0] = "不选择优惠券"
	for i, coupon := range coupons {
		couponOptions[i+1] = fmt.Sprintf("%d: 折扣 -%.2f", i+1, coupon.Discount) // 修改格式以便解析
	}

	var (
		total    float64
		discount float64
		payable  float64
		coupon   *Coupon
	)

	total = product.Price * float64(quantity)
	discount = 0.00
	payable = total - discount

	totalLabel := widget.NewLabel(fmt.Sprintf("总价： %.2f", total))
	discountLabel := widget.NewLabel(fmt.Sprintf("折扣： -%.2f", discount))
	payableLabel := widget.NewLabel(fmt.Sprintf("实付价： %.2f", payable))

	couponRadio := widget.NewRadioGroup(couponOptions, func(s string) {
		if s == "不选择优惠券" || s == "" {
			coupon = nil
			total = product.Price * float64(quantity)
			discount = 0.00
			payable = total - discount

			totalLabel.SetText(fmt.Sprintf("总价： %.2f", total))
			discountLabel.SetText(fmt.Sprintf("折扣： -%.2f", discount))
			payableLabel.SetText(fmt.Sprintf("实付价： %.2f", payable))
			return
		}

		// 解析用户选择的优惠券索引
		var selectedIndex int
		_, err := fmt.Sscanf(s, "%d:", &selectedIndex) // 提取索引
		if err != nil || selectedIndex < 1 || selectedIndex > len(coupons) {
			dialog.ShowError(fmt.Errorf("无效的优惠券选择"), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}

		// 获取对应的优惠券
		coupon = &coupons[selectedIndex-1]
		fmt.Println("Coupon:", *coupon)

		// 计算总价、折扣和实付价
		total = product.Price * float64(quantity)
		discount = coupon.Discount
		payable = total - discount

		if payable < 0 {
			totalLabel.SetText(fmt.Sprintf("总价： %.2f", total))
			discountLabel.SetText("不可用此优惠券")
			payableLabel.SetText(fmt.Sprintf("实付价： %.2f", total))
			return
		}
		totalLabel.SetText(fmt.Sprintf("总价： %.2f", total))
		discountLabel.SetText(fmt.Sprintf("折扣： -%.2f", discount))
		payableLabel.SetText(fmt.Sprintf("实付价： %.2f", payable))
	})

	couponRadio.Selected = "不选择优惠券"
	couponRadio.Horizontal = true

	dialog.ShowCustomConfirm("选择优惠券和数量", "确认", "取消", container.NewVBox(
		widget.NewLabel("选择优惠券:"),
		couponRadio,
		widget.NewSeparator(),
		totalLabel,
		discountLabel,
		payableLabel,
	), func(confirmed bool) {
		if !confirmed {
			return
		}
		if payable <= 0 {
			dialog.ShowError(fmt.Errorf("不可选择此优惠券"), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		if couponRadio.Selected == "不选择优惠券" || couponRadio.Selected == "" {
			sendCheckoutOnCouponCart(product, quantity, "", discount, payable, total)
			return
		}
		sendCheckoutOnCouponCart(product, quantity, coupon.Code, discount, payable, total)
	}, fyne.CurrentApp().Driver().AllWindows()[0])
}
