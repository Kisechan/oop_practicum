package control

import (
	"encoding/json"
	"net/http"
	"time"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
)

func GetProductsHandler(c *gin.Context) {
	// 尝试从 Redis 缓存中获取商品信息
	// cacheKey := "random_products"
	// cachedProducts, err := redisClient.Get(ctx, cacheKey).Result()
	// if err == nil {
	// 	// 如果缓存中存在商品信息，直接返回
	// 	var products []rep.Product
	// 	if err := json.Unmarshal([]byte(cachedProducts), &products); err == nil {
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"products": products,
	// 		})
	// 		return
	// 	}
	// }

	// 随机商品
	var products []rep.Product
	if err := rep.DB.Order("RAND()").Limit(20).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// 将商品信息存入 Redis 缓存
	// productsJSON, err := json.Marshal(products)
	// if err == nil {
	// 	redisClient.Set(ctx, cacheKey, productsJSON, time.Hour) // 缓存 1 小时
	// }

	// 返回商品信息
	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func GetAllProductsHandler(c *gin.Context) {
	// 尝试从 Redis 缓存中获取所有商品信息
	cacheKey := "all_products"
	cachedProducts, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// 如果缓存中存在商品信息，直接返回
		var products []rep.Product
		if err := json.Unmarshal([]byte(cachedProducts), &products); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"products": products,
			})
			return
		}
	}

	// 查询所有商品
	var products []rep.Product
	if err := rep.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// 将商品信息存入 Redis 缓存
	productsJSON, err := json.Marshal(products)
	if err == nil {
		redisClient.Set(ctx, cacheKey, productsJSON, time.Hour) // 缓存 1 小时
	}

	// 返回商品信息
	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

// 根据关键词搜索商品接口
func SearchProductsHandler(c *gin.Context) {
	// 获取查询参数
	name := c.Query("name") // 产品名称
	// category := c.Query("category")  // 产品类别
	// minPrice := c.Query("min_price") // 最低价格
	// maxPrice := c.Query("max_price") // 最高价格
	// seller := c.Query("seller")      // 卖家
	// isActive := c.Query("is_active") // 是否激活

	// 构建查询条件
	var products []rep.Product
	query := rep.DB.Model(&rep.Product{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	// if category != "" {
	// 	query = query.Where("category = ?", category)
	// }
	// if minPrice != "" {
	// 	query = query.Where("price >= ?", minPrice)
	// }
	// if maxPrice != "" {
	// 	query = query.Where("price <= ?", maxPrice)
	// }
	// if seller != "" {
	// 	query = query.Where("seller = ?", seller)
	// }
	// if isActive != "" {
	// 	query = query.Where("is_active = ?", isActive)
	// }

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
