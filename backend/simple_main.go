package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:4173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Basic routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "TPT Titan Backend API",
			"version": "0.1.0",
			"status":  "running",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "TPT Titan API",
			"version": "0.1.0",
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

		// Auth routes - mock responses
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"token": "mock-jwt-token",
					"user": gin.H{
						"id": 1,
						"email": "user@example.com",
						"name": "Demo User",
					},
				})
			})
			authGroup.POST("/register", func(c *gin.Context) {
				c.JSON(http.StatusCreated, gin.H{
					"message": "User registered successfully",
					"user": gin.H{
						"id": 1,
						"email": "user@example.com",
					},
				})
			})
		}

		// Protected routes (simplified - no auth required for demo)
		api.GET("/documents", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"documents": []gin.H{
					{
						"id": 1,
						"title": "Welcome to TPT Titan",
						"type": "text",
						"created_at": "2024-01-01T00:00:00Z",
					},
				},
			})
		})

		api.GET("/ai/models", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"models": []gin.H{
					{
						"id": "llama2",
						"name": "Llama 2",
						"provider": "Ollama",
						"status": "available",
					},
				},
			})
		})

		api.GET("/contacts", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"contacts": []gin.H{
					{
						"id": 1,
						"name": "John Doe",
						"email": "john@example.com",
					},
				},
			})
		})

		api.GET("/calendars", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"calendars": []gin.H{
					{
						"id": 1,
						"name": "Personal",
						"color": "#4285f4",
					},
				},
			})
		})

		api.GET("/emails", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"emails": []gin.H{
					{
						"id": 1,
						"subject": "Welcome to TPT Titan",
						"from": "system@tpt-titan.com",
						"received_at": "2024-01-01T00:00:00Z",
					},
				},
			})
		})

		api.GET("/chat/rooms", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"rooms": []gin.H{
					{
						"id": 1,
						"name": "General",
						"type": "channel",
					},
				},
			})
		})
	}

	// Start server
	serverAddr := "0.0.0.0:8080"
	log.Printf("Starting TPT Titan Backend on %s...", serverAddr)
	log.Printf("API available at http://localhost:8080/api/v1/health")
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
