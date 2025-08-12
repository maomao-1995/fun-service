package handler

import (
	"fun-service/internal/model"
	"fun-service/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserDTO struct {
	Username string `json:"username"`
	Phone    int64  `json:"phone"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Id       int64  `json:"id"`
}

// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"获取用户信息成功","data":UserDTO}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxxx"}"
// @Router /user/info [get]
func UserInfo(c *gin.Context) {
	username := c.MustGet("username").(string)
	userPhone := c.MustGet("userPhone").(int64)

	var userTemp UserDTO
	err := database.DB.
		Model(&model.User{}).
		Select("id,username, phone,email,nickname").
		Where("phone = ? OR username = ?", userPhone, username).
		Scan(&userTemp).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取用户信息成功", "data": userTemp})
}
