package rep

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createOrderHandler(c *gin.Context) {
	var order Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully"})
}

// 更新订单接口
func updateOrderHandler(c *gin.Context) {
	var order Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := DB.Model(&Order{}).Where("id = ?", order.ID).Updates(order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

func RepAPIInit() {
	r := gin.Default()

	r.POST("/api/orders/create", createOrderHandler)
	r.POST("/api/orders/update", updateOrderHandler)

	r.Run(":8081")

	fmt.Println("Rep API Start at 8081 Successfully")
}
