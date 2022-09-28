package utils

import (
	"github.com/dgrijalva/jwt-go"
	"lessons/pkg/common/config"
	"lessons/pkg/entity"
	"time"
)

func GenerateJWTToken(user *entity.User, cfg *config.Config) (string, error) {
	expiresAt := time.Now().Add(time.Duration(cfg.JWTLifeTime) * time.Second)
	claims := &jwt.StandardClaims{
		Subject:   user.ID.String(),
		ExpiresAt: expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeJWTToken(token string, cfg *config.Config) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
