package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/config"
	"github.com/yourusername/perfume-api/models"
	"github.com/yourusername/perfume-api/utils"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra email đã tồn tại trước khi tạo user mới
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email đã tồn tại"})
		return
	}

	hashed, _ := utils.HashPassword(input.Password)
	user := models.User{Name: input.Name, Email: input.Email, Password: hashed}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // trả về lỗi thực tế
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Đăng ký thành công"})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không đúng"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không đúng"})
		return
	}

	token, _ := utils.GenerateJWT(user.ID, user.Email, user.Role)

	// Sau khi tạo token, lưu thông tin session vào Redis
	session := map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
		"role": user.Role,
		"permissions": []string{
			"product:view", "order:create", // assign quyền
		},
	}

	raw, _ := json.Marshal(session)
	key := "session:" + token
	config.RedisClient.Set(c, key, raw, 72*time.Hour)

	// Trả về token + thông tin user
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
