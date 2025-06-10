package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yourusername/perfume-api/config"
	"github.com/yourusername/perfume-api/models"
	"github.com/yourusername/perfume-api/utils"
)

// SessionData là struct dùng cho Redis
type SessionData struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy header Authorization
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Thiếu hoặc sai header Bearer"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Giải mã JWT
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
			c.Abort()
			return
		}

		ctx := context.Background()
		key := "session:" + tokenString

		// 3. Thử lấy từ Redis
		raw, err := config.RedisClient.Get(ctx, key).Result()
		if err == nil {
			var session SessionData
			if jsonErr := json.Unmarshal([]byte(raw), &session); jsonErr == nil {
				// Refresh TTL
				config.RedisClient.Expire(ctx, key, 72*time.Hour)

				c.Set("user", session)
				c.Set("user_id", session.ID)
				c.Set("role", session.Role)
				c.Set("permissions", session.Permissions)
				c.Next()
				return
			}
		} else if err != redis.Nil {
			log.Printf("❌ Redis error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}

		// 4. Fallback nếu không có trong Redis
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token thiếu user_id"})
			c.Abort()
			return
		}

		var user models.User
		if err := config.DB.First(&user, uint(userIDFloat)).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User không tồn tại"})
			c.Abort()
			return
		}

		// 5. Gán vào context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("role", user.Role)
		c.Next()
	}
}
