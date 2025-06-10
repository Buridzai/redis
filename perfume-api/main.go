package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/config"
	"github.com/yourusername/perfume-api/models"
	"github.com/yourusername/perfume-api/routes"
	"github.com/yourusername/perfume-api/utils"
)

func main() {
	// 1) Kết nối database & migrate các bảng
	config.ConnectDatabase()
	config.ConnectRedis()
	config.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
	)

	// 2) Seed dữ liệu mẫu
	seedAdmin()
	seedProducts()

	// 3) Khởi tạo Gin
	r := gin.Default()

	// 4) Cấu hình CORS (BẮT BUỘC phải nằm trước khi khai báo route)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Hoặc []string{"*"} trong dev
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 5) Đăng ký các nhóm route
	api := r.Group("/api")
	{
		routes.AuthRoutes(api.Group("/auth"))
		routes.ProductRoutes(api.Group("/products")) // GET /api/products, GET /api/products/:id, v.v.
		routes.CartRoutes(api)                       // POST /api/cart, GET /api/cart, v.v.
		routes.OrderRoutes(api)                      // POST /api/orders, GET /api/orders, v.v.
	}
	// Đăng ký route admin (nếu có)
	routes.AdminProductRoutes(r) // /admin/products, /admin/products/:id, v.v.

	// 6) Chạy server ở port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Không thể chạy server: %v", err)
	}
}

func seedAdmin() {
	var count int64
	config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		admin := models.User{
			Name:     "Admin",
			Email:    "admin@gmail.com",
			Password: utils.HashPasswordNoErr("admin123"),
			Role:     "admin",
		}
		config.DB.Create(&admin)
		log.Println("✅ Đã tạo admin: admin@gmail.com / admin123")
	}
}

func seedProducts() {
	var count int64
	config.DB.Model(&models.Product{}).Count(&count)
	if count == 0 {
		products := []models.Product{
			{
				Name:        "Chanel No.5",
				Price:       2900000,
				Description: "Hương thơm kinh điển, quyến rũ & sang trọng bậc nhất thế giới.",
				Image:       "https://mir-s3-cdn-cf.behance.net/project_modules/1400_opt_1/af48ac95889975.5ea4495e81a85.jpg",
			},
			{
				Name:        "Dior Sauvage",
				Price:       2500000,
				Description: "Năng động, cá tính, cuốn hút – biểu tượng của nam tính hiện đại.",
				Image:       "https://blog.atome.id/wp-content/uploads/2022/06/Cek-harga-Dior-Sauvage-parfum-pria-dengan-wangi-maskulin.jpg",
			},
			{
				Name:        "Gucci Bloom",
				Price:       2200000,
				Description: "Hoa nhài và hoa huệ trắng thanh tao, ngọt dịu và nữ tính.",
				Image:       "https://th.bing.com/th/id/OIP.HMPSCQBUKsiJiWduRzk2MwAAAA?w=450&h=450&rs=1&pid=ImgDetMain",
			},
		}
		for _, p := range products {
			config.DB.Create(&p)
		}
		log.Println("✅ Đã seed sản phẩm mẫu vào bảng `products`")
	}
}
