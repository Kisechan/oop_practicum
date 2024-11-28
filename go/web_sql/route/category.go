package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(router *gin.Engine) {
	categoryGroup := router.Group("/categories")
	{
		categoryGroup.GET("/:id", control.GetCategory)                   // 获取单个分类
		categoryGroup.GET("/", control.ListCategories)                   // 查询所有分类
		categoryGroup.POST("/", control.CreateCategory)                  // 创建新分类
		categoryGroup.PUT("/:id", control.UpdateCategory)                // 更新分类
		categoryGroup.DELETE("/:id", control.DeleteCategory)             // 删除分类
		categoryGroup.GET("/parent/:parentid", control.GetSubcategories) // 获取子分类
	}
}
