package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupOrderItemRoutes(router *gin.Engine) {
	orderItemGroup := router.Group("/order-items")
	{
		orderItemGroup.GET("/:id", control.GetOrderItem)                      // 获取单个订单项
		orderItemGroup.GET("/", control.ListOrderItems)                       // 查询所有订单项（分页）
		orderItemGroup.POST("/", control.CreateOrderItem)                     // 创建新订单项
		orderItemGroup.PUT("/:id", control.UpdateOrderItem)                   // 更新订单项
		orderItemGroup.DELETE("/:id", control.DeleteOrderItem)                // 删除订单项
		orderItemGroup.GET("/order/:orderid", control.GetOrderItemsByOrderID) // 根据订单ID获取订单项
	}
}
