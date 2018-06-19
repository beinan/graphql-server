package utils

import (
	"os"
)

const (
	PROD    = "prod"
	STAGING = "staging"
	DEVEL   = "devel"
)

var Env = GetEnv("APP_ENV", DEVEL)

func IsProd() bool {
	return Env == PROD
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
