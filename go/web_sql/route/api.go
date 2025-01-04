package route

import (
	"fmt"
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func APIInit() {
	r := gin.Default()
	SetupRoutes(r)
	r.Run(":8080")

	fmt.Println("Web API Start at 8080 Successfully")
}
func SetupRoutes(r *gin.Engine) {
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", control.RegisterHandler)
		userRoutes.POST("/login", control.LoginHandler)
		userRoutes.GET("/profile/:id", control.GetUserInfoHandler)
		userRoutes.PUT("/profile", control.UpdateUserInfoHandler)
		userRoutes.PUT("/change-password", control.ChangePasswordHandler)
	}
	productRoutes := r.Group("/products")
	{
		productRoutes.GET("/", control.GetProductsHandler)
		productRoutes.GET("/:id", control.GetProductsHandler)
		productRoutes.POST("/search", control.SearchProductsHandler)
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
	orderRoutes := r.Group("/orders")
	{
		orderRoutes.POST("/checkout", control.CheckoutHandler)
		orderRoutes.POST("/checkout/result", control.CheckoutResultHandler)
		orderRoutes.GET("/checkout/result", control.GetCheckoutResultHandler)
		orderRoutes.GET("/:id", control.GetOrdersHandler)
	}
}
