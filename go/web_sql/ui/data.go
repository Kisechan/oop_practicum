package ui

import (
	"web_sql/rep"

	"fyne.io/fyne/v2"
)

type Table struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	// Tables defines the metadata for each Table
	Tables = map[string]Table{
		"welcome": {"欢迎",
			"欢迎使用Shop - APP后台管理系统",
			welcomeScreen,
		},
		"users": {"用户表",
			"存储用户信息",
			tableScreen[rep.User],
			// welcomeScreen,
		},
		"shippings": {"物流信息表",
			"储存物流信息",
			// shippingsScreen,
			welcomeScreen,
		},
		"reviews": {"评论表",
			"储存评论信息",
			// reviewsScreen,
			welcomeScreen,
		},
		"products": {"商品表",
			"储存商品信息",
			// productsScreen,
			welcomeScreen,
		},
		"orders": {"订单表",
			"储存订单信息",
			// ordersScreen,
			welcomeScreen,
		},
		"order_items": {"单项商品表",
			"储存订单中的单项商品",
			// order_itemsScreen,
			welcomeScreen,
		},
		"delivery_addresses": {"收货地址表",
			"储存用户的收货地址",
			// delivery_addressesScreen,
			welcomeScreen,
		},
		"coupons": {"优惠券表",
			"储存优惠券",
			// couponsScreen,
			welcomeScreen,
		},
		"categories": {"商品种类表",
			"储存商品种类信息",
			// categoriesScreen,
			welcomeScreen,
		},
		"carts": {"购物车表",
			"储存购物车信息",
			// cartsScreen,
			welcomeScreen,
		},
	}

	TablesIndex = map[string][]string{
		"": {"welcome", "users", "shippings", "reviews", "products", "orders", "order_items", "delivery_addresses", "coupons", "categories", "carts"},
		// "users": {"shippings_users","orders_users","delivery_addresses_users","coupons_users","carts_users"},
		// "orders":  {"detailed_orders"},
		// "widgets":     {"accordion", "activity", "button", "card", "entry", "form", "input", "progress", "text", "toolbar"},
	}
)
