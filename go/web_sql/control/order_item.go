package control

import (
	"net/http"
	"strconv"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取单个订单项
func GetOrderItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order item ID"})
		return
	}

	var orderItem rep.OrderItem
	if err := rep.DB.Preload("Order").Preload("Product").First(&orderItem, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, orderItem)
}

// 查询所有订单项（分页）
func ListOrderItems(c *gin.Context) {
	var orderItems []rep.OrderItem
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	if err := rep.DB.Offset(offset).Limit(pageSize).Preload("Order").Preload("Product").Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
		return
	}

	c.JSON(http.StatusOK, orderItems)
}

// 创建新订单项
func CreateOrderItem(c *gin.Context) {
	var orderItem rep.OrderItem
	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 计算总价 = 数量 * 单价
	orderItem.TotalPrice = float64(orderItem.Quantity) * orderItem.UnitPrice

	if err := rep.DB.Create(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderItem)
}

// 更新订单项
func UpdateOrderItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order item ID"})
		return
	}

	var orderItem rep.OrderItem
	if err := rep.DB.First(&orderItem, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var updateData struct {
		Quantity  int     `json:"quantity" binding:"required"`
		UnitPrice float64 `json:"unit_price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 更新字段并重新计算总价
	orderItem.Quantity = updateData.Quantity
	orderItem.UnitPrice = updateData.UnitPrice
	orderItem.TotalPrice = float64(updateData.Quantity) * updateData.UnitPrice

	if err := rep.DB.Save(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order item updated successfully", "order_item": orderItem})
}

// 删除订单项
func DeleteOrderItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order item ID"})
		return
	}

	if err := rep.DB.Delete(&rep.OrderItem{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order item deleted successfully"})
}

// 根据订单ID获取订单项
func GetOrderItemsByOrderID(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("orderid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var orderItems []rep.OrderItem
	if err := rep.DB.Where("order_id = ?", orderID).Preload("Product").Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
		return
	}

	c.JSON(http.StatusOK, orderItems)
}
