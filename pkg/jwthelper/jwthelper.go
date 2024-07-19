package jwthelper

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	*jwt.MapClaims
	Username string `json:"username"`
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func GetClaims(authorization string) (jwt.MapClaims, error) {
	tokenString := authorization[len("Bearer "):]
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}
