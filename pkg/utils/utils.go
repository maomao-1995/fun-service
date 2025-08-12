package utils

import (
	"fmt"
	"fun-service/pkg/redisMain"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

// 生成随机昵称：前缀 + 6位随机数字
func GenerateRandomNickname() string {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 生成6位随机数字（100000-999999之间）
	randomNum := 100000000 + rand.Intn(900000000)

	// 拼接前缀和随机数字
	return fmt.Sprintf("%s_%d", "用户", randomNum)
}

func CheckHash(hash string) (string, bool) {
	// 检查哈希值是否已存在
	filename, err := redisMain.Rdb.Get(redisMain.Ctx, hash).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		fmt.Println("查询哈希值失败:", err)
		return "", false
	}
	return filename, true
}

// 保存哈希值到 Redis
func SaveHash(hash, filename string) {
	err := redisMain.Rdb.Set(redisMain.Ctx, hash, filename, 0).Err()
	if err != nil {
		fmt.Println("保存哈希值失败:", err)
	}
}
