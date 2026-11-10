package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 打印每个请求的方法、路径、状态码、耗时
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		if query != "" {
			path += "?" + query
		}
		log.Printf("[%d] %s %s %v", c.Writer.Status(), c.Request.Method, path, latency)
	}
}
