package database

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	cfg := config.AppConfig
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Model{},
		&model.Favorite{},
		&model.Announcement{},
		&model.Session{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrated successfully")
	return nil
}

func Seed() error {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin := model.User{
		ID:            uuid.New(),
		Email:         "admin@ysmmc.local",
		PasswordHash:  string(passwordHash),
		Username:      "admin",
		Role:          "admin",
		ProfileStatus: "approved",
		EmailVerified: true,
	}

	if err := DB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	announcement := model.Announcement{
		ID:      uuid.New(),
		Title:   "Welcome to YSM Model Station",
		Content: "This is a model sharing platform where you can upload, download and share your model works.",
		IsActive: true,
	}

	if err := DB.Create(&announcement).Error; err != nil {
		return fmt.Errorf("failed to create announcement: %w", err)
	}

	log.Println("Database seeded successfully")
	return nil
}
