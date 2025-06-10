package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/controllers"
	"github.com/yourusername/perfume-api/middlewares"
)

func CartRoutes(r *gin.RouterGroup) {
	r.Use(middlewares.JWTAuthMiddleware())

	r.GET("/cart", controllers.GetCart)
	r.POST("/cart", controllers.AddToCart)
	r.PUT("/cart/:id", controllers.UpdateCartItem)
	r.DELETE("/cart/:id", controllers.DeleteCartItem)
}
