package control

import (
	"net/http"
	"strconv"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取用户购物车列表
func ListCartItems(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var cartItems []rep.Cart
	if err := rep.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	c.JSON(http.StatusOK, cartItems)
}

// 添加商品到购物车
func AddToCart(c *gin.Context) {
	var cartItem rep.Cart
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 检查商品是否存在
	var product rep.Product
	if err := rep.DB.First(&product, cartItem.ProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 检查库存是否足够
	if product.Stock < cartItem.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// 检查购物车中是否已存在该商品
	var existingItem rep.Cart
	if err := rep.DB.Where("user_id = ? AND product_id = ?", cartItem.UserID, cartItem.ProductID).First(&existingItem).Error; err == nil {
		// 如果存在，更新数量
		existingItem.Quantity += cartItem.Quantity
		if err := rep.DB.Save(&existingItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Cart updated successfully", "cart": existingItem})
		return
	}

	// 添加新记录到购物车
	if err := rep.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cartItem)
}

// 更新购物车中的商品
func UpdateCartItem(c *gin.Context) {
	cartID, err := strconv.Atoi(c.Param("cartid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	var cartItem rep.Cart
	if err := rep.DB.First(&cartItem, cartID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var updateData struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 检查库存是否足够
	var product rep.Product
	if err := rep.DB.First(&product, cartItem.ProductID).Error; err == nil {
		if product.Stock < updateData.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
	}

	cartItem.Quantity = updateData.Quantity
	if err := rep.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated successfully", "cart": cartItem})
}

// 删除购物车中的商品
func RemoveCartItem(c *gin.Context) {
	cartID, err := strconv.Atoi(c.Param("cartid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	if err := rep.DB.Delete(&rep.Cart{}, cartID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item removed successfully"})
}

// 清空用户的购物车
func ClearCart(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := rep.DB.Where("user_id = ?", userID).Delete(&rep.Cart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
