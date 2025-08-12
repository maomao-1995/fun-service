package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
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

// 加载已上传文件的哈希值
// 全局哈希集合，用于存储已上传文件的哈希值和文件名
var UploadedHashes = make(map[string]string) // 哈希值映射到文件名
const HashFilePath = "./uploads/hashes.json"

func LoadHashes() {

	file, err := os.Open(HashFilePath)
	if err != nil {
		// 如果文件不存在，创建一个空文件
		file, err = os.Create(HashFilePath)
		if err != nil {
			fmt.Println("创建哈希文件失败:", err)
			return
		}
		defer file.Close()
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&UploadedHashes); err != nil {
		fmt.Println("加载哈希值失败:", err)
	}
}

// 保存哈希值到文件
func SaveHashes() {
	file, err := os.Create(HashFilePath)
	if err != nil {
		fmt.Println("保存哈希值失败:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(UploadedHashes); err != nil {
		fmt.Println("保存哈希值失败:", err)
	}
}
