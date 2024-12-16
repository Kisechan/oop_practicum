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
		productRoutes.GET("/", control.GetAllProductsHandler)
		productRoutes.GET("/:id", control.GetProductsHandler)
		productRoutes.GET("/search", control.SearchProductsHandler)
	}
	cartRoutes := r.Group("/cart")
	{
		cartRoutes.POST("/items", control.AddCartHandler)
		cartRoutes.DELETE("/items/:id", control.RemoveCartHandler)
		cartRoutes.GET("/items", control.GetCartHandler)
	}
	reviewRoutes := r.Group("/reviews")
	{
		reviewRoutes.POST("/", control.CreateReviewHandler)
		reviewRoutes.GET("/product/:id", control.GetProductReviewsHandler)
	}
	couponRoutes := r.Group("/coupons")
	{
		couponRoutes.GET("/user/:id", control.GetUserCouponsHandler)
		couponRoutes.POST("/use", control.UseCouponHandler)
	}
}
