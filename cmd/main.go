package main

import (
	"log"
	"test-service/internal/config"
	"test-service/internal/delivery/http"
	"test-service/internal/repository"
	"test-service/internal/usecase"
	"test-service/pkg/database"
	"test-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize logger
	logr := logger.NewLogger()

	// Connect database
	db, err := database.NewTiDBConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	// Init layer dependencies
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, logr)
	userHandler := http.NewUserHandler(userUsecase)

	// Setup router
	r := gin.Default()

	http.RegisterUserRoutes(r, userHandler)

	logr.Infof("Starting user-service on port %s 🚀", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
