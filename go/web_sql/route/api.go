package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", control.GetAllUsers)
		userGroup.GET("/:id", control.GetUserByID)
		userGroup.POST("/", control.CreateUser)
		userGroup.PUT("/:id", control.UpdateUser)
		userGroup.DELETE("/:id", control.DeleteUser)
	}
	// cartGroup := r.Group("/carts")
	// {
	// 	cartGroup.GET("/", control.GetAllCarts)
	// 	cartGroup.GET("/:id", control.GetCartByID)
	// 	cartGroup.POST("/", control.CreateCart)
	// 	cartGroup.PUT("/:id", control.UpdateCart)
	// 	cartGroup.DELETE("/:id", control.DeleteCart)
	// }
}
