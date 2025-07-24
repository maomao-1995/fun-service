package main

import (
	"fmt"
	"fun-service/config"
	"fun-service/internal/api/router"
	"fun-service/pkg/database"
	"fun-service/pkg/logger"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic("加载配置失败: " + err.Error())
	}
	fmt.Println("配置加载成功")

	// 2. 初始化日志
	logCfg := logger.Config{
		Level: cfg.Log.Level,
		Path:  cfg.Log.Path,
		// 添加其他字段映射（如果有的话）
	}
	logger.Init(logCfg)
	fmt.Println("日志系统初始化成功")

	// 3. 初始化数据库
	database.InitMySQL(cfg.MySQL)
	fmt.Println("数据库连接成功")

	// 4. 注册路由
	r := router.SetupRouter()
	fmt.Println("路由注册成功")

	// 5. 启动服务
	logger.Info("服务启动: " + cfg.Server.Addr)
	if err := r.Run(cfg.Server.Addr); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}
