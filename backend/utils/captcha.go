package utils

import (
	"bytes"
	"encoding/base64"
	"sync"
	"time"

	"github.com/dchest/captcha"
)

// captchaValueStore 存储验证码的实际值（用于测试环境）
type captchaValueStore struct {
	sync.RWMutex
	values map[string]captchaValue
}

type captchaValue struct {
	value     string
	expiresAt time.Time
}

var valueStore = &captchaValueStore{
	values: make(map[string]captchaValue),
}

// 清理过期的验证码值
func (s *captchaValueStore) cleanup() {
	s.Lock()
	defer s.Unlock()
	
	now := time.Now()
	for id, v := range s.values {
		if now.After(v.expiresAt) {
			delete(s.values, id)
		}
	}
}

// 保存验证码值
func (s *captchaValueStore) set(id, value string) {
	s.Lock()
	defer s.Unlock()
	
	s.values[id] = captchaValue{
		value:     value,
		expiresAt: time.Now().Add(5 * time.Minute),
	}
}

// 获取验证码值
func (s *captchaValueStore) get(id string) string {
	s.RLock()
	defer s.RUnlock()
	
	if v, ok := s.values[id]; ok && time.Now().Before(v.expiresAt) {
		return v.value
	}
	return ""
}

// GenerateCaptcha 生成验证码
// 返回: id(验证码ID), base64Img(图片base64), captchaStr(验证码实际值), error
func GenerateCaptcha() (id string, base64Img string, captchaStr string, error error) {
	// 定期清理过期的验证码值
	go valueStore.cleanup()
	
	// 1. 生成验证码数字
	digits := captcha.RandomDigits(4)
	
	// 2. 将数字转换为字符串（这是验证码的实际值）
	captchaStr = ""
	for _, d := range digits {
		captchaStr += string('0' + d)
	}
	
	// 3. 生成验证码ID
	id = captcha.New()
	
	// 4. 使用自定义存储，将我们的数字存储进去
	// 这样 captcha.WriteImage 和 captcha.VerifyString 都会使用相同的数字
	store := captcha.NewMemoryStore(captcha.CollectNum, captcha.Expiration)
	store.Set(id, digits)
	captcha.SetCustomStore(store)
	
	// 5. 将实际值也存储到我们的 valueStore（用于测试环境返回）
	valueStore.set(id, captchaStr)

	// 6. 生成验证码图片（会使用我们设置的 digits）
	var buf bytes.Buffer
	if err := captcha.WriteImage(&buf, id, 120, 40); err != nil {
		return "", "", "", err
	}

	// 7. 转换为 base64
	base64Img = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	return id, base64Img, captchaStr, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id, digits string) bool {
	return captcha.VerifyString(id, digits)
}

// GetCaptchaValue 获取验证码的实际值（仅用于测试）
func GetCaptchaValue(id string) string {
	return valueStore.get(id)
}
