package handler

import (
	"fmt"
	"fun-service/internal/model"
	"fun-service/pkg/database"
	"fun-service/pkg/redis"
	"fun-service/pkg/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// user
type UserReq struct {
	Username  string `json:"username" binding:"required"`
	Birthdate string `json:"birthdate"`
	Phone     string `json:"phone" binding:"required"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required"`
	Nickname  string `json:"nickname"`
}	

// 用户注册
func UserRegister(c *gin.Context) {
	var req UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	//初始化参数
	if req.Nickname == "" {
		req.Nickname = utils.GenerateRandomNickname()
	}

	birthdate, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "出生日期格式错误，请使用YYYY-MM-DD"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost, // 默认成本系数10，范围4-31
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 400, "error": "密码加密失败"})
		return
	}
	hashedPasswordString := string(hashedPassword)

	//创建model实例
	newUser := model.User{
		Username:  req.Username,
		Birthdate: birthdate,
		Phone:     req.Phone,
		Email:     req.Email,
		Password:  hashedPasswordString,
		Nickname:  req.Nickname,
	}

	// var tempUser model.User
	if err := database.DB.Where("phone = ?", req.Phone).First(&newUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "手机号已注册"})
		return
	}
	if err := database.DB.Where("username = ?", req.Username).First(&newUser).Error; err == nil {
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

// 发送验证码
func SendCode(c *gin.Context) {
	// 获取手机号
	phone := c.Query("phone")
	// err := redis.Rdb.Set(redis.Ctx, "name", "张三", 1*time.Minute).Err()
	val, err := redis.Rdb.Get(redis.Ctx, "name").Result()
	if  err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "获取验证码失败", "error": err.Error()})
	}
	fmt.Println("Redis Set and Get:", val,phone)
}