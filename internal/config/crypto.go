package config

import (
	"encoding/base64"
)

// 加密盐值
const encryptionSalt = "mgg_salt_"

// Encrypt 加密字符串（使用 salt + base64）
func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// 添加盐值后进行 base64 编码
	salted := encryptionSalt + plaintext
	encoded := base64.StdEncoding.EncodeToString([]byte(salted))

	return encoded, nil
}

// Decrypt 解密字符串
func Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		// 解码失败，可能是明文密码（向后兼容）
		return ciphertext, nil
	}

	// 检查并移除盐值
	decodedStr := string(decoded)
	if len(decodedStr) > len(encryptionSalt) && decodedStr[:len(encryptionSalt)] == encryptionSalt {
		return decodedStr[len(encryptionSalt):], nil
	}

	// 没有盐值前缀，可能是旧数据，直接返回
	return ciphertext, nil
}
