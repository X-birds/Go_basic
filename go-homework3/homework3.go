package main

import (
	"fmt"

	"github.com/libtnb/sqlite"
	"gorm.io/gorm"
)

// 题目1：模型定义
// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
// 要求 ：
// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
// 编写Go代码，使用Gorm创建这些模型对应的数据库表。
type User struct {
	gorm.Model
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"uniqueIndex"`
	Age       uint8  `gorm:"not null"`
	Postcount int
	Posts     []Post
}
type Post struct {
	gorm.Model
	Title         string `gorm:"size:200;not null"`
	Content       string `gorm:"type:text"`
	UserID        uint   `gorm:"not null;index"`
	User          User
	CommentStatus string
	Comments      []Comment
}
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	PostID  uint   `gorm:"not null;index"`
	Post    Post
}

// 题目2：关联查询
// 基于上述博客系统的模型定义。
// 要求 ：
// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
func QueryUserPosts_With_Comments(db *gorm.DB, userID uint) {
	var user User

	result := db.Preload("Posts").Preload("Posts.Comments").First(&user, userID)
	if result.Error != nil {
		fmt.Println("查询失败：", result.Error)
		return
	}
	fmt.Printf("\n 用户：%s（ID：%d）发布了%d篇文章.\n", user.Name, user.ID, len(user.Posts))
	for i, post := range user.Posts {
		fmt.Printf("文章 %d -- %s，评论数为：%d\n", i+1, post.Title, len(post.Comments))
		fmt.Printf("文章内容： %s\n", post.Content)
		if len(post.Comments) > 0 {
			for j, comment := range post.Comments {
				fmt.Printf("评论 %d: %s\n", j+1, comment.Content)
			}
		}
	}
}

func QureyMaxcomments_Post(db *gorm.DB) {
	var posts []Post
	db.Preload("Comments").Find(&posts)
	maxcount := 0

	var bestpost Post
	for _, post := range posts {
		commentcount := len(post.Comments)

		if commentcount > maxcount {
			maxcount = commentcount
			bestpost = post
		}
	}

	fmt.Printf("评论最多的文章：%s，评论数为：%d\n", bestpost.Title, maxcount)

}

// 题目3：钩子函数
// 继续使用博客系统的模型。
// 要求 ：
// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，
// 如果评论数量为 0，则更新文章的评论状态为 "无评论"。
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	fmt.Println("新增文章--更新用户文章数量：")
	tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("Postcount",
		gorm.Expr("Postcount + 1"))
	fmt.Printf("用户：%d 的文章数量已更新\n", p.UserID)
	return nil
}
func (p *Post) AfterDelete(tx *gorm.DB) error {
	fmt.Println("删除文章--更新用户文章数量：")
	tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("Postcount",
		gorm.Expr("Postcount - 1"))
	fmt.Printf("用户：%d 的文章数量已更新\n", p.UserID)
	return nil
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	fmt.Printf("钩子触发：评论删除 (PostID: %d)\n", c.PostID)

	// 1. 查询该文章还剩多少评论
	var count int64
	tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count)

	fmt.Printf("文章 %d 剩余评论数：%d\n", c.PostID, count)

	// 2. 如果评论数为 0，更新文章状态
	if count == 0 {
		tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Update("comment_status", "无评论")
		fmt.Printf("文章 %d 状态已更新为：无评论\n", c.PostID)
	}

	return nil
}

func main() {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}

	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("迁移失败：" + err.Error())
	}
	fmt.Println("数据库表创建成功！")

	users := []User{
		{Name: "张三", Email: "123@123.com", Age: 18},
		{Name: "李四", Email: "345@345.com", Age: 18},
	}
	if err := db.Create(&users).Error; err != nil {
		panic("导入用户数据失败：" + err.Error())
	}
	posts := []Post{
		{Title: "张三的第一篇博客", Content: "111111111111", UserID: users[0].ID},
		{Title: "李四的第一篇博客", Content: "222222222222", UserID: users[1].ID},
		{Title: "李四的第二篇博客", Content: "333333333333", UserID: users[1].ID},
	}
	if err := db.Create(&posts).Error; err != nil {
		panic("导入文章数据失败：" + err.Error())
	}
	comments := []Comment{
		{Content: "张三写的好！", PostID: posts[0].ID},
		{Content: "李四写的好！", PostID: posts[1].ID},
		{Content: "第一次评论好！", PostID: posts[2].ID},
		{Content: "第二次评论好！", PostID: posts[2].ID},
	}
	if err := db.Create(&comments).Error; err != nil {
		panic("导入评论数据失败：" + err.Error())
	}

	fmt.Println("测试数据插入成功！")

	fmt.Println("=== 测试2：删除文章1的所有剩余评论 ===")
	db.Where("post_id = ?", posts[0].ID).Delete(&Comment{})

	QueryUserPosts_With_Comments(db, 1)

	QureyMaxcomments_Post(db)
}
