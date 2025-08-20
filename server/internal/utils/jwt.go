package utils

import (
	"context"
	"fmt"

	"server/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func ValidateToken(ctx context.Context, tokenString string) (*CustomClaims, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
