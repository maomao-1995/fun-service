package middleware

import (
	"time"

	"fun-service/pkg/logger" // 引入之前实现的日志工具

	"github.com/gin-gonic/gin"
)

// Logger 记录请求日志的中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 请求开始时间
		startTime := time.Now()

		// 2. 处理请求
		c.Next()

		// 3. 请求结束后记录日志
		endTime := time.Now()
		latency := endTime.Sub(startTime) // 耗时
		method := c.Request.Method        // 请求方法
		path := c.Request.URL.Path        // 请求路径
		statusCode := c.Writer.Status()   // 状态码
		clientIP := c.ClientIP()          // 客户端IP

		// 记录日志
		logger.Infof(
			"method=%s path=%s status=%d ip=%s latency=%v",
			method, path, statusCode, clientIP, latency,
		)
	}
}
