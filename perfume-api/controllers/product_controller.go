package controllers

import (
	"net/http"

	"github.com/yourusername/perfume-api/config"
	"github.com/yourusername/perfume-api/models"

	"github.com/gin-gonic/gin"
)

// @Summary Lấy tất cả sản phẩm
// @Produce json
// @Success 200 {array} models.Product
// @Router /api/products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product
	config.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}

// @Summary Tạo sản phẩm mới
// @Accept json
// @Produce json
// @Param product body models.Product true "Thông tin sản phẩm"
// @Success 200 {object} models.Product
// @Router /api/products [post]
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&product)
	c.JSON(http.StatusOK, product)
}

// @Summary Lấy sản phẩm theo ID
// @Produce json
// @Param id path int true "ID sản phẩm"
// @Success 200 {object} models.Product
// @Router /api/products/{id} [get]
func GetProductByID(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sản phẩm"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// @Summary Cập nhật sản phẩm
// @Accept json
// @Produce json
// @Param id path int true "ID sản phẩm"
// @Param product body models.Product true "Thông tin mới"
// @Success 200 {object} models.Product
// @Router /api/products/{id} [put]
func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sản phẩm"})
		return
	}
	var data models.Product
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Model(&product).Updates(data)
	c.JSON(http.StatusOK, product)
}

// @Summary Xóa sản phẩm
// @Produce json
// @Param id path int true "ID sản phẩm"
// @Success 200 {object} map[string]string
// @Router /api/products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sản phẩm"})
		return
	}
	config.DB.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Xóa sản phẩm thành công"})
}

func GetProductDetail(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy sản phẩm"})
		return
	}
	c.JSON(200, product)
}
