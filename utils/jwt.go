package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("AllYourBase")

type JWTClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func GenerateJWT(userId string) (string, error) {
	// Create the Claims
	claims := JWTClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(1440) * time.Minute).Unix(),
			Issuer:    "Beinan's Auth Service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ParseJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method.")
		}
		return signingKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user"].(string), nil
	} else {
		return "", errors.New("Invalid token.")
	}
}
