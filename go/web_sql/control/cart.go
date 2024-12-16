package control

import (
	"github.com/gin-gonic/gin"
)

// 加购物车接口
func AddCartHandler(c *gin.Context) {
	// username := c.PostForm("username")
	// productID := c.PostForm("product_id")

	// cartMutex.Lock()
	// defer cartMutex.Unlock()

	// carts[username] = append(carts[username], productID)
	// c.JSON(http.StatusOK, gin.H{
	// 	"status":  "success",
	// 	"message": "添加成功",
	// })
}

// 删除购物车接口
func RemoveCartHandler(c *gin.Context) {
	// username := c.PostForm("username")
	// productID := c.PostForm("product_id")

	// cartMutex.Lock()
	// defer cartMutex.Unlock()

	// for i, id := range carts[username] {
	// 	if id == productID {
	// 		carts[username] = append(carts[username][:i], carts[username][i+1:]...)
	// 		break
	// 	}
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"status":  "success",
	// 	"message": "删除成功",
	// })
}
