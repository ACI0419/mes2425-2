package configs

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// GetDefaultRedisConfig 获取默认Redis配置
func GetDefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

// InitRedis 初始化Redis连接
func InitRedis(config *RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Redis connected successfully")
	return rdb, nil
}