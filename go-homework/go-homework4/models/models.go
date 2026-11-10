package models

import (
	"time"

	"gorm.io/gorm"
)

// ---------- 数据库模型 ----------
// User 用户表
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Password string `gorm:"not null;size:255" json:"-"`
	Email    string `gorm:"uniqueIndex;not null;size:100" json:"-"`

	Posts []Post `json:"-"`
}

// 文章表
type Post struct {
	gorm.Model
	Title   string `gorm:"not null;size:100" json:"title"`
	Content string `gorm:"not null;type:text" json:"content"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`

	User     User      `json:"user"`
	Comments []Comment `json:"comments,omitempty"`
}

// 评论表
type Comment struct {
	gorm.Model
	Content string `gorm:"not null;type:text" json:"content"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	PostID  uint   `gorm:"not null;index" json:"post_id"`

	User User `json:"user"`
	Post Post `json:"-"`
}

// ---------- 请求结构（前端传进来的 JSON） ----------

// 注册用户请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
}

// 登录用户请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// 文章请求
type PostRequest struct {
	Title   string `json:"title" binding:"required,max=100"`
	Content string `json:"content" binding:"required"`
}

// 评论请求
type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// ---------- 响应结构（返回给前端的 JSON） ----------
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 把数据库模型转成安全的响应结构
func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
