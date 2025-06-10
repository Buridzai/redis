package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/config"
	"github.com/yourusername/perfume-api/models"
)

func GetCart(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var cart []models.Cart
	config.DB.Preload("Product").Where("user_id = ?", userID).Find(&cart)
	c.JSON(http.StatusOK, cart)
}

func AddToCart(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var input struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cart := models.Cart{
		UserID:    userID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}
	config.DB.Create(&cart)
	c.JSON(http.StatusOK, cart)
}

func UpdateCartItem(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Model(&models.Cart{}).Where("id = ?", id).Update("quantity", input.Quantity)
	c.JSON(http.StatusOK, gin.H{"message": "Đã cập nhật"})
}

func DeleteCartItem(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Cart{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Đã xoá khỏi giỏ"})
}
