package jwtour

import (
	"denettest/internal/config"
	"denettest/internal/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewToken(user domain.User, duration time.Duration) (string, error) {
	secretKey := config.GetConfig().SecretKey
	claims := jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(duration).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}
