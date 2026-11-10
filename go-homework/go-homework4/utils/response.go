package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 预定义的“业务错误”，handler 里直接用，不用每次拼字符串
var (
	ErrNotFound  = errors.New("资源不存在")
	ErrForbidden = errors.New("没有权限执行此操作")
	ErrConflict  = errors.New("资源已存在")
)

// Response 是所有接口统一的返回格式
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success 返回 200
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "成功", Data: data})
}

// Created 返回 201（新建资源成功）
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{Code: 0, Message: "创建成功", Data: data})
}

// Error 返回指定状态码和错误信息
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{Code: status, Message: message})
}

// HandleError 把内部错误翻译成合适的 HTTP 状态码
func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrForbidden):
		Error(c, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrConflict):
		Error(c, http.StatusConflict, err.Error())
	default:
		Error(c, http.StatusInternalServerError, err.Error())
	}
}
