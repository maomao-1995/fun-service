package middleware

import (
	"fun-service/pkg/jwtMain"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "no token"})
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		claims, claimsErr := jwtMain.ParseToken(tokenStr)
		if claimsErr != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "非法或过期 token", "error": claimsErr.Error()})
			return
		}
		c.Set("username", claims.Username)
		c.Set("userPhone", claims.UserPhone)
		c.Next()
	}
}
