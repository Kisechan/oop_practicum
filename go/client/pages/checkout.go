package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func GenerateOrderNumber() string {
	// 获取当前时间戳（格式为年月日时分秒）
	timestamp := time.Now().Format("20060102150405") // 例如：20231025143045

	// 生成一个随机数（范围：1000-9999）
	rand.Seed(time.Now().UnixNano())       // 确保每次运行随机数不同
	randomNumber := rand.Intn(9000) + 1000 // 生成 1000-9999 的随机数

	// 组合成订单号
	orderNumber := fmt.Sprintf("ORDER_%s_%d", timestamp, randomNumber)
	return orderNumber
}

func placeOrderCart(userID, productID, quantity int, couponCode string, discount, payable, total float64) (string, error) {
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

	resp, err := http.Post("http://localhost:8080/orders/checkout", "application/json", bytes.NewBuffer(orderJSON))
	if err != nil {
		return "", fmt.Errorf("发送订单请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下单失败: %s", resp.Status)
	}

	return orderNumber, nil
}

func pollOrderResultCart(orderNumber string) (string, string, error) {
	for {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/orders/checkout/result/%s", orderNumber))
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

func showPurchaseResult(status string, message string) {
	// 显示购买结果
	dialog.ShowInformation("购买结果", Message[status]+"\n"+Message[message], fyne.CurrentApp().Driver().AllWindows()[1])
}

func sendCheckoutOnCoupon(product Product, quantity int, couponCode string, discount float64, payable float64, originalTotal float64) {
	orderNumber, err := placeOrderCart(currentUser.ID, product.ID, quantity, couponCode, discount, payable, originalTotal)

	if err != nil {
		dialog.ShowError(fmt.Errorf("下单失败: %v", err), fyne.CurrentApp().Driver().AllWindows()[1])
		return
	}
	fmt.Println("OrderNumber is", orderNumber)
	// 轮询订单结果
	go func() {
		status, message, err := pollOrderResultCart(orderNumber)
		if err != nil {
			dialog.ShowError(fmt.Errorf("获取订单结果失败: %v", err), fyne.CurrentApp().Driver().AllWindows()[1])
			return
		}

		// 显示购买结果
		showPurchaseResult(status, message)
	}()
}

func showCouponSelectionDialog(userID int, product Product) {
	coupons, err := getAvailableCoupons(userID)
	quantity := 1
	if err != nil {
		dialog.ShowError(fmt.Errorf("获取优惠券失败: %v", err), fyne.CurrentApp().Driver().AllWindows()[1])
		return
	}
	fmt.Println("Get Coupons:", coupons)
	couponOptions := make([]string, len(coupons)+1)
	couponOptions[0] = "不选择优惠券"
	for i, coupon := range coupons {
		couponOptions[i+1] = fmt.Sprintf("%d: 折扣: - %.2f", i+1, coupon.Discount)
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

	quantityEntry := widget.NewEntry()
	quantityEntry.SetText("1")

	totalLabel := widget.NewLabel(fmt.Sprintf("总价： %.2f", total))
	discountLabel := widget.NewLabel(fmt.Sprintf("折扣： -%.2f", discount))
	payableLabel := widget.NewLabel(fmt.Sprintf("实付价： %.2f", payable))

	couponRadio := widget.NewRadioGroup(couponOptions, func(s string) {

		if qty, err := fmt.Sscanf(quantityEntry.Text, "%d", &quantity); err != nil || qty != 1 {
			dialog.ShowInformation("无效数量", "请输入有效的购买数量", fyne.CurrentApp().Driver().AllWindows()[1])
			return
		}

		if s == "不选择优惠券" || s == "" {
			coupon = nil
			total = product.Price * float64(quantity)
			discount = 0.00
			payable = total - discount

			totalLabel.Text = fmt.Sprintf("总价： %.2f", total)
			totalLabel.Refresh()
			discountLabel.Text = fmt.Sprintf("折扣： -%.2f", discount)
			discountLabel.Refresh()
			payableLabel.Text = fmt.Sprintf("实付价： %.2f", payable)
			payableLabel.Refresh()
			return
		}
		var (
			selectedIndex  int
			couponSelected string
		)
		fmt.Sscanf(s, "%d%s", &selectedIndex, &couponSelected)
		coupon = &coupons[selectedIndex-1]
		fmt.Println("Coupon:", *coupon)

		total = product.Price * float64(quantity)
		discount = coupon.Discount
		payable = total - discount
		if payable <= 0 {
			totalLabel.Text = fmt.Sprintf("总价： %.2f", total)
			totalLabel.Refresh()
			discountLabel.Text = "不可选择此优惠券"
			discountLabel.Refresh()
			payableLabel.Text = fmt.Sprintf("实付价： %.2f", total)
			payableLabel.Refresh()
			return
		}
		totalLabel.Text = fmt.Sprintf("总价： %.2f", total)
		totalLabel.Refresh()
		discountLabel.Text = fmt.Sprintf("折扣： -%.2f", discount)
		discountLabel.Refresh()
		payableLabel.Text = fmt.Sprintf("实付价： %.2f", payable)
		payableLabel.Refresh()
	})
	couponRadio.Selected = "不选择优惠券"
	couponRadio.Horizontal = true

	quantityEntry.OnChanged = func(s string) {
		quantityStr := quantityEntry.Text
		if quantityStr == "" {
			return
		}
		q, err := strconv.Atoi(quantityStr)
		if err != nil || quantity <= 0 {
			quantityEntry.SetText("1")
			dialog.ShowError(fmt.Errorf("请输入有效的数量"), fyne.CurrentApp().Driver().AllWindows()[1])
			return
		}
		if q > 999 {
			q = 999
			quantityEntry.SetText("999")
			dialog.ShowError(fmt.Errorf("购买数量不能大于999"), fyne.CurrentApp().Driver().AllWindows()[1])
		}
		quantity = q
		total = product.Price * float64(quantity)
		if coupon != nil {
			discount = coupon.Discount
		} else {
			discount = 0.00
		}
		payable = total - discount
		if payable <= 0 {
			totalLabel.Text = fmt.Sprintf("总价： %.2f", total)
			totalLabel.Refresh()
			discountLabel.Text = "不可选择此优惠券"
			discountLabel.Refresh()
			payableLabel.Text = fmt.Sprintf("实付价： %.2f", total)
			payableLabel.Refresh()
			return
		}
		totalLabel.Text = fmt.Sprintf("总价： %.2f", total)
		totalLabel.Refresh()
		discountLabel.Text = fmt.Sprintf("折扣： -%.2f", discount)
		discountLabel.Refresh()
		payableLabel.Text = fmt.Sprintf("实付价： %.2f", payable)
		payableLabel.Refresh()
	}

	dialog.ShowCustomConfirm("选择优惠券和数量", "确认", "取消", container.NewVBox(
		widget.NewLabel("选择优惠券:"),
		couponRadio,
		widget.NewLabel("购买数量:"),
		quantityEntry,
		widget.NewSeparator(),
		totalLabel,
		discountLabel,
		payableLabel,
	), func(confirmed bool) {
		if !confirmed {
			return
		}
		if payable < 0 {
			dialog.ShowError(fmt.Errorf("不可选择此优惠券"), fyne.CurrentApp().Driver().AllWindows()[1])
			return
		}
		if quantityEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("请输入购买数量"), fyne.CurrentApp().Driver().AllWindows()[1])
			return
		}
		if couponRadio.Selected == "不选择优惠券" || couponRadio.Selected == "" {
			sendCheckoutOnCoupon(product, quantity, "", discount, payable, total)
			return
		}
		sendCheckoutOnCoupon(product, quantity, coupon.Code, discount, payable, total)
	}, fyne.CurrentApp().Driver().AllWindows()[1])
}
