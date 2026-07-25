package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(user_id uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	issuedTime := time.Now()
	expiredOn := issuedTime.Add(expiresIn)

	iat := jwt.NewNumericDate(issuedTime)
	exp := jwt.NewNumericDate(expiredOn)

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		Subject:   user_id.String(),
		ExpiresAt: exp,
		IssuedAt:  iat,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var user_id uuid.UUID

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("Invalid token")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	user_id, err = uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid user id")
	}
	return user_id, nil
}
