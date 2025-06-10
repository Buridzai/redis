package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/controllers"
)

// Sửa lại để nhận *gin.RouterGroup
func AuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", controllers.Register)
	rg.POST("/login", controllers.Login)
}
