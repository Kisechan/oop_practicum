package pages

import (
	"bytes"
	"client/str"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Seller      string  `json:"seller"`
	IsActive    string  `json:"is_active"`
	Icon        *string `json:"icon"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}

type Review struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Rating    string    `json:"rating"`
	Comment   string    `json:"comment"`
	Time      time.Time `json:"time"`
}

type ReviewsResponse struct {
	Reviews []Review `json:"reviews"`
}

// 主页面
func CreateHomePage() fyne.CanvasObject {
	// 搜索框
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索商品...")

	// 商品列表
	productCards := createProductCardsFromAPI("")
	// 创建可滑动的卡片群
	scrollableCards := container.NewVScroll(container.NewGridWithColumns(1, productCards...))

	// 搜索按钮
	searchButton := widget.NewButtonWithIcon("搜索", theme.SearchIcon(), func() {
		productCards = createProductCardsFromAPI(searchEntry.Text)
		scrollableCards.Content = container.NewGridWithColumns(1, productCards...)
		scrollableCards.Refresh() // 刷新卡片群
		fmt.Println("scrollableCards Refreshed Successfully")
	})

	// 布局：顶部固定搜索框，下方为可滑动的卡片群
	return container.NewBorder(
		container.NewBorder(nil, nil, nil, searchButton, searchEntry), // 顶部搜索框
		nil,             // 底部无内容
		nil,             // 左侧无内容
		nil,             // 右侧无内容
		scrollableCards, // 下方可滑动的卡片群
	)
}

// 从API获取商品数据并创建商品卡片
func createProductCardsFromAPI(search string) []fyne.CanvasObject {
	// 调用API获取商品数据
	products, err := fetchProductsFromAPI(search)
	if err != nil {
		fmt.Println("Error fetching products:", err)
		return nil
	}

	// 创建卡片
	var cards []fyne.CanvasObject
	for _, product := range products {
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
		image.SetMinSize(fyne.NewSize(100, 100))

		nameLabel := widget.NewRichText(
			&widget.TextSegment{
				Text: str.CutStr(product.Name, 17),
				Style: widget.RichTextStyle{
					SizeName: theme.SizeNameHeadingText,
					// Inline: true,
				},
			},
		)

		priceLabel := widget.NewRichText(
			&widget.TextSegment{
				Text: fmt.Sprintf("￥%.2f", product.Price),
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					SizeName:  theme.SizeNameHeadingText,
					Alignment: fyne.TextAlignTrailing,
					TextStyle: fyne.TextStyle{
						Bold:   true,
						Italic: true,
					},
				},
			},
		)
		// sellerLabel := widget.NewLabel(product.Seller)

		// 创建卡片内容
		cardContent := container.NewVBox(
			image,
			nameLabel,
			priceLabel,
			// sellerLabel,
			widget.NewButton("详情", func() {
				// 打开新窗口展示商品详细页面
				showProductDetailWindow(product)
			}),
		)

		// 创建卡片
		card := widget.NewCard("", "", cardContent)
		cards = append(cards, card)
	}
	fmt.Println("Cards Created Finished Successfully")
	return cards
}

// 调用API获取商品数据
func fetchProductsFromAPI(search string) ([]Product, error) {
	var resp *http.Response
	var err error
	if search != "" {
		searchRequest := map[string]string{
			"name": search,
		}
		searchJSON, err := json.Marshal(searchRequest)
		if err != nil {
			fmt.Println("Error encoding search request:", err)
			return nil, err
		}

		// 发送搜索请求
		resp, err = http.Post("http://localhost:8080/products/search", "application/json", bytes.NewBuffer(searchJSON))
		if err != nil {
			fmt.Println("Error get request:", err)
			return nil, err
		}
	} else {
		resp, err = http.Get("http://localhost:8080/products")
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

	// 打印 API 返回的原始数据
	// fmt.Println("API Response:", string(body))

	var productsResponse ProductsResponse
	if err := json.Unmarshal(body, &productsResponse); err != nil {
		fmt.Println("Error unmarshal reqbody:", err)
		return nil, err
	}

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

	// 创建评论列表
	reviewList := container.NewVBox()
	for _, review := range reviews {
		reviewItem := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("评分: %s", review.Rating)),
			widget.NewLabel(fmt.Sprintf("评论: %s", review.Comment)),
			widget.NewLabel(fmt.Sprintf("时间: %s", review.Time.Format("2006-01-02 15:04:05"))),
		)
		reviewList.Add(reviewItem)
	}

	// 创建评论输入框
	commentEntry := widget.NewEntry()
	commentEntry.SetPlaceHolder("输入评论...")
	ratingEntry := widget.NewEntry()
	ratingEntry.SetPlaceHolder("输入评分 (1-5)...")

	// 创建提交按钮
	submitButton := widget.NewButton("提交评论", func() {
		// 发送评论
		sendReview(product.ID, ratingEntry.Text, commentEntry.Text)
		// 刷新评论列表
		reviews, err := fetchReviewsFromAPI(product.ID)
		if err != nil {
			fmt.Println("Error fetching reviews:", err)
		} else {
			reviewList.Objects = nil // 清空评论列表
			for _, review := range reviews {
				reviewItem := container.NewVBox(
					widget.NewLabel(fmt.Sprintf("评分: %s", review.Rating)),
					widget.NewLabel(fmt.Sprintf("评论: %s", review.Comment)),
					widget.NewLabel(fmt.Sprintf("时间: %s", review.Time.Format("2006-01-02 15:04:05"))),
				)
				reviewList.Add(reviewItem)
			}
			reviewList.Refresh() // 刷新评论列表
		}
	})

	// 创建评论输入部分
	commentInput := container.NewVBox(
		widget.NewLabel("添加评论"),
		ratingEntry,
		commentEntry,
		submitButton,
	)

	// 创建详细页面内容
	detailContent := container.NewVBox(
		image,
		widget.NewLabel("商品名称: "+product.Name),
		widget.NewLabel(fmt.Sprintf("价格: ￥%.2f", product.Price)),
		widget.NewLabel("商家: "+product.Seller),
		widget.NewLabel("描述:"),
		descriptionLabel,
		widget.NewLabel("评论:"),
		reviewList,
		commentInput,
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
		}),
		widget.NewButton("购买", func() {
			fmt.Println("购买:", product.Name)
		}),
	)

	// 将详细内容和按钮组合在一起
	scrollableContent := container.NewVScroll(detailContent)
	return container.NewBorder(nil, buttons, nil, nil, scrollableContent)
}

// 获取评论
func fetchReviewsFromAPI(productID int) ([]Review, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/reviews/product/%d", productID))
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
	review := Review{
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

	resp, err := http.Post("http://localhost:8080/reviews/", "application/json", bytes.NewBuffer(reviewJSON))
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
