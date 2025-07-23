package main

import (
	"fun-service/config"
	"fun-service/internal/api/router"
	"fun-service/pkg/logger"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic("加载配置失败: " + err.Error())
	}

	// 2. 初始化日志
	logCfg := logger.Config{
		Level: cfg.Log.Level,
		Path:  cfg.Log.Path,
		// 添加其他字段映射（如果有的话）
	}
	logger.Init(logCfg)

	// 3. 初始化数据库
	// database.InitMySQL(cfg.MySQL)

	// 4. 注册路由
	r := router.SetupRouter()

	// 5. 启动服务
	logger.Info("服务启动: " + cfg.Server.Addr)
	if err := r.Run(cfg.Server.Addr); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}
