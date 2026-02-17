package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret            string
	JWTExpireHours       int
	JWTRefreshExpireDays int

	ServerPort string
	GinMode    string

	UploadPath  string
	MaxFileSize int64

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	FrontendURL string
}

var AppConfig *Config

func LoadConfig() error {
	godotenv.Load()

	jwtExpireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	jwtRefreshExpireDays, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_DAYS", "7"))
	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "104857600"), 10, 64)
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "8r44r699r"),
		DBName:     getEnv("DB_NAME", "ysmmc"),

		JWTSecret:            getEnv("JWT_SECRET", "ysmmc-jwt-secret-key"),
		JWTExpireHours:       jwtExpireHours,
		JWTRefreshExpireDays: jwtRefreshExpireDays,

		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),

		UploadPath:  getEnv("UPLOAD_PATH", "../uploads"),
		MaxFileSize: maxFileSize,

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     smtpPort,
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
