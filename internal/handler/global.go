package handler

import (
	"fmt"
	"fun-service/internal/model"
	"fun-service/pkg/database"
	"strconv"

	"fun-service/pkg/jwtMain"
	"fun-service/pkg/redis"
	"fun-service/pkg/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type UserParams struct {
	Username  string `json:"username" binding:"required"`
	Birthdate string `json:"birthdate"`
	Phone     string `json:"phone" binding:"required"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required"`
	Nickname  string `json:"nickname"`
	Code      string `json:"code" binding:"required"` // 验证码
}

// @Summary 用户注册
// @Description 用户注册
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserReq true "用户信息"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"注册成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /register [post]
func Register(c *gin.Context) {
	var params UserParams

	jsonErr := c.ShouldBindJSON(&params)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误", "error": jsonErr.Error()})
		return
	}
	//初始化参数
	//Nickname
	if params.Nickname == "" {
		params.Nickname = utils.GenerateRandomNickname()
	}
	//Birthdate
	birthdate, err := time.Parse("2006-01-02", params.Birthdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "出生日期格式错误，请使用YYYY-MM-DD"})
		return
	}
	//Password
	hashedPassword, passwordErr := bcrypt.GenerateFromPassword(
		[]byte(params.Password),
		bcrypt.DefaultCost, // 默认成本系数10，范围4-31
	)
	if passwordErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 400, "error": "密码加密失败"})
		return
	}
	hashedPasswordString := string(hashedPassword)

	code, getCodeErr := redis.Rdb.Get(redis.Ctx, params.Phone).Result()
	if getCodeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请先发送验证码", "error": getCodeErr.Error()})
		return
	}
	codeTrue := code != params.Code
	if codeTrue {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码错误"})
		return
	}

	//创建model实例
	newUser := model.User{
		Username:  params.Username,
		Birthdate: birthdate,
		Phone:     params.Phone,
		Email:     params.Email,
		Password:  hashedPasswordString,
		Nickname:  params.Nickname,
	}

	selectErr01 := database.DB.Where("phone = ?", params.Phone).First(&newUser).Error
	if selectErr01 == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "手机号已注册"})
		return
	}

	selectErr02 := database.DB.Where("username = ?", params.Username).First(&newUser).Error
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

type SendCodeParams struct {
	Phone string `json:"phone" binding:"required"`
}

// @Summary 发送注册手机验证码
// @Description 发送注册手机验证码
// @Tags users
// @Accept json
// @Produce json
// @Param user body SendCodeReq true "手机号"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"验证码发送成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"参数错误"}"
// @Router /sendCode [post]
func SendCode(c *gin.Context) {
	var params SendCodeParams

	jsonErr := c.ShouldBindJSON(&params)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": jsonErr.Error(), "msg": "参数错误"})
		return
	}

	var userTemp model.User
	selectErr := database.DB.Where("phone = ?", params.Phone).First(&userTemp).Error
	if selectErr == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "手机号已注册"})
		return
	}

	_, getErr := redis.Rdb.Get(redis.Ctx, params.Phone).Result()
	if getErr == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码未过期，请稍后再试"})
		return
	}

	setErr := redis.Rdb.Set(redis.Ctx, params.Phone, "123", 1*time.Minute).Err()
	if setErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "设置验证码失败", "error": setErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "验证码发送成功"})
}

type LoginParams struct {
	Phone    string `json:"phone" binding:""`
	Code     string `json:"code" binding:""`
	Username string `json:"username" binding:""`     // 用户名
	Password string `json:"password" binding:""`     // 密码
	Type     string `json:"type" binding:"required"` // 登录类型，phone 或 username
}

// @Summary 用户登录
// @Description 用户登录
// @Tags users
// @Accept json
// @Produce json
// @Param user body LoginReq true "登录信息"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"登录成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxxx",token:"xxxx"}"
// @Router /login [post]
func Login(c *gin.Context) {
	var params LoginParams

	jsonErr := c.ShouldBindJSON(&params)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": jsonErr.Error(), "msg": "参数错误"})
		return
	}

	var userTemp model.User
	if params.Type == "phone" {
		selectErr01 := database.DB.Where("phone = ?", params.Phone).First(&userTemp).Error
		if selectErr01 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "账号错误"})
			return
		}

	}
	if params.Type == "username" {
		selectErr02 := database.DB.Where("username = ?", params.Username).First(&userTemp).Error
		if selectErr02 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户不存在"})
			return
		}

		passWordErr := bcrypt.CompareHashAndPassword([]byte(userTemp.Password), []byte(params.Password))
		switch {
		case passWordErr == bcrypt.ErrMismatchedHashAndPassword:
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
			return
		case passWordErr != nil:
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码校验失败"})
			return
		default:
		}
	}

	PhoneTemp, PhoneTempErr := strconv.ParseInt(userTemp.Phone, 10, 64)
	if PhoneTempErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "手机号转换失败", "error": PhoneTempErr.Error()})
		return
	}
	fmt.Println("PhoneTemp:", PhoneTemp)

	token, _ := jwtMain.GenerateToken(PhoneTemp, userTemp.Username)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "token": token})
}
