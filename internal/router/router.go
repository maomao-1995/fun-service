package router

import (
	"fun-service/internal/handler"
	"fun-service/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Logger())

	// 注册Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务，监听 8080 端口
	// 注意：默认是阻塞式的，会一直运行直到被中断

	// 用户路由
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", handler.UserRegister) // POST /user/register
		userGroup.POST("/sendCode", handler.SendCode)     // POST /user/sendCode
	}
	return r
}
