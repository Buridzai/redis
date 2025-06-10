package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/controllers"
	"github.com/yourusername/perfume-api/middlewares"
)

func OrderRoutes(api *gin.RouterGroup) {
	orders := api.Group("/orders")
	{
		// Áp middleware xác thực JWT cho tất cả route trong nhóm /api/orders
		orders.Use(middlewares.JWTAuthMiddleware())

		// POST /api/orders     -> controllers.CreateOrder
		orders.POST("", controllers.CreateOrder)

		// GET  /api/orders     -> controllers.GetOrders (lấy lịch sử đơn hàng user)
		orders.GET("", controllers.GetOrders)
	}
}
