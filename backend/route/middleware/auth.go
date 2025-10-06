package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"bk_kms/lib"
	"bk_kms/model/dto"
	"bk_kms/utils"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code: 401,
				Msg:  "未授权，请先登录",
			})
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code: 401,
				Msg:  "Authorization 格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 token
		claims, err := utils.ParseToken(tokenString, lib.GlobalConfig.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code: 401,
				Msg:  "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}