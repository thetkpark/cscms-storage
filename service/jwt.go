package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtManager struct {
	secret []byte
}

func (j *JwtManager) GenerateUserJWT(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userId,
	})

	return token.SignedString(j.secret)
}

func (j *JwtManager) ValidateUserJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return j.secret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.StandardClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", fmt.Errorf("claims is not valid")
}