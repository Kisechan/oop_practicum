package control

import (
	"fmt"
	"net/http"
	"strconv"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) {
	var loginRequest struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	// 绑定 JSON 数据到 loginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user rep.User
	if err := rep.DB.Preload("Carts").Preload("Orders").Preload("Coupons").Where("phone = ?", loginRequest.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Phone number invalid"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
		return
	}

	fmt.Println("User:", user, "Login Successfully")
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// 注册接口
func RegisterHandler(c *gin.Context) {
	var user rep.User

	// 绑定 JSON 数据到 user 结构体
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hash Password Failed"})
		return
	}
	user.Password = string(hashedPassword)

	// 创建用户
	if err := rep.Create[rep.User](rep.DB, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Register Failed!"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// 查看个人信息接口
func GetUserInfoHandler(c *gin.Context) {
	// 获取路径参数 id
	userID, _ := strconv.Atoi(c.Param("id"))

	// 查询用户
	var user rep.User
	if err := rep.DB.Preload("Carts").Preload("Orders").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// 修改个人信息接口
func UpdateUserInfoHandler(c *gin.Context) {
	var updateRequest struct {
		ID       int     `json:"id"`
		Username string  `json:"username"`
		Email    *string `json:"email"`
		Phone    string  `json:"phone"`
		Address  *string `json:"address"`
	}

	// 绑定 JSON 数据到 updateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user rep.User
	if err := rep.DB.First(&user, updateRequest.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 更新用户信息
	user.Username = updateRequest.Username
	user.Email = updateRequest.Email
	user.Phone = updateRequest.Phone
	user.Address = updateRequest.Address

	if err := rep.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func ChangePasswordHandler(c *gin.Context) {
	var changePasswordRequest struct {
		UserID      int    `json:"user_id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&changePasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user rep.User
	if err := rep.DB.First(&user, changePasswordRequest.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePasswordRequest.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	// 对新密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}
	user.Password = string(hashedPassword)

	// 更新用户密码
	if err := rep.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}
