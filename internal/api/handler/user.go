package handler

import (
	"fmt"
	"fun-service/internal/model"
	"fun-service/pkg/database"
	"fun-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// user
type UserReq struct {
	Username string `json:"username" binding:"required"`
	Age      int    `json:"age" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
}

// 用户注册
func UserRegister(c *gin.Context) {
	var req UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}
	fmt.Println("Received registration request:", req)
	//初始化默认值
	if req.Nickname == "" {
		req.Nickname = utils.GenerateRandomNickname()
	}

	newUser := model.User{
		Username: req.Username,
		Age:      req.Age,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	}

	var tempUser model.User
	if err := database.DB.Where("phone = ?", req.Phone).First(&tempUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "手机号已注册"})
		return
	}
	if err := database.DB.Where("username = ?", req.Username).First(&tempUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户名已存在"})
		return
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "注册失败，请稍后重试",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}
