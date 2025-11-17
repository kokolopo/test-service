package main

import (
	"log"
	"test-service/internal/config"
	httpDelivery "test-service/internal/delivery/http"
	"test-service/internal/repository"
	"test-service/internal/usecase"
	"test-service/pkg/database"
	"test-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	logr := logger.NewLogger()

	db, err := database.NewTiDBConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, logr)
	userHandler := httpDelivery.NewUserHandler(userUsecase)

	r := gin.New()

	// access.log middleware
	r.Use(logger.GinAccessLogger(logr))

	httpDelivery.RegisterUserRoutes(r, userHandler)

	logr.Infof("Starting user-service on port %s 🚀", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
