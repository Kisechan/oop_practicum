package control

import (
	"net/http"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
)

// 查看商品接口
func GetProductsHandler(c *gin.Context) {
	// 获取路径参数 id
	productID := c.Param("id")

	// 查询产品
	var product rep.Product
	if err := rep.DB.Preload("Reviews").First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 返回产品信息
	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func GetAllProductsHandler(c *gin.Context) {
	var products []rep.Product
	if err := rep.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

// 根据关键词搜索商品接口
func SearchProductsHandler(c *gin.Context) {
	// 获取查询参数
	name := c.Query("name")          // 产品名称
	category := c.Query("category")  // 产品类别
	minPrice := c.Query("min_price") // 最低价格
	maxPrice := c.Query("max_price") // 最高价格
	seller := c.Query("seller")      // 卖家
	isActive := c.Query("is_active") // 是否激活

	// 构建查询条件
	var products []rep.Product
	query := rep.DB.Model(&rep.Product{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}
	if seller != "" {
		query = query.Where("seller = ?", seller)
	}
	if isActive != "" {
		query = query.Where("is_active = ?", isActive)
	}

	// 执行查询
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
