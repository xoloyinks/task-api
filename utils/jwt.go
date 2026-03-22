// utils/jwt.go
package utils

import (
	"fmt"
	"time"

	"task-tracker-api/config"
	"task-tracker-api/models"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClaimsKey contextKey = "claims"

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

var cfg = config.Load()
var secret = cfg.JwtSecret

func GenerateJWT(user *models.User) (string, error) {

	// 1. define the claims (payload)
	claims := Claims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // expires in 24hrs
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // when it was issued
		},
	}

	// 2. create the token with signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. sign the token with your secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error signing token")
	}

	return tokenString, nil

}

// utils/jwt.go
func VerifyJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// make sure signing method is HMAC and not something else
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	// extract and validate claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
