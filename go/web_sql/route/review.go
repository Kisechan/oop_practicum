package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupReviewRoutes(router *gin.Engine) {
	reviewGroup := router.Group("/reviews")
	{
		reviewGroup.GET("/:id", control.GetReview)                       // 获取单个评论
		reviewGroup.GET("/", control.ListReviews)                        // 查询所有评论（分页）
		reviewGroup.POST("/", control.CreateReview)                      // 创建新评论
		reviewGroup.PUT("/:id", control.UpdateReview)                    // 更新评论
		reviewGroup.DELETE("/:id", control.DeleteReview)                 // 删除评论
		reviewGroup.GET("/product/:productid", control.ReviewsByProduct) // 获取某产品的评论
		reviewGroup.GET("/user/:userid", control.ReviewsByUser)          // 获取某用户的评论
	}
}
