package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/config"
	"github.com/yourusername/perfume-api/models"
)

// Dùng chung kiểu input cho CreateOrder
type OrderItemInput struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type OrderRequest struct {
	Items []OrderItemInput `json:"items"`
}

// CreateOrder: tạo đơn hàng mới cho user hiện tại
func CreateOrder(c *gin.Context) {
	// Lấy user_id từ context do JWTAuthMiddleware gán vào
	userID := c.MustGet("user_id").(uint)

	var input OrderRequest
	if err := c.ShouldBindJSON(&input); err != nil || len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	var total float64
	var orderItems []models.OrderItem

	// Xử lý từng item trong request
	for _, item := range input.Items {
		var product models.Product
		if err := config.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sản phẩm không tồn tại"})
			return
		}

		// Nếu bạn có dùng product.Stock, kiểm tra ở đây (ví dụ:)
		// if item.Quantity > product.Stock {
		//     c.JSON(http.StatusBadRequest, gin.H{"error": "Không đủ hàng"})
		//     return
		// }

		// Cộng tiền (ép kiểu float64)
		total += product.Price * float64(item.Quantity)

		// Tạo slice orderItems dùng struct đã định nghĩa trong models
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price, // đã là float64
		})

		// Nếu có giảm tồn kho, uncomment:
		// config.DB.Model(&product).Update("stock", product.Stock-item.Quantity)
	}

	// Tạo đối tượng Order (models.Order.Total là float64)
	order := models.Order{
		UserID: userID,
		Total:  total,
		Items:  orderItems,
	}

	// Lưu Order + các OrderItem liên quan
	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo đơn hàng"})
		return
	}

	// Xóa toàn bộ Cart của user sau khi đặt
	config.DB.Where("user_id = ?", userID).Delete(&models.Cart{})

	c.JSON(http.StatusOK, gin.H{"message": "Đặt hàng thành công", "order": order})
}

// GetOrders: trả về mảng Order của user hiện tại (lịch sử đơn hàng)
func GetOrders(c *gin.Context) {
	// Lấy user_id từ context
	userID := c.MustGet("user_id").(uint)

	var orders []models.Order
	// Preload("Items") để tự động tải mảng OrderItem mỗi Order
	err := config.DB.
		Preload("Items").
		Where("user_id = ?", userID).
		Find(&orders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy lịch sử đơn hàng"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
