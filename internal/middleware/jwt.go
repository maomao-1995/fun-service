package middleware

import (
	"fmt"
	"fun-service/pkg/jwtMain"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"msg": "no token"})
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		fmt.Println("tokenStr:", tokenStr)
		claims, claimsErr := jwtMain.ParseToken(tokenStr)
		if claimsErr != nil {
			c.AbortWithStatusJSON(401, gin.H{"msg": "非法或过期 token", "error": claimsErr.Error()})
			return
		}
		c.Set("username", claims.Username)
		c.Set("userPhone", claims.UserPhone)
		c.Next()
	}
}
