package main

import (
	"client/pages"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func main() {
	// 创建应用
	a := app.NewWithID("com.kisechan.shopclientnt")
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("网购客户端")

	// 加载图片资源
	cartIconResource, err := loadImageResource("image/cart.png")
	if err != nil {
		log.Fatalf("Failed to load image resource: %v", err)
	}
	if cartIconResource == nil {
		log.Fatalf("cartIconResource is nil")
	}
	fmt.Println("cartIconResources is not nil")
	// 导航栏
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("主页", theme.HomeIcon(), pages.CreateHomePage()),
		container.NewTabItemWithIcon("购物车", cartIconResource, pages.CreateCartPage()),
		container.NewTabItemWithIcon("订单", theme.DocumentIcon(), pages.CreateOrdersPage()),
		container.NewTabItemWithIcon("我的", theme.AccountIcon(), pages.CreateProfilePage()),
		container.NewTabItemWithIcon("设置", theme.SettingsIcon(), pages.CreateSettingsPage()),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	// 设置窗口内容
	w.SetContent(tabs)
	// 360, 780为一般的手机屏幕大小
	w.Resize(fyne.NewSize(400, 480))
	w.SetFixedSize(true)
	w.ShowAndRun()
}

// 加载图片资源
func loadImageResource(path string) (fyne.Resource, error) {
	// 读取图片文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 将图片数据转换为 Fyne 资源
	return fyne.NewStaticResource("cart.jpg", data), nil
}
