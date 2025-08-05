package handler

import (
	"fun-service/internal/model"
	"fun-service/pkg/database"
	"fun-service/pkg/redis"
	"fun-service/pkg/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type UserReq struct {
	Username  string `json:"username" binding:"required"`
	Birthdate string `json:"birthdate"`
	Phone     string `json:"phone" binding:"required"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required"`
	Nickname  string `json:"nickname"`
	Code      string `json:"code" binding:"required"` // 验证码
}
// register godoc
// @Summary 用户注册
// @Description 用户注册
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserReq true "用户信息"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"注册成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /user/register [post]
func UserRegister(c *gin.Context) {
	var req UserReq

	jsonErr := c.ShouldBindJSON(&req)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误", "error": jsonErr.Error()})
		return
	}

	//初始化参数
	//Nickname
	if req.Nickname == "" {
		req.Nickname = utils.GenerateRandomNickname()
	}
	//Birthdate
	birthdate, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "出生日期格式错误，请使用YYYY-MM-DD"})
		return
	}
	//Password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost, // 默认成本系数10，范围4-31
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 400, "error": "密码加密失败"})
		return
	}
	hashedPasswordString := string(hashedPassword)

	code, getCodeErr := redis.Rdb.Get(redis.Ctx, req.Phone).Result()
	if getCodeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请先发送验证码", "error": getCodeErr.Error()})
		return
	}
	codeTrue := code != req.Code
	if codeTrue {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码错误"})
		return
	}

	//创建model实例
	newUser := model.User{
		Username:  req.Username,
		Birthdate: birthdate,
		Phone:     req.Phone,
		Email:     req.Email,
		Password:  hashedPasswordString,
		Nickname:  req.Nickname,
	}

	selectErr01 := database.DB.Where("phone = ?", req.Phone).First(&newUser).Error
	if selectErr01 == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "手机号已注册"})
		return
	}

	selectErr02 := database.DB.Where("username = ?", req.Username).First(&newUser).Error
	if selectErr02 == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户名已存在"})
		return
	}

	selectErr03 := database.DB.Create(&newUser).Error
	if selectErr03 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败，请稍后重试", "error": selectErr03.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

type SendCodeReq struct {
	Phone string `json:"phone" binding:"required"`
}
// SendCode godoc
// @Summary 发送注册手机验证码
// @Description 发送注册手机验证码
// @Tags users
// @Accept json
// @Produce json
// @Param user body SendCodeReq true "手机号"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"验证码发送成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"参数错误"}"
// @Router /user/sendCode [post]
func SendCode(c *gin.Context) {
	var req SendCodeReq

	jsonErr := c.ShouldBindJSON(&req)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": jsonErr.Error(), "msg": "参数错误"})
		return
	}

	//创建model实例
	newUser := model.User{
		Phone: req.Phone,
	}
	selectErr := database.DB.Where("phone = ?", req.Phone).First(&newUser).Error
	if selectErr == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "手机号已注册"})
		return
	}

	_, getErr := redis.Rdb.Get(redis.Ctx, req.Phone).Result()
	if getErr == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码未过期，请稍后再试"})
		return
	}

	setErr := redis.Rdb.Set(redis.Ctx, req.Phone, "123", 1*time.Minute).Err()
	if setErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "设置验证码失败", "error": setErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "验证码发送成功"})
}
