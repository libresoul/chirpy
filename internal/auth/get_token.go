package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	headerVal := headers.Get("Authorization")
	if headerVal == "" {
		return "", fmt.Errorf("Missing auth header")
	}

	if !strings.HasPrefix(headerVal, "Bearer ") && len(headerVal) < 20 {
		return "", fmt.Errorf("Missing token")
	}

	token := strings.Fields(headerVal)[1]
	return token, nil
}
