package route

import (
	"github.com/gin-gonic/gin"

	"bk_kms/controller"
	"bk_kms/route/middleware"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 认证相关路由（无需认证）
	authController := controller.NewAuthController()
	r.GET("/api/v1/captcha", authController.GetCaptcha)
	r.POST("/api/v1/auth/login", authController.Login)

	// API v1 路由组（需要认证）
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		bookmarkController := controller.NewBookmarkController()
		tagController := controller.NewTagController()

		// 书签相关路由
		v1.GET("/bookmarks", bookmarkController.List)
		v1.POST("/bookmark", bookmarkController.Create)
		v1.PUT("/bookmarks", bookmarkController.Update)
		v1.DELETE("/bookmark", bookmarkController.Delete)
		v1.GET("/bookmark/:id/content", bookmarkController.GetContent)

		// 书签导入（SSE 流式响应）
		v1.POST("/bookmarks/import", bookmarkController.Import)

		// 标签相关路由
		v1.GET("/tags", tagController.List)
		v1.PUT("/tag/:id", tagController.Update)
	}

	// 测试路由
	testController := &controller.TestController{}
	r.GET("/api/v1/test/hello", testController.Hello)

	return r
}