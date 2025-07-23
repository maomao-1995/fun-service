package router

import (
	"fun-service/internal/api/handler"
	"fun-service/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Logger())

	// 启动服务，监听 8080 端口
	// 注意：默认是阻塞式的，会一直运行直到被中断

	// 用户路由
	userHandler := handler.NewUserHandler()
	userGroup := r.Group("/users")
	{
		userGroup.POST("", userHandler.CreateUser) // POST /users
		userGroup.GET("/:id", userHandler.GetUser) // GET /users/:id
	}

	return r
}
