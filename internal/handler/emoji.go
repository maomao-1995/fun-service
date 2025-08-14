package handler

import (
	"encoding/json"
	"fun-service/internal/model"
	"fun-service/pkg/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)
type EmojiDetailRequest struct {
	ID uint `json:"id" binding:"required"`
}
type EmojiDetailDTO struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	ViewCount        int    `json:"view_count"`
	CollectionCount  int    `json:"collection_count"`
	Tags             datatypes.JSON	 `json:"tags"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// @Summary 获取表情详情
// @Description 获取表情详情
// @Tags emoji
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param emoji body EmojiDetailRequest true "Emoji ID"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"获取表情详情成功","data":EmojiDetailDTO}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emoji/detail [get]
func EmojiDetail(c *gin.Context) {
	var params EmojiDetailRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	var emoji model.Emoji
	var emojiDetailTemp EmojiDetailDTO
	if err := database.DB.Model(&model.Emoji{}).Where("id = ?", params.ID).First(&emoji).Scan(&emojiDetailTemp).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "msg": "表情未找到", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "获取表情详情成功", "data": emojiDetailTemp})
}

type EmojiAddRequest struct {
	Name string   `json:"name" binding:"required"`
	URL  string   `json:"url" binding:"required"`
	Tags []string `json:"tags" binding:"omitempty,dive,required"`
}

// @Summary 添加表情
// @Description 添加表情
// @Tags emoji
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param emoji body EmojiAddRequest true "Emoji信息"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"添加表情成功","data":Emoji}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emoji/add [post]
func EmojiAdd(c *gin.Context) {
	var params EmojiAddRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	//处理参数
	temp, _ := json.Marshal(params.Tags)
	tagsTemp := datatypes.JSON(json.RawMessage(temp))

	newEmoji := model.Emoji{
		Name: params.Name,
		URL:  params.URL,
		Tags: tagsTemp,
		AuthorUUID: c.GetString("uuid"), // 从上下文中获取用户UUID
	}

	if err := database.DB.Model(&model.Emoji{}).Create(&newEmoji).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建emoji失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建emoji成功",})

}

type EmojiDeleteRequest struct {
	ID uint `json:"id" binding:"required"`
}

// @Summary 删除表情
// @Description 删除表情
// @Tags emoji
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param emoji body EmojiDeleteRequest true "Emoji ID"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"删除表情成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emoji/delete [post]
func EmojiDelete(c *gin.Context) {
	var params EmojiDeleteRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	if err := database.DB.Model(&model.Emoji{}).Where("id = ?", params.ID).Delete(&model.Emoji{}).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "msg": "表情未找到", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "删除表情成功"})

}
