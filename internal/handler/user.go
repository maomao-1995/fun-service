package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	// auth := c.GetHeader("Authorization")
	// claims, err := ParseToken(auth)
	// if err != nil {
	// 	c.JSON(401, gin.H{"msg": "非法或过期 token"})
	// 	return
	// }
	// c.JSON(200, claims)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取用户信息成功"})
}
