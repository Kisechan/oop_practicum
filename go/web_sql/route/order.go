package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.Engine) {
	orderGroup := router.Group("/orders")
	{
		orderGroup.GET("/:id", control.GetOrder)               // 获取单个订单
		orderGroup.GET("/", control.ListOrders)                // 查询所有订单（分页）
		orderGroup.POST("/", control.CreateOrder)              // 创建新订单
		orderGroup.PUT("/:id", control.UpdateOrder)            // 更新订单状态
		orderGroup.DELETE("/:id", control.DeleteOrder)         // 删除订单
		orderGroup.GET("/user/:userid", control.GetUserOrders) // 获取某用户的所有订单
	}
}
