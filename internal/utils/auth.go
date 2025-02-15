package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct {
	SecretKey []byte
	TokenExp  time.Duration
}

// GenerateToken generates a JWT token with the user ID as part of the claims
func (a *AuthConfig) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(a.TokenExp).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.SecretKey)
}

// VerifyToken verifies a token JWT validate
func (a *AuthConfig) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.SecretKey, nil
	})

	// Check for errors
	if err != nil {
		return nil, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	return claims, nil
}
