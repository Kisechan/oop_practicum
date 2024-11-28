package route

import (
	"web_sql/control"

	"github.com/gin-gonic/gin"
)

func SetupCouponRoutes(router *gin.Engine) {
	couponGroup := router.Group("/coupons")
	{
		couponGroup.GET("/:id", control.GetCoupon)                 // 获取单个优惠券
		couponGroup.GET("/", control.ListCoupons)                  // 查询所有优惠券（分页）
		couponGroup.POST("/", control.CreateCoupon)                // 创建新优惠券
		couponGroup.PUT("/:id", control.UpdateCoupon)              // 更新优惠券
		couponGroup.DELETE("/:id", control.DeleteCoupon)           // 删除优惠券
		couponGroup.GET("/user/:userid", control.GetCouponsByUser) // 查询某用户的优惠券
		couponGroup.GET("/check/:code", control.CheckCouponByCode) // 根据优惠码检查优惠券
	}
}
