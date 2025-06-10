package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/controllers"
	"github.com/yourusername/perfume-api/middlewares"
)

func ProductRoutes(rg *gin.RouterGroup) {
	// Tất cả phải qua AuthMiddleware
	rg.Use(middlewares.AuthMiddleware())

	rg.GET("/", controllers.GetProducts) // ai cũng xem được?
	rg.POST("/", middlewares.Authorize("create:product"), controllers.CreateProduct)
	rg.PUT("/:id", middlewares.Authorize("update:product"), controllers.UpdateProduct)
	rg.DELETE("/:id", middlewares.Authorize("delete:product"), controllers.DeleteProduct)
	rg.GET("/:id", controllers.GetProductByID)
}
