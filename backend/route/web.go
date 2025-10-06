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

	// API 路由组（需要认证）
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		bookmarkController := controller.NewBookmarkController()
		tagController := controller.NewTagController()

		// 书签相关路由
		api.GET("/bookmarks", bookmarkController.List)
		api.POST("/bookmark", bookmarkController.Create)
		api.PUT("/bookmarks", bookmarkController.Update)
		api.DELETE("/bookmark", bookmarkController.Delete)

		// 标签相关路由
		api.GET("/tags", tagController.List)

		// API v1 路由组
		v1 := api.Group("/v1")
		{
			v1.PUT("/tag/:id", tagController.Update)
			// 书签导入（SSE 流式响应）
			v1.POST("/bookmarks/import", bookmarkController.Import)
		}
	}

	// 书签内容查看（无需认证，根据 OpenAPI 文档）
	r.GET("/bookmark/:id/content", controller.NewBookmarkController().GetContent)

	// 测试路由
	testController := &controller.TestController{}
	r.GET("/api/v1/test/hello", testController.Hello)

	return r
}