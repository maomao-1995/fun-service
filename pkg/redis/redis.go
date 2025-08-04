package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// 全局Redis客户端实例
var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "",               // 密码，如果没有则为空
		DB:       0,                // 使用的数据库编号
		// 连接池设置
		PoolSize:        10,
		MinIdleConns:    5,
		MaxConnAge:      30 * time.Minute,
		IdleTimeout:     5 * time.Minute,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond,
		MaxRetryBackoff: 1 * time.Second,
	})

	// 测试连接
	_, err := Rdb.Ping(Ctx).Result()
	return err
}


