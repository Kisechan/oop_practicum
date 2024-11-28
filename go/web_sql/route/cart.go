package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupCartRoutes(router *gin.Engine) {
	cartGroup := router.Group("/carts")
	{
		cartGroup.GET("/:userid", control.ListCartItems)      // 获取用户的购物车列表
		cartGroup.POST("/", control.AddToCart)                // 添加商品到购物车
		cartGroup.PUT("/:cartid", control.UpdateCartItem)     // 更新购物车中的商品
		cartGroup.DELETE("/:cartid", control.RemoveCartItem)  // 删除购物车中的商品
		cartGroup.DELETE("/clear/:userid", control.ClearCart) // 清空用户的购物车
	}
}
