package dto

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`   // 用户名
	Pwd       string `json:"pwd" binding:"required"`        // 密码
	Captcha   string `json:"captcha" binding:"required"`    // 验证码内容
	CaptchaID string `json:"captcha_id" binding:"required"` // 验证码ID
}

// LoginData 登录响应数据
type LoginData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data LoginData `json:"data"`
}

// CaptchaData 验证码响应数据
type CaptchaData struct {
	Captcha    string `json:"captcha"`              // 图形验证码，base64编码
	CaptchaID  string `json:"captcha_id"`           // 图形验证码ID
	CaptchaStr string `json:"captcha_str,omitempty"` // 验证码实际值（仅非release环境）
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data CaptchaData `json:"data"`
}
