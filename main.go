package main

import (
	"fmt"
	"fun-service/config"
	_ "fun-service/docs"
	_ "fun-service/internal/handler"
	"fun-service/internal/router"
	"fun-service/pkg/database"
	"fun-service/pkg/logger"
	"fun-service/pkg/redis"
	"fun-service/pkg/utils"
)

// @title Fun Service API
// @version 1.0
// @host localhost:8800
// @BasePath

func main() {
	//加载配置
	cfg, cfgErr := config.Load()
	if cfgErr != nil {
		panic("加载配置失败: " + cfgErr.Error())
	} else {
		fmt.Println("配置加载成功")
	}

	//初始化日志
	logCfg := logger.Config{
		Level: cfg.Log.Level,
		Path:  cfg.Log.Path,
		// 添加其他字段映射（如果有的话）
	}
	logger.Init(logCfg)
	fmt.Println("日志系统初始化成功")

	//初始化数据库
	database.InitMySQL(cfg.MySQL)

	//连接redis
	redisErr := redis.InitRedis()
	if redisErr != nil {
		panic("Redis连接失败: " + redisErr.Error())
	} else {
		fmt.Println("Redis连接成功")
	}

	//加载已上传文件的哈希值
	utils.LoadHashes()

	// 启动服务
	routerInstance := router.SetupRouter()
	serverErr := routerInstance.Run(cfg.Server.Addr)
	if serverErr != nil {
		panic("启动服务失败: " + serverErr.Error())
	}
}
