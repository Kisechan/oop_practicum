package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("/:id", control.GetUser)
		userGroup.POST("/", control.CreateUser)
		userGroup.PUT("/:id", control.UpdateUser)
		userGroup.DELETE("/:id", control.DeleteUser)

		userGroup.GET("/", control.ListUsers)              // 查询用户列表
		userGroup.DELETE("/", control.DeleteUsers)         // 批量删除用户
		userGroup.GET("/export", control.ExportUsers)      // 导出用户数据
		userGroup.GET("/paginate", control.PaginatedUsers) // 分页查询用户
	}
}
