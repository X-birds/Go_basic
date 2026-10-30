package handlers

import (
	"go-homework4/middleware"
	"go-homework4/models"
	"go-homework4/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

// NewCommentHandler 创建一个评论相关的接口
func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

// Create 给文章发布评论
func (h *CommentHandler) Create(c *gin.Context) {
	postID, ok := parseID(c)
	if !ok {
		return
	}
	//解析请求参数
	var req models.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	//确认文章存在
	var count int64
	if err := h.db.Model(&models.Post{}).Where("id = ?", postID).Count(&count).Error; err != nil {
		utils.HandleError(c, err)
		return
	}
	if count == 0 {
		utils.Error(c, http.StatusNotFound, "文章不存在")
		return
	}

	// 创建评论
	comment := models.Comment{
		PostID:  postID,
		UserID:  middleware.GetUserID(c),
		Content: req.Content,
	}
	if err := h.db.Create(&comment).Error; err != nil {
		utils.HandleError(c, err)
		return
	}
	//带评论人信息返回
	h.db.Preload("User").First(&comment, comment.ID)
	utils.Created(c, comment)
}

// List 获取某个文章的评论列表
func (h *CommentHandler) List(c *gin.Context) {
	postID, ok := parseID(c)
	if !ok {
		return
	}
	//查询评论列表
	var comments []models.Comment
	if err := h.db.Preload("User").Where("post_id = ?", postID).Order("created_at ASC").Find(&comments).Error; err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, comments)
}
