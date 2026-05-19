package configs

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	JWTSecret                 string
	JWTExpirationHours        int
	JWTRefreshSecret          string
	JWTRefreshExpirationHours int

	APIKey     string
	DevKey     string
	AllowedIPs []string

	RedisHost   string
	RedisPort   string
	RedisPass   string
	RedisPrefix string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	jwtExp, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	jwtRefreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRATION_HOURS", "168"))

	AppConfig = &Config{
		AppPort:                   getEnv("APP_PORT", "8080"),
		DBHost:                    getEnv("DB_HOST", "127.0.0.1"),
		DBPort:                    getEnv("DB_PORT", "3306"),
		DBUser:                    getEnv("DB_USER", "root"),
		DBPass:                    getEnv("DB_PASSWORD", ""),
		DBName:                    getEnv("DB_NAME", "djm_db"),
		JWTSecret:                 getEnv("JWT_SECRET", "secret"),
		JWTExpirationHours:        jwtExp,
		JWTRefreshSecret:          getEnv("JWT_REFRESH_SECRET", "refresh_secret"),
		JWTRefreshExpirationHours: jwtRefreshExp,
		APIKey:                    getEnv("API_KEY", ""),
		DevKey:                    getEnv("DEV_KEY", ""),
		AllowedIPs:                parseCSV(getEnv("ALLOWED_IPS", "127.0.0.1")),
		RedisHost:                 getEnv("REDIS_HOST", "localhost"),
		RedisPort:                 getEnv("REDIS_PORT", "6379"),
		RedisPass:                 getEnv("REDIS_PASSWORD", ""),
		RedisPrefix:               getEnv("REDIS_PREFIX", "djm_api:"),
	}

	// Validate critical security configuration after loading
	if AppConfig.APIKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}
	if AppConfig.DevKey == "" {
		log.Fatal("DEV_KEY environment variable is required")
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseCSV(value string) []string {
	if value == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	result := []string{}
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
