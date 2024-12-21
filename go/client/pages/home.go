package pages

import (
	"client/str"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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

// 主页面
func CreateHomePage() fyne.CanvasObject {
	// 搜索框
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索商品...")

	// 搜索按钮
	searchButton := widget.NewButtonWithIcon("搜索", theme.SearchIcon(), func() {
		// 这里可以添加搜索逻辑
	})

	// 商品列表
	productCards := createProductCardsFromAPI()

	// 创建可滑动的卡片群
	scrollableCards := container.NewVScroll(container.NewGridWithColumns(1, productCards...))

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
func createProductCardsFromAPI() []fyne.CanvasObject {
	// 调用API获取商品数据
	products, err := fetchProductsFromAPI()
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

	return cards
}

// 调用API获取商品数据
func fetchProductsFromAPI() ([]Product, error) {
	resp, err := http.Get("http://localhost:8080/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var productsResponse ProductsResponse
	if err := json.Unmarshal(body, &productsResponse); err != nil {
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

	// 创建详细页面内容
	detailContent := container.NewVBox(
		image,
		widget.NewLabel("商品名称: "+product.Name),
		widget.NewLabel(fmt.Sprintf("价格: ￥%.2f", product.Price)),
		widget.NewLabel("商家: "+product.Seller),
		widget.NewLabel("描述:"),
		descriptionLabel,
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
