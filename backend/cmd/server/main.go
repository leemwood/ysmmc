package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/router"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	if err := database.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	gin.SetMode(config.AppConfig.GinMode)

	r := gin.Default()

	router.Setup(r)

	log.Printf("Server starting on port %s...", config.AppConfig.ServerPort)
	if err := r.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
