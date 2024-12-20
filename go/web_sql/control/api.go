package control

import (
	"fmt"
	"net/http"
	"time"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
)

func createOrderHandler(c *gin.Context) {
	var req OrderRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	order := rep.Order{
		UserID:      req.UserID,
		Total:       &req.Total,
		Status:      "pending",
		CreatedTime: time.Now(),
		UpdateTime:  time.Now(),
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
	}
	if err := rep.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully"})
}

// 更新订单接口
// func updateOrderHandler(c *gin.Context) {
// 	var order Order
// 	if err := c.BindJSON(&order); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	if err := DB.Model(&Order{}).Where("id = ?", order.ID).Updates(order).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
// }

func RepAPIInit() {
	r := gin.Default()

	r.POST("/api/orders/create", createOrderHandler)
	// r.POST("/api/orders/update", updateOrderHandler)

	r.Run(":8081")

	fmt.Println("Rep API Start at 8081 Successfully")
}
