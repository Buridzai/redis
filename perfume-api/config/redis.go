package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	RedisCtx    = context.Background()
)

func ConnectRedis() {
	addr := os.Getenv("REDIS_ADDR") // <-- lấy từ biến môi trường
	if addr == "" {
		addr = "redis:6379" // fallback nếu biến không được đặt
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "", // nếu có password thì thêm vào đây
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	if err := RedisClient.Ping(RedisCtx).Err(); err != nil {
		panic("Không thể kết nối Redis: " + err.Error())
	}

	log.Println("✅ Đã kết nối Redis thành công!")
}
