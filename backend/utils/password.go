package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
)

// GenerateSalt 生成随机盐值
func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// HashPassword 使用 MD5 + 盐值加密密码
func HashPassword(password, salt string) string {
	h := md5.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyPassword 验证密码
func VerifyPassword(password, salt, hashedPassword string) bool {
	return HashPassword(password, salt) == hashedPassword
}
