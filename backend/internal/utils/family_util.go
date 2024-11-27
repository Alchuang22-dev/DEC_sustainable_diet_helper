package utils

import (
	"encoding/hex"
)

// 生成一个唯一的家庭 Token
func GenerateFamilyToken() string {
    tokenBytes := make([]byte, 32)
    return hex.EncodeToString(tokenBytes)
}