package control

import (
	"net/http"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
)

func CreateReviewHandler(c *gin.Context) {
	var review rep.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rep.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

func GetProductReviewsHandler(c *gin.Context) {
	productID := c.Param("id")

	var reviews []rep.Review
	if err := rep.DB.Preload("Users").Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
	})
}
