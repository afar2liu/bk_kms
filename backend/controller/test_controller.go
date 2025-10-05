package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bk_kms/lib"
)

// TestController 测试控制器
type TestController struct{}

// Hello 测试方法
func (tc *TestController) Hello(c *gin.Context) {
	// 输出日志
	lib.Logger.Info("Hello World")

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": "Hello World!",
	})
}
