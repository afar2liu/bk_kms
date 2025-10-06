package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"bk_kms/lib"
	"bk_kms/model/dto"
	"bk_kms/repo"
	"bk_kms/utils"
)

type AuthController struct {
	userRepo *repo.UserRepo
}

func NewAuthController() *AuthController {
	return &AuthController{
		userRepo: &repo.UserRepo{},
	}
}

// GetCaptcha 获取图形验证码
func (ac *AuthController) GetCaptcha(c *gin.Context) {
	id, base64Img, captchaStr, err := utils.GenerateCaptcha()
	if err != nil {
		lib.Logger.Error("生成验证码失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "生成验证码失败",
		})
		return
	}

	// 构造响应数据
	data := dto.CaptchaData{
		Captcha:   base64Img,
		CaptchaID: id,
	}

	// 非 release 环境返回验证码实际值（用于自动化测试）
	ginMode := lib.GlobalConfig.Server.GinMode
	if ginMode != "release" {
		data.CaptchaStr = captchaStr
	}

	// 返回验证码和验证码ID给客户端
	c.JSON(http.StatusOK, dto.CaptchaResponse{
		Code: 0,
		Msg:  "成功",
		Data: data,
	})
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	// 验证验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.Captcha) {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "验证码错误",
		})
		return
	}

	// 查找用户
	user, err := ac.userRepo.FindByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, dto.Response{
				Code: 1,
				Msg:  "用户名或密码错误",
			})
			return
		}
		lib.Logger.Error("查询用户失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "登录失败",
		})
		return
	}

	// 验证密码
	if !utils.VerifyPassword(req.Pwd, user.Salt, user.Password) {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "用户名或密码错误",
		})
		return
	}

	// 生成 JWT token
	expDuration, _ := time.ParseDuration(lib.GlobalConfig.JWT.Exp)
	if expDuration == 0 {
		expDuration = 24 * time.Hour
	}

	token, err := utils.GenerateToken(user.ID, user.Username, lib.GlobalConfig.JWT.Secret, expDuration)
	if err != nil {
		lib.Logger.Error("生成token失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "登录失败",
		})
		return
	}

	lib.Logger.Info("用户登录成功: " + user.Username)

	c.JSON(http.StatusOK, dto.LoginResponse{
		Code: 0,
		Msg:  "登录成功",
		Data: dto.LoginData{
			ID:       user.ID,
			Username: user.Username,
			Token:    token,
		},
	})
}
