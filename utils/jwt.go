package utils

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("AllYourBase")

var logger = DefaultLogger

type AppClaims struct {
	UserId      string   `json:"user"`
	IsAdmin     bool     `json:"isAdmin,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type JWTClaims struct {
	AppClaims
	jwt.StandardClaims
}

func GenerateJWT(appClaims AppClaims) (string, error) {
	// Create the Claims
	claims := JWTClaims{
		appClaims,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(1440) * time.Minute).Unix(),
			Issuer:    "Beinan's Auth Service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method.")
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("Invalid token.")
	}
}

const authContextKey = "AUTH_OBJECT_CONTEXT_KEY"

func AuthAttach(ctx context.Context, token string) context.Context {
	authObject, err := ParseJWT(token)
	if err != nil {
		logger.Debugf("Parsing JWT error: %v", err)
	}
	return context.WithValue(
		ctx,
		authContextKey,
		authObject,
	)
}

func GetAuthObject(ctx context.Context) *JWTClaims {
	authObj := ctx.Value(authContextKey).(*JWTClaims)
	logger.Debugf("Auth object getting from ctx: %v", authObj)
	return authObj
}
