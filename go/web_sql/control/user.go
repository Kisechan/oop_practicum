// package control

// import (
// 	"net/http"
// 	"strconv"
// 	"web_sql/rep"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// // 获取用户信息
// func GetUser(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
// 		return
// 	}

// 	// 查询数据库
// 	user, err := rep.GetID[rep.User](rep.DB, id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// 返回响应
// 	c.JSON(http.StatusOK, user)
// }

// // 获取所有用户（分页）
// func ListUsers(c *gin.Context) {
// 	var users []rep.User
// 	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	if err != nil {
// 		page = 1
// 	}
// 	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
// 	if err != nil {
// 		pageSize = 10
// 	}

// 	offset := (page - 1) * pageSize
// 	if err := rep.DB.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, users)
// }

// // 创建新用户
// func CreateUser(c *gin.Context) {
// 	var user rep.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	if err := rep.DB.Create(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, user)
// }

// // 更新用户信息
// func UpdateUser(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
// 		return
// 	}

// 	var updateData rep.User
// 	if err := c.ShouldBindJSON(&updateData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	// 查询用户
// 	user, err := rep.GetID[rep.User](rep.DB, id)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		}
// 		return
// 	}

// 	// 更新数据
// 	if err := rep.DB.Model(&user).Updates(updateData).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
// }

// // 删除用户
// func DeleteUser(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
// 		return
// 	}

// 	_, err = rep.GetID[rep.User](rep.DB, id)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		}
// 		return
// 	}

// 	// 删除用户
// 	if err := rep.DeleteID[rep.User](rep.DB, id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
// }

// // 批量删除用户
// func DeleteUsers(c *gin.Context) {
// 	var ids []int
// 	if err := c.ShouldBindJSON(&ids); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	if len(ids) == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "No user IDs provided"})
// 		return
// 	}

// 	if err := rep.DB.Delete(&rep.User{}, ids).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Users deleted successfully"})
// }

// // 导出用户数据
// func ExportUsers(c *gin.Context) {
// 	var users []rep.User
// 	if err := rep.DB.Find(&users).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "导出用户数据失败"})
// 		return
// 	}

// 	// 假设你有一个导出功能，这里可以将用户数据转换成 CSV 或其它格式
// 	// 此处仅返回数据作为示例
// 	c.JSON(http.StatusOK, gin.H{"data": users})
// }

// // 分页查询用户
// func PaginatedUsers(c *gin.Context) {
// 	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	if err != nil {
// 		page = 1
// 	}
// 	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
// 	if err != nil {
// 		pageSize = 10
// 	}

// 	var users []rep.User
// 	offset := (page - 1) * pageSize
// 	if err := rep.DB.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "分页查询失败"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": users, "page": page, "pageSize": pageSize})
// }

package control

import (
	"net/http"
	"strconv"
	"web_sql/utils"

	"github.com/gin-gonic/gin"
)

// 获取所有用户
func GetAllUsers(c *gin.Context) {
	// 调用 C++ 程序的接口
	users, err := utils.CallCppAPI("getAllUsers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// 获取单个用户
func GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	// 调用 C++ 程序的接口
	user, err := utils.CallCppAPI("getUserByID", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// 创建用户
func CreateUser(c *gin.Context) {
	var user map[string]interface{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 调用 C++ 程序的接口
	result, err := utils.CallCppAPI("createUser", user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

// 更新用户
func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user map[string]interface{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 调用 C++ 程序的接口
	result, err := utils.CallCppAPI("updateUser", id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// 删除用户
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	// 调用 C++ 程序的接口
	_, err := utils.CallCppAPI("deleteUser", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}
