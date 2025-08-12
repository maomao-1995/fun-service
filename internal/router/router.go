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

	// 配置静态文件服务
	r.Static("/uploads", "./uploads")

	// 全局中间件
	r.Use(middleware.Logger())

	// 注册Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务，监听 8080 端口
	// 注意：默认是阻塞式的，会一直运行直到被中断

	//全局路由
	r.POST("/register", handler.Register)
	r.POST("/sendCode", handler.SendCode)
	r.POST("/login", handler.Login)
	r.GET("/refresh", handler.Refresh)
	r.POST("/upload", handler.Upload)
	// 用户路由
	userGroup := r.Group("/user")
	userGroup.Use(middleware.ParseToken())
	{
		userGroup.GET("/info", handler.UserInfo)
	}
	//emoji路由
	emojiGroup := r.Group("/emoji")
	emojiGroup.Use(middleware.ParseToken())
	{
		emojiGroup.GET("/detail", handler.EmojiDetail)
		emojiGroup.POST("/add", handler.EmojiAdd)
		emojiGroup.POST("/delete", handler.EmojiDelete)
	}
	return r
}
