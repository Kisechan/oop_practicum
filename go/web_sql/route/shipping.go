package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupShippingRoutes(router *gin.Engine) {
	shippingGroup := router.Group("/shipping")
	{
		shippingGroup.GET("/:id", control.GetShipping)                                  // 获取单个运单
		shippingGroup.GET("/", control.ListShipping)                                    // 查询所有运单（分页）
		shippingGroup.POST("/", control.CreateShipping)                                 // 创建新运单
		shippingGroup.PUT("/:id", control.UpdateShipping)                               // 更新运单信息
		shippingGroup.DELETE("/:id", control.DeleteShipping)                            // 删除单个运单
		shippingGroup.GET("/track/:trackingNumber", control.TrackShipping)              // 根据追踪号查询运单状态
		shippingGroup.GET("/byOrderItem/:orderItemID", control.ListShippingByOrderItem) // 查询某订单项的所有运单
	}
}
