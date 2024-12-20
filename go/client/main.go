package main

import (
	"client/pages"

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
	// 导航栏
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("主页", theme.HomeIcon(), pages.CreateHomePage()),
		container.NewTabItemWithIcon("购物车", theme.ContentAddIcon(), pages.CreateCartPage()),
		container.NewTabItemWithIcon("订单", theme.DocumentIcon(), pages.CreateOrdersPage()),
		container.NewTabItemWithIcon("个人主页", theme.AccountIcon(), pages.CreateProfilePage()),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	// 设置窗口内容
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(360, 780))
	// w.SetFixedSize(true)
	w.ShowAndRun()
}
