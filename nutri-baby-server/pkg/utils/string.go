package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"strings"
)

// GenerateShortCode 生成6位短码（字母数字组合）
// 使用 Base36 编码（0-9, A-Z），避免混淆字符（如0和O，1和I）
func GenerateShortCode() (string, error) {
	// 字符集：去除容易混淆的字符 0,O,1,I,l
	const charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	const codeLength = 6

	result := make([]byte, codeLength)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < codeLength; i++ {
		// 生成随机索引
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

// GenerateToken 生成指定长度的随机 token (十六进制字符串)
// length 参数指定字节数，实际生成的字符串长度为 length * 2
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ParseScene 解析小程序码 scene 参数
// 格式: "key1=val1&key2=val2"
// 返回: map[string]string
func ParseScene(scene string) map[string]string {
	result := make(map[string]string)
	if scene == "" {
		return result
	}

	pairs := strings.Split(scene, "&")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}

	return result
}

// BuildScene 构建小程序码 scene 参数
// 格式: "key1=val1&key2=val2"
// 注意: 总长度不能超过32个字符
func BuildScene(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	var parts []string
	for key, value := range params {
		parts = append(parts, key+"="+value)
	}

	return strings.Join(parts, "&")
}
