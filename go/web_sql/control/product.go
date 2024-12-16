package control

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 查看商品接口
func GetProductsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		// "products": products,
	})
}

// 根据关键词搜索商品接口
func SearchProductsHandler(c *gin.Context) {
	// keyword := c.Query("keyword")
	// var results []map[string]string

	// for _, product := range products {
	// 	if product["name"] == keyword {
	// 		results = append(results, product)
	// 	}
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"status":  "success",
	// 	"results": results,
	// })
}
