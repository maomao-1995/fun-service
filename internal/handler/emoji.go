package handler

import (
	"encoding/json"
	"fmt"
	"fun-service/internal/model"
	"fun-service/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func EmojiDetail(c *gin.Context) {
	// 处理获取表情详情的逻辑
}

type EmojiAddRequest struct {
	Name string   `json:"name" binding:"required"`
	URL  string   `json:"url" binding:"required"`
	Tags []string `json:"tags" binding:"omitempty,dive,required"`
}

func EmojiAdd(c *gin.Context) {
	var params EmojiAddRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	//处理参数
	temp, _ := json.Marshal(params.Tags)
	tagsTemp := datatypes.JSON(json.RawMessage(temp))
	fmt.Println("tagsTemp", tagsTemp)

	newEmoji := model.Emoji{
		Name: params.Name,
		URL:  params.URL,
		Tags: tagsTemp,
	}

	if err := database.DB.Model(&model.Emoji{}).Create(&newEmoji).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建emoji失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建emoji成功", "data": newEmoji})

}

func EmojiDelete(c *gin.Context) {}
