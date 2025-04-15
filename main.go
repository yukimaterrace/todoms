package main

import (
	"fmt"
	"log"

	"github.com/yukimaterrace/todoms/config"
	"github.com/yukimaterrace/todoms/controller"
	"github.com/yukimaterrace/todoms/repository"
	"github.com/yukimaterrace/todoms/service"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Connect to database
	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authConfig := &config.AuthConfig{
		JWTSecret:          repository.GetEnvOrDefault("JWT_SECRET", "your-secret-key-change-me-in-production"),
		AccessTokenExpiry:  config.DefaultAccessTokenExpiry,
		RefreshTokenExpiry: config.DefaultRefreshTokenExpiry,
	}
	userService := service.NewUserService(userRepo, logger)
	authService := service.NewJWTAuthService(userRepo, authConfig, logger)

	// Setup Echo using controller package
	e := controller.SetupEcho(userService, authService)

	// Start server
	port := repository.GetEnvOrDefault("PORT", "8080")
	fmt.Printf("Server is running on port %s...\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
