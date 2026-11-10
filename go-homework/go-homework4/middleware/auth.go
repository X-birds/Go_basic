package middleware

import (
	"go-homework4/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 上下文中存用户名和 ID 用的键，定义成常量避免手写字符串写错
const (
	ContextUserIDKey   = "userID"
	ContextUsernameKey = "username"
)

// Auth 返回一个"验票中间件"；secret 是验证 token 用的密钥
func Auth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ① 从请求头取 Authorization
		header := c.GetHeader("Authorization")
		if header == "" {
			utils.Error(c, http.StatusUnauthorized, "请先登录")
			c.Abort()
			return
		}
		// ② 校验格式必须是 "Bearer <token>"
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			utils.Error(c, http.StatusUnauthorized, "Authorization 格式错误为 Bearer <token>")
			c.Abort()
			return
		}
		// ③ 去掉前缀，拿到纯 token
		tokenString := strings.TrimSpace(strings.TrimPrefix(header, prefix))
		if tokenString == "" {
			utils.Error(c, http.StatusUnauthorized, "token 不能为空")
			c.Abort()
			return
		}
		// ④ 验证并解析
		claims, err := utils.ParseToken(tokenString, secret)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "token 无效或过期")
			c.Abort()
			return
		}
		// ⑤ 把用户信息放进上下文，后面的 handler 能取到
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)

		// ⑥ 放行
		c.Next()
	}
}

// GetUserID 从上下文取当前登录用户 ID（未登录时为 0）
func GetUserID(c *gin.Context) uint {
	return c.GetUint(ContextUserIDKey)
}
