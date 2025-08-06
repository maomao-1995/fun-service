package handler

import (
	"fun-service/internal/model"
	"fun-service/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	username := c.MustGet("username").(string)
	userPhone := c.MustGet("userPhone").(int64)

	var userTemp model.User
	selectErr01 := database.DB.Where("phone = ? OR username = ?", userPhone, username).First(&userTemp).Error
	if selectErr01 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取用户信息成功", "data": userTemp})
}
