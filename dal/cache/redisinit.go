package cache

import (
	"context"
	"time"

	"github.com/WoWBytePaladin/go-mall/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func Redis() *redis.Client {
	return redisClient
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:         config.Redis.Addr,
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolTimeout:  30 * time.Second,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		// 连接不上redis 让项目停止启动
		panic(err)
	}
}
