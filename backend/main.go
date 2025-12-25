package main

import (
	"log"
	"net/http"
	"tpt-titan/backend/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Connect to database
	if err := config.ConnectDatabase(&cfg.Database); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer config.CloseDatabase()

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Basic routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "TPT Titan Backend API",
			"version": "0.1.0",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"database": "connected",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"service": "TPT Titan API",
				"version": "0.1.0",
			})
		})

		// TODO: Add authentication routes
		// TODO: Add user management routes
		// TODO: Add file management routes
		// TODO: Add email routes
	}

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting TPT Titan Backend on %s...", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
