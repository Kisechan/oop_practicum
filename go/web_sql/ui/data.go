package ui

import (
	"strings"

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
			usersScreen,
			// welcomeScreen,
		},
		// "shippings": {"物流信息表",
		// 	"储存物流信息",
		// 	// shippingsScreen,
		// 	welcomeScreen,
		// },
		"reviews": {"评论表",
			"储存评论信息",
			reviewsScreen,
		},
		"products": {"商品表",
			"储存商品信息",
			productsScreen,
		},
		"orders": {"订单表",
			"储存订单信息",
			ordersScreen,
		},
		// "order_items": {"单项商品表",
		// 	"储存订单中的单项商品",
		// 	// order_itemsScreen,
		// 	welcomeScreen,
		// },
		// "delivery_addresses": {"收货地址表",
		// 	"储存用户的收货地址",
		// 	// delivery_addressesScreen,
		// 	welcomeScreen,
		// },
		"coupons": {"优惠券表",
			"储存优惠券",
			couponsScreen,
		},
		// "categories": {"商品种类表",
		// 	"储存商品种类信息",
		// 	// categoriesScreen,
		// 	welcomeScreen,
		// },
		"carts": {"购物车表",
			"储存购物车信息",
			cartsScreen,
		},
	}

	TablesIndex = map[string][]string{
		"":         {"welcome", "users", "products"},
		"users":    {"orders", "coupons", "carts"},
		"products": {"reviews"},
	}
)

func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, char := range str {
		if char >= 'A' && char <= 'Z' {
			// 如果是大写字母，且不是第一个字符，添加下划线
			if i > 0 {
				result.WriteRune('_')
			}
			// 将大写字母转换为小写
			result.WriteRune(char + 32)
		} else {
			// 其他字符直接写入
			result.WriteRune(char)
		}
	}
	return result.String()
}
