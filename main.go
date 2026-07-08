package main

import (
	"log"
	"os"
	"tapinvest_api/config"
	"tapinvest_api/handlers"
	"tapinvest_api/repository"
	"tapinvest_api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database Connection
	config.ConnectDatabase()
	defer config.DB.Close()

	// Initialize Repositories
	bondRepo := repository.NewBondRepository(config.DB)
	wishlistRepo := repository.NewWishlistRepository(config.DB)

	// Initialize Handlers
	bondHandler := handlers.NewBondHandler(bondRepo)
	wishlistHandler := handlers.NewWishlistHandler(wishlistRepo)

	// Setup Gin Router
	router := gin.Default()
	
	// Register Routes
	routes.RegisterRoutes(router, bondHandler, wishlistHandler)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}