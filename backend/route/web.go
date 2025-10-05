package route

import (
	"github.com/gin-gonic/gin"

	"bk_kms/controller"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 测试路由
		testGroup := v1.Group("/test")
		{
			testController := &controller.TestController{}
			testGroup.GET("/hello", testController.Hello)
		}
	}

	return r
}