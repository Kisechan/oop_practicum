package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(router *gin.Engine) {
	productGroup := router.Group("/products")
	{
		productGroup.GET("/:id", control.GetProduct)                          // 获取单个商品
		productGroup.GET("/", control.ListProducts)                           // 查询所有商品（分页）
		productGroup.POST("/", control.CreateProduct)                         // 创建新商品
		productGroup.PUT("/:id", control.UpdateProduct)                       // 更新商品
		productGroup.DELETE("/:id", control.DeleteProduct)                    // 删除商品
		productGroup.GET("/category/:categoryid", control.ProductsByCategory) // 获取某分类下的商品
		productGroup.GET("/active", control.ListActiveProducts)               // 查询所有上架商品
	}
}
