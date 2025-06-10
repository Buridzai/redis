package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/yourusername/perfume-api/config"
	"github.com/yourusername/perfume-api/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) Lấy header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2) Validate JWT và lấy claims
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userIDFloat, ok := claims["UserID"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		userID := uint(userIDFloat)

		// 4) Lấy session data từ Redis
		key := "session:" + tokenString
		raw, err := config.RedisClient.Get(context.Background(), key).Result()
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			return
		}

		var sess SessionData
		if err := json.Unmarshal([]byte(raw), &sess); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
			return
		}

		// 5) Đưa vào Gin context để downstream dùng
		c.Set("user_id", userID)
		c.Set("user_role", sess.Role)
		c.Set("user_perms", sess.Permissions)

		c.Next()
	}
}
