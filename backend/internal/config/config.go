package config

import (
	"fmt"
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

	UploadPath          string
	MaxFileSize         int64
	MaxDiskUsage        int
	EnableDatePartition bool

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	FrontendURL    string
	AllowedOrigins []string
}

var AppConfig *Config

func LoadConfig() error {
	godotenv.Load()

	jwtExpireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	jwtRefreshExpireDays, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_DAYS", "7"))
	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "104857600"), 10, 64)
	maxDiskUsage, _ := strconv.Atoi(getEnv("MAX_DISK_USAGE", "90"))
	enableDatePartition := getEnv("ENABLE_DATE_PARTITION", "false") == "true"
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "ysmmc"),

		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTExpireHours:       jwtExpireHours,
		JWTRefreshExpireDays: jwtRefreshExpireDays,

		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),

		UploadPath:          getEnv("UPLOAD_PATH", "./uploads"),
		MaxFileSize:         maxFileSize,
		MaxDiskUsage:        maxDiskUsage,
		EnableDatePartition: enableDatePartition,

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     smtpPort,
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),

		FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:5173"),
		AllowedOrigins: parseAllowedOrigins(getEnv("ALLOWED_ORIGINS", "http://localhost:5173")),
	}

	return nil
}

func Validate() error {
	if AppConfig.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD environment variable is required")
	}
	if AppConfig.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	if len(AppConfig.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}
	return nil
}

func parseAllowedOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}
	result := []string{}
	start := 0
	for i := 0; i <= len(origins); i++ {
		if i == len(origins) || origins[i] == ',' {
			origin := origins[start:i]
			if origin != "" {
				result = append(result, origin)
			}
			start = i + 1
		}
	}
	return result
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
