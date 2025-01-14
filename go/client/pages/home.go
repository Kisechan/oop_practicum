package pages

import (
	"bytes"
	"client/str"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 图片缓存
var (
	imageCache     = make(map[string]*fyne.StaticResource)
	imageCacheLock sync.RWMutex
)

// 主页面
func CreateHomePage() fyne.CanvasObject {
	// 搜索框
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索商品...")

	// 商品列表
	var (
		products    []Product
		productList *widget.List
	)

	// 创建可滑动的列表
	productList = widget.NewList(
		func() int {
			return len(products)
		},
		func() fyne.CanvasObject {
			// 创建列表项的模板
			image := canvas.NewImageFromResource(theme.FileImageIcon())
			image.FillMode = canvas.ImageFillContain
			image.SetMinSize(fyne.NewSize(70, 70))

			nameLabel := widget.NewLabel("商品名称")
			priceLabel := widget.NewLabel("价格")

			return container.NewGridWithRows(
				2,
				image,
				container.NewGridWithRows(
					3,
					nameLabel,
					priceLabel,
					widget.NewButton("详情", func() {}),
				),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			product := products[i]
			items := o.(*fyne.Container).Objects

			image := items[0].(*canvas.Image)
			nameLabel := items[1].(*fyne.Container).Objects[0].(*widget.Label)
			priceLabel := items[1].(*fyne.Container).Objects[1].(*widget.Label)
			detailButton := items[1].(*fyne.Container).Objects[2].(*widget.Button)

			// 异步加载图片（带缓存）
			if product.Icon != nil && *product.Icon != "" {
				loadImageWithCache(*product.Icon, image)
			} else {
				image.Resource = theme.FileImageIcon()
				image.Refresh()
			}

			nameLabel.SetText(str.CutStr(product.Name, 17))
			nameLabel.TextStyle.Bold = true
			priceLabel.SetText(fmt.Sprintf("￥ %.2f", product.Price))
			priceLabel.Alignment = fyne.TextAlignTrailing
			priceLabel.TextStyle.Italic = true
			priceLabel.Theme().Color(theme.ColorNameError, theme.VariantLight)
			detailButton.OnTapped = func() {
				showProductDetailWindow(product)
			}
		},
	)

	products, _ = fetchProductsFromAPI(searchEntry.Text)
	productList.Refresh() // 刷新列表
	fmt.Println("Product List Refreshed Successfully")
	// 搜索按钮
	searchButton := widget.NewButtonWithIcon("搜索", theme.SearchIcon(), func() {
		products = nil
		products, _ = fetchProductsFromAPI(searchEntry.Text)
		productList.Refresh() // 刷新列表
		fmt.Println("Product List Refreshed Successfully")
	})

	// 布局：顶部固定搜索框，下方为可滑动的列表
	return container.NewBorder(
		container.NewBorder(nil, nil, nil, searchButton, searchEntry), // 顶部搜索框
		nil,         // 底部无内容
		nil,         // 左侧无内容
		nil,         // 右侧无内容
		productList, // 下方可滑动的列表
	)
}

// 异步加载图片（带缓存）
func loadImageWithCache(url string, image *canvas.Image) {
	// 检查缓存
	imageCacheLock.RLock()
	if resource, ok := imageCache[url]; ok {
		image.Resource = resource
		image.Refresh()
		imageCacheLock.RUnlock()
		return
	}
	imageCacheLock.RUnlock()

	// 异步下载图片
	go func() {
		resp, err := http.Get(url)
		if err != nil {
			fyne.LogError("Failed to download image", err)
			return
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fyne.LogError("Failed to read image data", err)
			return
		}

		resource := fyne.NewStaticResource("image.jpg", data)

		// 更新缓存
		imageCacheLock.Lock()
		imageCache[url] = resource
		imageCacheLock.Unlock()

		// 更新图片
		image.Resource = resource
		image.Refresh()
	}()
}

// 从API获取商品数据
func fetchProductsFromAPI(search string) ([]Product, error) {
	var resp *http.Response
	var err error
	if search != "" {
		// searchRequest := map[string]string{
		// 	"name": search,
		// }
		// searchJSON, err := json.Marshal(searchRequest)
		// if err != nil {
		// 	fmt.Println("Error encoding search request:", err)
		// 	return nil, err
		// }

		// 发送搜索请求
		resp, err = http.Post(ServerAddress+"products/search?name="+search, "application/json", nil)
		if err != nil {
			fmt.Println("Error get request:", err)
			return nil, err
		}
	} else {
		resp, err = http.Get(ServerAddress + "products")
		if err != nil {
			fmt.Println("Error get request:", err)
			return nil, err
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var productsResponse ProductsResponse
	if err := json.Unmarshal(body, &productsResponse); err != nil {
		fmt.Println("Error unmarshal reqBody:", err)
		return nil, err
	}
	// fmt.Println("Search:", search)
	// fmt.Println("ProductsResponse:", productsResponse)
	return productsResponse.Products, nil
}

// 显示商品详细窗口
func showProductDetailWindow(product Product) {
	// 创建新窗口
	detailWindow := fyne.CurrentApp().NewWindow("商品详情")

	// 创建商品详细页面
	detailPage := createProductDetailPage(product)

	// 设置窗口内容
	detailWindow.SetContent(detailPage)
	detailWindow.Resize(fyne.NewSize(360, 780)) // 设置窗口大小
	detailWindow.Show()
}

// 创建商品详细页面
func createProductDetailPage(product Product) fyne.CanvasObject {
	// 加载图片
	var image *canvas.Image
	if product.Icon != nil && *product.Icon != "" {
		parsedURL, err := url.Parse(*product.Icon)
		if err != nil {
			fmt.Println("Error parsing image URL:", err)
			image = canvas.NewImageFromResource(theme.FileImageIcon())
		} else {
			uri := storage.NewURI(parsedURL.String())
			if uri == nil {
				fmt.Println("Invalid URI for product:", product.ID, "Name:", product.Name)
				image = canvas.NewImageFromResource(theme.FileImageIcon())
			} else {
				image = canvas.NewImageFromURI(uri)
			}
		}
	} else {
		image = canvas.NewImageFromResource(theme.FileImageIcon())
	}
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(200, 200))

	// 创建商品描述的多行文本
	descriptionLabel := widget.NewLabel(product.Description)
	descriptionLabel.Wrapping = fyne.TextWrapWord // 自动换行

	// 获取评论
	reviews, err := fetchReviewsFromAPI(product.ID)
	if err != nil {
		fmt.Println("Error fetching reviews:", err)
	}
	fmt.Println("获取评论:", reviews)
	// 创建评论列表
	reviewList := container.NewVBox()
	for _, review := range reviews {
		commentItem := widget.NewLabel(review.Comment)
		commentItem.Wrapping = fyne.TextWrapWord
		usernameItem := widget.NewLabel(review.User.Username)
		usernameItem.Wrapping = fyne.TextWrapWord
		usernameItem.TextStyle.Bold = true
		timeItem := widget.NewLabel(review.Time.Format("2006-01-02 15:04:05"))
		timeItem.TextStyle.Italic = true
		timeItem.Alignment = fyne.TextAlignTrailing
		ratingItem := widget.NewLabel(fmt.Sprintf("%s ✰", review.Rating))
		ratingItem.Alignment = fyne.TextAlignTrailing
		reviewItem := container.NewVBox(
			container.NewGridWithColumns(
				2,
				usernameItem,
				ratingItem,
			),
			commentItem,
			timeItem,
			widget.NewSeparator(),
		)
		reviewList.Add(reviewItem)
	}

	// 创建评论输入框
	commentEntry := widget.NewEntry()
	commentEntry.SetPlaceHolder("输入评论...")

	ratingLabel := widget.NewLabel("5")
	ratingSlider := widget.NewSlider(1, 5)
	ratingSlider.Step = 1
	ratingSlider.OnChanged = func(value float64) {
		fmt.Println("用户选择的评分:", int(value))
		ratingLabel.SetText(fmt.Sprintf("%d分", int(value)))
	}
	ratingSlider.SetValue(5)

	// 提交按钮
	submitButton := widget.NewButton("提交评论", func() {
		// 发送评论
		if commentEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("评论不能为空"), fyne.CurrentApp().Driver().AllWindows()[1])
			return
		}
		sendReview(product.ID, fmt.Sprintf("%d", int(ratingSlider.Value)), commentEntry.Text)
		// 刷新评论列表
		reviews, err := fetchReviewsFromAPI(product.ID)
		fmt.Println("获取评论:", reviews)
		if err != nil {
			fmt.Println("Error fetching reviews:", err)
		} else {
			reviewList.Objects = nil // 清空评论列表
			for _, review := range reviews {
				commentItem := widget.NewLabel(review.Comment)
				commentItem.Wrapping = fyne.TextWrapWord
				usernameItem := widget.NewLabel(review.User.Username)
				usernameItem.Wrapping = fyne.TextWrapWord
				usernameItem.TextStyle.Bold = true
				timeItem := widget.NewLabel(review.Time.Format("2006-01-02 15:04:05"))
				timeItem.TextStyle.Italic = true
				timeItem.Alignment = fyne.TextAlignTrailing
				ratingItem := widget.NewLabel(fmt.Sprintf("%s ✰", review.Rating))
				ratingItem.Alignment = fyne.TextAlignTrailing
				reviewItem := container.NewVBox(
					container.NewGridWithColumns(
						2,
						usernameItem,
						ratingItem,
					),
					commentItem,
					timeItem,
					widget.NewSeparator(),
				)
				reviewList.Add(reviewItem)
			}
			reviewList.Refresh() // 刷新评论列表
		}
	})

	// 创建评论输入部分
	commentInput := container.NewGridWithRows(
		2,
		commentEntry,
		container.NewGridWithColumns(
			4,
			widget.NewLabel("评分:"),
			ratingSlider,
			ratingLabel,
			submitButton,
		),
	)

	// 创建详细页面内容
	detailContent := container.NewVBox(
		image,
		widget.NewLabel("商品名称: "+product.Name),
		widget.NewLabel(fmt.Sprintf("价格: ￥%.2f", product.Price)),
		widget.NewLabel("卖家: "+product.Seller),
		widget.NewLabel("商品描述:"),
		descriptionLabel,
		widget.NewSeparator(),
		widget.NewLabel("评论区"),
		commentInput,
		widget.NewSeparator(),
		reviewList,
	)

	// 创建底部按钮
	buttons := container.NewGridWithColumns(
		3,
		widget.NewButton("返回", func() {
			// 关闭当前窗口
			detailWindow := fyne.CurrentApp().Driver().AllWindows()[1]
			detailWindow.Close()
		}),
		widget.NewButton("加入购物车", func() {
			fmt.Println("加入购物车:", product.Name)
			if currentUser == nil {
				// 显示提示框
				dialog.ShowInformation("未登录", "请先登录以加入购物车", fyne.CurrentApp().Driver().AllWindows()[1])
				return
			}

			// 创建数量选择窗口
			quantityEntry := widget.NewEntry()
			quantityEntry.SetPlaceHolder("输入数量")
			quantityEntry.SetText("1") // 默认数量为 1
			quantityEntry.OnChanged = func(s string) {
				quantityStr := quantityEntry.Text
				if quantityStr == "" {
					return
				}
				quantity, err := strconv.Atoi(quantityStr)
				if err != nil || quantity <= 0 {
					quantityEntry.SetText("1")
					dialog.ShowError(fmt.Errorf("请输入有效的数量"), fyne.CurrentApp().Driver().AllWindows()[1])
					return
				}
				if quantity > 999 {
					quantity = 1
					quantityEntry.SetText("1")
					dialog.ShowError(fmt.Errorf("加入购物车数量不能大于999"), fyne.CurrentApp().Driver().AllWindows()[1])
					return
				}
			}
			// 创建确认按钮
			confirmButton := widget.NewButton("确认", func() {
				quantityStr := quantityEntry.Text
				quantity, err := strconv.Atoi(quantityStr)
				if err != nil || quantity <= 0 {
					dialog.ShowError(fmt.Errorf("请输入有效的数量"), fyne.CurrentApp().Driver().AllWindows()[1])
					return
				}

				// 构建购物车项
				cartItem := map[string]interface{}{
					"user_id":    currentUser.ID,
					"product_id": product.ID,
					"quantity":   quantity, // 使用用户输入的数量
				}

				// 发送加入购物车请求
				cartItemJSON, err := json.Marshal(cartItem)
				if err != nil {
					dialog.ShowError(fmt.Errorf("编码购物车项失败: %v", err), fyne.CurrentApp().Driver().AllWindows()[1])
					return
				}

				resp, err := http.Post(ServerAddress+"cart/items", "application/json", bytes.NewBuffer(cartItemJSON))
				if err != nil {
					dialog.ShowError(fmt.Errorf("发送请求失败: %v", err), fyne.CurrentApp().Driver().AllWindows()[1])
					return
				}
				defer resp.Body.Close()

				// 处理响应
				if resp.StatusCode == http.StatusCreated {
					dialog.ShowInformation("成功", "商品已加入购物车", fyne.CurrentApp().Driver().AllWindows()[1])
				} else {
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						dialog.ShowError(fmt.Errorf("读取响应失败: %v", err), fyne.CurrentApp().Driver().AllWindows()[1])
						return
					}
					dialog.ShowError(fmt.Errorf("加入购物车失败: %s", string(body)), fyne.CurrentApp().Driver().AllWindows()[1])
				}
			})

			// 创建对话框内容
			dialogContent := container.NewGridWithRows(
				2,
				widget.NewLabel("加入购物车的数量:"),
				container.NewGridWithColumns(
					2,
					quantityEntry,
					confirmButton,
				),
			)

			// 显示自定义对话框
			dialog.ShowCustom("选择数量", "关闭", dialogContent, fyne.CurrentApp().Driver().AllWindows()[1])
		}),
		widget.NewButton("立即购买", func() {
			// 检查用户是否登录
			if currentUser == nil {
				dialog.ShowInformation("未登录", "请先登录以进行购买", fyne.CurrentApp().Driver().AllWindows()[1])
				return
			}

			showCouponSelectionDialog(currentUser.ID, product)
		}),
	)

	// 将详细内容和按钮组合在一起
	scrollableContent := container.NewVScroll(detailContent)
	return container.NewBorder(nil, buttons, nil, nil, scrollableContent)
}

// 获取评论
func fetchReviewsFromAPI(productID int) ([]Review, error) {
	resp, err := http.Get(fmt.Sprintf(ServerAddress+"reviews/product/%d", productID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reviewsResponse ReviewsResponse
	if err := json.Unmarshal(body, &reviewsResponse); err != nil {
		return nil, err
	}

	return reviewsResponse.Reviews, nil
}

// 发送评论
func sendReview(productID int, rating, comment string) {
	if currentUser == nil {
		fmt.Println("用户未登录")
		dialog.ShowError(fmt.Errorf("用户未登录\n请先登录"), fyne.CurrentApp().Driver().AllWindows()[1])
		return
	}
	review := Review{
		UserID:    currentUser.ID,
		ProductID: productID,
		Rating:    rating,
		Comment:   comment,
		Time:      time.Now(),
	}

	reviewJSON, err := json.Marshal(review)
	if err != nil {
		fmt.Println("Error encoding review:", err)
		return
	}

	resp, err := http.Post(ServerAddress+"reviews/", "application/json", bytes.NewBuffer(reviewJSON))
	if err != nil {
		fmt.Println("Error sending review:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("评论发送成功")
	} else {
		fmt.Println("评论发送失败")
	}
}
