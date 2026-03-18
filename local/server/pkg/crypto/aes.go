package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// AESEncrypt AES-GCM加密
func AESEncrypt(plaintext []byte, key []byte) (string, error) {
	// 使用SHA256将密钥转换为32字节
	hash := sha256.Sum256(key)
	aesKey := hash[:]

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt AES-GCM解密
func AESDecrypt(ciphertext string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	// 使用SHA256将密钥转换为32字节
	hash := sha256.Sum256(key)
	aesKey := hash[:]

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// GenerateKey 生成随机密钥
func GenerateKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// DeriveKey 从密码派生密钥（简单实现）
func DeriveKey(password string, salt []byte) []byte {
	combined := append([]byte(password), salt...)
	hash := sha256.Sum256(combined)
	return hash[:]
}

// EncryptMessage 加密消息内容
func EncryptMessage(content []byte, conversationKey string) (string, error) {
	return AESEncrypt(content, []byte(conversationKey))
}

// DecryptMessage 解密消息内容
func DecryptMessage(encryptedContent string, conversationKey string) ([]byte, error) {
	return AESDecrypt(encryptedContent, []byte(conversationKey))
}

// GenerateConversationKey 生成会话密钥
// 使用双方用户ID生成确定性密钥
func GenerateConversationKey(userID1, userID2 string) string {
	// 确保顺序一致
	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}
	combined := userID1 + ":" + userID2
	hash := sha256.Sum256([]byte(combined))
	return base64.StdEncoding.EncodeToString(hash[:])
}
