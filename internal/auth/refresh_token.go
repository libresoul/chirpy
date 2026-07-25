package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	data := make([]byte, 32)
	rand.Read(data)

	hexStr := hex.EncodeToString(data)
	return hexStr
}
