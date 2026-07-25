package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	headerVal := headers.Get("Authorization")
	if headerVal == "" {
		return "", fmt.Errorf("Missing auth header")
	}

	if !strings.HasPrefix(headerVal, "ApiKey ") && len(headerVal) < 15 {
		return "", fmt.Errorf("Missing API key")
	}

	key := strings.Fields(headerVal)[1]
	return key, nil
}
