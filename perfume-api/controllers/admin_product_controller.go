package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/config"
	"github.com/yourusername/perfume-api/models"
)

// Tạo sản phẩm mới
func AdminCreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo sản phẩm"})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// Cập nhật sản phẩm
func AdminUpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sản phẩm"})
		return
	}
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Model(&product).Updates(input)
	c.JSON(http.StatusOK, product)
}

// Lấy danh sách tất cả sản phẩm
func AdminGetAllProducts(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách sản phẩm"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// Xóa sản phẩm
func AdminDeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Xóa thất bại"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}
