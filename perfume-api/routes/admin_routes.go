package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/controllers"
	"github.com/yourusername/perfume-api/middlewares"
)

func AdminProductRoutes(r *gin.Engine) {
	admin := r.Group("/admin", middlewares.RequireAdmin())
	{
		admin.POST("/products", controllers.AdminCreateProduct)
		admin.GET("/products", controllers.AdminGetAllProducts)
		admin.GET("/users", controllers.GetAllUsers)
		admin.PUT("/products/:id", controllers.AdminUpdateProduct)
		admin.DELETE("/products/:id", controllers.AdminDeleteProduct)
		// ...các route admin khác nếu cần
	}
}
