package database

import (
	"go-homework4/models"

	"github.com/libtnb/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 数据库连接、创建数据库表
func Connect(path string) (*gorm.DB, error) {
	//打开数据库文件；文件不存在会自动创建
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	//  把 Go 结构体同步成数据库表（建表 / 补字段）
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		return nil, err
	}
	return db, nil
}
