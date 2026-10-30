package handlers

import (
	"errors"
	"go-homework4/models"
	"go-homework4/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler 把 handler 需要的依赖挂在结构体上
type AuthHandler struct {
	db        *gorm.DB
	jwtSecret []byte
}

// NewAuthHandler 构造 handler；main.go 里会调用它
func NewAuthHandler(db *gorm.DB, secret []byte) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: secret}
}

// Register 处理 POST /api/v1/register
func (h *AuthHandler) Register(c *gin.Context) {
	// ① 解析 + 校验请求体
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	// ② 检查用户名是否已被占用
	var count int64
	if err := h.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		utils.HandleError(c, err)
		return
	}
	if count > 0 {
		utils.Error(c, http.StatusConflict, "用户名已被占用")
		return
	}
	// ③ 检查邮箱是否已被注册
	if err := h.db.Model(&models.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		utils.HandleError(c, err)
		return
	}
	if count > 0 {
		utils.Error(c, http.StatusConflict, "邮箱已被注册")
		return
	}
	// ④ 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	// ⑤ 组装模型并入库
	user := models.User{
		Username: req.Username,
		Password: string(hash),
		Email:    req.Email,
	}
	if err := h.db.Create(&user).Error; err != nil {
		utils.HandleError(c, err)
		return
	}

	// ⑥ 返回安全的用户信息（不含密码）
	utils.Created(c, user.ToResponse())
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	// ① 按用户名查用户
	var user models.User
	err := h.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，和密码错误返回同一句话
			utils.Error(c, http.StatusUnauthorized, "用户名或密码错误")
			return
		}
		utils.HandleError(c, err)
		return
	}
	// ② 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// ③ 签发 token
	token, err := utils.GenerateToken(h.jwtSecret, user.ID, user.Username)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "token 生成失败")
		return
	}
	// ④ 返回 token + 用户信息
	utils.Success(c, gin.H{
		"token": token,
		"user":  user.ToResponse(),
	})
}
