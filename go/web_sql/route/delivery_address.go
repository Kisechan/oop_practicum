package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupDeliveryAddressRoutes(router *gin.Engine) {
	addressGroup := router.Group("/delivery-addresses")
	{
		addressGroup.GET("/:id", control.GetDeliveryAddress)            // 获取单个地址
		addressGroup.GET("/", control.ListDeliveryAddresses)            // 查询所有地址（分页）
		addressGroup.POST("/", control.CreateDeliveryAddress)           // 创建新地址
		addressGroup.PUT("/:id", control.UpdateDeliveryAddress)         // 更新地址
		addressGroup.DELETE("/:id", control.DeleteDeliveryAddress)      // 删除地址
		addressGroup.GET("/user/:userid", control.GetAddressesByUserID) // 根据用户ID查询地址
	}
}
