package utils

import (
	"fmt"
	"math/rand"
	"time"
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
