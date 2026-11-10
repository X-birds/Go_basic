package handlers

import (
	"errors"
	"go-homework4/middleware"
	"go-homework4/models"
	"go-homework4/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostHandler 文章相关的接口
type PostHandler struct {
	db *gorm.DB
}

// NewPostHandler 创建一个文章相关的接口
func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{db: db}
}

// 设置分页常量
const (
	defaultPageSize = 10
	maxPageSize     = 50
)

// 解析url中的:id参数，parseID 解析文章ID，失败时写好400错误响应
func parseID(c *gin.Context) (uint, bool) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 {
		utils.Error(c, http.StatusBadRequest, "无效的文章ID")
		return 0, false
	}
	return uint(id64), true
}

// 获取文章list，并支持分页
func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	//兜底，防止传入 0、负数或特别大的值
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	// 统计文章数量
	var total int64
	if err := h.db.Model(&models.Post{}).Count(&total).Error; err != nil {
		utils.HandleError(c, err)
		return
	}

	//分页查询文章
	var posts []models.Post
	err := h.db.Preload("User").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&posts).Error
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	//返回结果
	utils.Success(c, gin.H{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// get 获取文章详情（内容、作者、评论等）
func (h *PostHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	//查询文章
	var post models.Post
	err := h.db.Preload("User").Preload("Comments.User").First(&post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(c, utils.ErrNotFound)
			return
		}
		utils.HandleError(c, err)
		return
	}
	//返回结果
	utils.Success(c, post)
}

// Create 创建文章
func (h *PostHandler) Create(c *gin.Context) {
	var req models.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "请求参数错误："+err.Error())
		return
	}

	//创建文章
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  middleware.GetUserID(c),
	}
	if err := h.db.Create(&post).Error; err != nil {
		utils.HandleError(c, err)
		return
	}
	//重新查询一次，吧作责信息也带上再返回
	h.db.Preload("User").First(&post, post.ID)
	utils.Created(c, post)
}

// 更新文章Update(作者本人)
func (h *PostHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	//接收请求参数
	var req models.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "请求参数错误："+err.Error())
		return
	}
	//查询文章
	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(c, utils.ErrNotFound)
			return
		}
		utils.HandleError(c, err)
		return
	}
	//检查作者是否是当前登录用户
	if post.UserID != middleware.GetUserID(c) {
		utils.HandleError(c, utils.ErrForbidden)
		return
	}
	//更新文章
	if err := h.db.Model(&post).Updates(map[string]any{
		"title":   req.Title,
		"content": req.Content,
	}).Error; err != nil {
		utils.HandleError(c, err)
		return
	}
	//重新查询一次，拿到更新后的内容和作者信息再返回
	h.db.Preload("User").First(&post, post.ID)
	//返回结果
	utils.Success(c, post)
}

// delete 删除文章
func (h *PostHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	// 查询文章
	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(c, utils.ErrNotFound)
			return
		}
		utils.HandleError(c, err)
		return
	}
	//检查作者是否为当前登录用户
	if post.UserID != middleware.GetUserID(c) {
		utils.HandleError(c, utils.ErrForbidden)
		return
	}
	//删除文章
	if err := h.db.Delete(&post).Error; err != nil {
		utils.HandleError(c, err)
		return
	}
	//返回数据
	utils.Success(c, gin.H{
		"deleted_id": post.ID,
	})
}
