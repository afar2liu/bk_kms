package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

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

	// 3. 设置 Gin 运行模式
	ginMode := config.Server.GinMode
	if ginMode == "" {
		ginMode = gin.DebugMode // 默认为 debug 模式
	}
	gin.SetMode(ginMode)
	lib.Logger.Info(fmt.Sprintf("Gin 运行模式: %s", ginMode))

	// 4. 初始化数据库连接
	if err := lib.InitDatabase(config); err != nil {
		lib.Logger.Fatal(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 5. 初始化路由
	router := route.InitRouter()

	// 6. 启动 HTTP 服务器
	addr := fmt.Sprintf(":%d", config.Server.Port)
	lib.Logger.Info(fmt.Sprintf("HTTP 服务器启动在端口: %d", config.Server.Port))

	if err := router.Run(addr); err != nil {
		lib.Logger.Fatal(fmt.Sprintf("启动 HTTP 服务器失败: %v", err))
	}
}
