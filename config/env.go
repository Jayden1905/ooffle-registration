package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHosts            string
	Port                   string
	DBUser                 string
	DBPasswd               string
	DBAddr                 string
	DBName                 string
	DBHost                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	ISProduction           bool
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHosts: getEnv("PUBLIC_HOSTS", "http://localhost"),
		Port:        getEnv("PORT", "8080"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPasswd:    getEnv("DB_PASSWD", "root"),
		DBHost:      getEnv("DB_HOST", ""),
		DBAddr: fmt.Sprintf(
			"%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306"),
		),
		DBName:                 getEnv("DB_NAME", "event"),
		JWTSecret:              getEnv("JWT_SECRET", "not-secret-anymore?"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
		ISProduction:           getEnvAsBool("IS_PRODUCTION", false),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}

		return b
	}

	return fallback
}
