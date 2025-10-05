package main

import (
	"fmt"
	"log"

	"bk_kms/lib"
	"bk_kms/route"
)

func main() {
	// 1. 加载配置文件
	config, err := lib.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 2. 初始化日志
	if err := lib.InitLogger(config.Log.Level); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer lib.Logger.Sync()

	lib.Logger.Info("项目启动中...")

	// 3. 初始化路由
	router := route.InitRouter()

	// 4. 启动 HTTP 服务器
	addr := fmt.Sprintf(":%d", config.Server.Port)
	lib.Logger.Info(fmt.Sprintf("HTTP 服务器启动在端口: %d", config.Server.Port))

	if err := router.Run(addr); err != nil {
		lib.Logger.Fatal(fmt.Sprintf("启动 HTTP 服务器失败: %v", err))
	}
}
