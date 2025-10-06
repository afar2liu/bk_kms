package main

import (
	"fmt"
	"log"
	"time"

	"bk_kms/lib"
	"bk_kms/model/db"
	"bk_kms/utils"
)

func main() {
	// 加载配置
	config, err := lib.LoadConfig("../config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化日志
	if err := lib.InitLogger(config.Log.Level); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer lib.Logger.Sync()

	// 初始化数据库连接
	if err := lib.InitDatabase(config); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 自动迁移表结构
	fmt.Println("开始创建数据库表...")
	if err := lib.DB.AutoMigrate(
		&db.User{},
		&db.Bookmark{},
		&db.Tag{},
		&db.BookmarkTag{},
	); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}
	fmt.Println("数据库表创建成功！")

	// 创建默认管理员用户
	fmt.Println("创建默认管理员用户...")
	var count int64
	lib.DB.Model(&db.User{}).Count(&count)
	if count == 0 {
		salt, _ := utils.GenerateSalt()
		password := utils.HashPassword("admin123", salt)

		admin := &db.User{
			Username:  "admin",
			Password:  password,
			Salt:      salt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := lib.DB.Create(admin).Error; err != nil {
			log.Fatalf("创建管理员用户失败: %v", err)
		}
		fmt.Println("管理员用户创建成功！")
		fmt.Println("用户名: admin")
		fmt.Println("密码: admin123")
	} else {
		fmt.Println("用户已存在，跳过创建")
	}

	fmt.Println("数据库初始化完成！")
}
