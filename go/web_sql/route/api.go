package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func APIInit() {
	r := gin.Default()
	SetupRoutes(r)
	r.Run(":8080")
}
func SetupRoutes(r *gin.Engine) {
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", control.RegisterHandler)
		userRoutes.POST("/login", control.LoginHandler)
		userRoutes.GET("/profile/:id", control.GetUserInfoHandler)
		userRoutes.PUT("/profile", control.UpdateUserInfoHandler)
	}
	productRoutes := r.Group("/products")
	{
		productRoutes.GET("/:id", control.GetProductsHandler)
		// productRoutes.GET("/:id", control.GetProduct)
		// productRoutes.GET("/search", control.SearchProducts)
	}
	// orderRoutes := r.Group("/orders")
	// {
	// 	orderRoutes.POST("/", control.CreateOrder)
	// 	orderRoutes.GET("/:id", control.GetOrder)
	// }
	cartRoutes := r.Group("/cart")
	{
		cartRoutes.POST("/items", control.AddCartHandler)
		// cartRoutes.GET("/items", control.)
		cartRoutes.DELETE("/items/:id", control.RemoveCartHandler)
	}
}
