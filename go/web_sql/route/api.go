package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// userRoutes := r.Group("/users")
	{
		// userRoutes.POST("/register", control.Register)
		// userRoutes.POST("/login", control.Login)
		// userRoutes.GET("/profile", control.GetProfile)
		// userRoutes.PUT("/profile", control.UpdateProfile)
	}
	productRoutes := r.Group("/products")
	{
		productRoutes.GET("/", control.ListProducts)
		// productRoutes.GET("/:id", control.GetProduct)
		// productRoutes.GET("/search", control.SearchProducts)
	}
	orderRoutes := r.Group("/orders")
	{
		orderRoutes.POST("/", control.CreateOrder)
		orderRoutes.GET("/:id", control.GetOrder)
		// orderRoutes.PUT("/:id/cancel", control.CancelOrder)
	}
	// cartRoutes := r.Group("/cart")
	{
		// cartRoutes.POST("/items", cartController.AddItem)
		// cartRoutes.GET("/items", cartController.GetCartItems)
		// cartRoutes.DELETE("/items/:id", cartController.RemoveItem)
	}
}
