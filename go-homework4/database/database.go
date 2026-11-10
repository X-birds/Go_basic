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
	// SQLite 的外键约束默认是关闭的，必须在每个连接上显式开启，
	// 否则 schema 里的 FOREIGN KEY 只是摆设，能插入 user_id 不存在的孤儿数据。
	// _pragma=foreign_keys(1) 会对连接池中的所有连接生效。
	dsn := path + "?_pragma=foreign_keys(1)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
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
