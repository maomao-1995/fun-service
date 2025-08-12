package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"fun-service/internal/model"
	"fun-service/pkg/database"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fun-service/pkg/jwtMain"
	"fun-service/pkg/redisMain"
	"fun-service/pkg/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Tags global
// @Accept json
// @Produce json
// @Param user body UserParams true "用户信息"
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
	birthdate, birthdateErr := time.Parse("2006-01-02", params.Birthdate)
	if birthdateErr != nil {
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

	code, getCodeErr := redisMain.Rdb.Get(redisMain.Ctx, params.Phone).Result()
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
// @Tags global
// @Accept json
// @Produce json
// @Param user body SendCodeParams true "手机号"
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

	_, getErr := redisMain.Rdb.Get(redisMain.Ctx, params.Phone).Result()
	if getErr == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码未过期，请稍后再试"})
		return
	}

	setErr := redisMain.Rdb.Set(redisMain.Ctx, params.Phone, "123", 1*time.Minute).Err()
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
// @Tags global
// @Accept json
// @Produce json
// @Param user body LoginParams true "登录信息"
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

	token, _ := jwtMain.GenerateToken(PhoneTemp, userTemp.Username, time.Now().Add(50*time.Minute))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data": gin.H{"token": token}})
}

// @Summary 刷新重置token
// @Description 刷新重置token
// @Tags global
// @Accept json
// @Produce json
// @Param refresh_token header string true "refresh_token" // 需要在请求头中传递 refresh_token
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"刷新重置token成功","token":"xxxx"}"
// @Failure 400 {object} map[string]interface{} "{"code":401,"msg":"xxxxx"}"
// @Router /refresh [get]
func Refresh(c *gin.Context) {
	tokenStr := c.GetHeader("refresh_token")
	if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
		c.AbortWithStatusJSON(401, gin.H{"code": 400, "msg": "no token"})
		return
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	claims, claimsErr := jwtMain.ParseToken(tokenStr)
	if claimsErr != nil {
		c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "非法或过期 token", "error": claimsErr.Error()})
		return
	}

	token01, _ := jwtMain.GenerateToken(claims.UserPhone, claims.Username, time.Now().Add(5*time.Minute))
	token02, _ := jwtMain.GenerateToken(claims.UserPhone, claims.Username, time.Now().Add(7*24*time.Minute))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "刷新重置token成功", "data": gin.H{"accessToken": token01, "refreshToken": token02}})
}

// @Summary 文件上传
// @Description 文件上传
// @Tags global
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"文件上传成功","url":"xxxx"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxxx"}"
// @Router /upload [post]
func Upload(c *gin.Context) {
	// 获取图片文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "未上传文件"})
		return
	}
	defer file.Close()

	// 计算文件的 MD5 哈希值
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "计算文件哈希值失败"})
		return
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	// 重置文件指针
	file.Seek(0, 0)

	// 检查文件是否已存在
	if existingFileName, ok := utils.CheckHash(fileHash); ok {
		// 文件内容已存在，返回文件的访问 URL
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "文件已存在，无需重复上传",
			"url":     fmt.Sprintf("/uploads/%s", existingFileName),
		})
		return
	}

	// 创建文件夹
	dirPath := "./uploads"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm) // 创建目录
	}

	// 生成新的文件名（加入 UUID）
	fileExtension := filepath.Ext(header.Filename)
	newFileName := uuid.New().String() + fileExtension
	newFilePath := filepath.Join(dirPath, newFileName)

	// 保存文件
	if err := c.SaveUploadedFile(header, newFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "保存文件失败"})
		return
	}

	// 保存哈希值到 Redis
	utils.SaveHash(fileHash, newFileName)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
		"url":     fmt.Sprintf("/uploads/%s", newFileName),
	})
}
