package main

import (
	"fmt"
	"go-homework4/config"
	"go-homework4/database"
	"go-homework4/handlers"
	"go-homework4/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 读取配置
	cfg := config.LoadConfig()

	//连接数据库，并建立表
	db, err := database.Connect(cfg.DB.Path)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	//构造handler，把依赖注入
	secret := []byte(cfg.JWT.Secret)
	authH := handlers.NewAuthHandler(db, secret)
	postH := handlers.NewPostHandler(db)
	commentH := handlers.NewCommentHandler(db)

	// 创建引擎，并挂全局中间件
	r := gin.New()
	r.Use(middleware.Logger(), gin.Recovery())

	//公开接口
	api := r.Group("/api")
	{
		api.POST("/register", authH.Register)
		api.POST("/login", authH.Login)

		api.GET("/posts", postH.List)
		api.GET("/posts/:id", postH.Get)
		api.GET("/posts/:id/comments", commentH.List)
	}
	//需要认证的接口
	protected := api.Group("")
	protected.Use(middleware.Auth(secret))
	{
		protected.POST("/posts", postH.Create)
		protected.PUT("/posts/:id", postH.Update)
		protected.DELETE("/posts/:id", postH.Delete)
		protected.POST("/posts/:id/comments", commentH.Create)
	}

	//启动
	addr := ":" + cfg.Server.Port
	fmt.Printf("服务已启动: http://localhost%s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}

}
