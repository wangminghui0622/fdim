package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// VerifyPassword 验证密码是否匹配
// plainPassword: 用户输入的明文密码
// hashedPassword: 数据库中存储的加密密码
func VerifyPassword(plainPassword, hashedPassword string) bool {
	// 使用 SHA256 加密明文密码
	hash := sha256.Sum256([]byte(plainPassword))
	encrypted := hex.EncodeToString(hash[:])
	
	// 比较加密后的密码
	return encrypted == hashedPassword
}

// HashPassword 加密密码（用于注册时）
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("密码不能为空")
	}
	
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:]), nil
}

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return errors.New("密码长度不能少于6位")
	}
	
	if len(password) > 32 {
		return errors.New("密码长度不能超过32位")
	}
	
	return nil
}
