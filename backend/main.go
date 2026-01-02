package main

import (
	"log"
	"net/http"
	"tpt-titan/backend/config"
	"tpt-titan/backend/middleware"
	"tpt-titan/backend/routes"
	"tpt-titan/backend/routes/auth"
	"tpt-titan/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	// Initialize cache service first (needed for auth service)
	cacheService, err := services.NewCacheService("redis://localhost:6379")
	if err != nil {
		log.Printf("Failed to connect to Redis, continuing without cache: %v", err)
		cacheService = nil
	}

	// Initialize auth service
	authService := services.NewAuthService(config.GetDatabase(), cfg.JWT.Secret, cacheService)

	// Initialize AI service
	routes.InitAIService(&cfg.AI)

	// Initialize contact service
	routes.InitContactService(config.GetDatabase())

	// Initialize calendar service
	routes.InitCalendarService(config.GetDatabase())

	// Initialize email service
	routes.InitEmailService(config.GetDatabase())

	// Initialize chat service
	routes.InitChatService(config.GetDatabase())

	// Initialize database optimizer
	dbOptimizer := services.NewDatabaseOptimizer(config.GetDatabase())
	dbOptimizer.OptimizeConnectionPool()

	// Initialize monitoring service
	monitoringService := services.NewMonitoringService(cacheService, dbOptimizer)

	// Initialize WebSocket hub
	hub := routes.InitWebSocketHub(routes.GetChatService())

	// Initialize security middleware
	securityMiddleware := middleware.NewSecurityMiddleware()

	// Initialize Gin router
	r := gin.Default()

	// Global middleware - applied to all routes
	r.Use(securityMiddleware.RequestIDMiddleware())
	r.Use(securityMiddleware.CORSMiddleware())
	r.Use(securityMiddleware.SecurityHeadersMiddleware())
	r.Use(securityMiddleware.RateLimitMiddleware())
	r.Use(securityMiddleware.InputValidationMiddleware())
	r.Use(securityMiddleware.SQLInjectionProtectionMiddleware())
	r.Use(securityMiddleware.AuditMiddleware())

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

	// Get database instance
	db := config.GetDatabase()

	// API routes
	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"service": "TPT Titan API",
				"version": "0.1.0",
			})
		})

		// Authentication routes (public)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/request-password-reset", auth.RequestPasswordReset)
			authGroup.POST("/reset-password", auth.ResetPassword)
		}

		// Protected routes (require authentication)
		protected := api.Group("/")
		protected.Use(auth.AuthMiddleware(db))
		{
			// User profile routes
			protected.GET("/profile", auth.GetUserProfile)
			protected.PUT("/profile", auth.UpdateUserProfile)
			protected.POST("/logout", auth.Logout)

			// Two-factor authentication routes
			protected.POST("/auth/enable-totp", auth.EnableTOTP)
			protected.POST("/auth/verify-totp", auth.VerifyAndEnableTOTP)
			protected.POST("/auth/disable-totp", auth.DisableTOTP)

			// Password management routes
			protected.POST("/auth/change-password", auth.ChangePassword)

			// Spreadsheet routes - TODO: Implement
			// spreadsheetGroup := protected.Group("/spreadsheets")

			// Advanced Form routes - TODO: Implement
			// formsAdvancedGroup := protected.Group("/forms")

			// Document routes
			documentGroup := protected.Group("/documents")
			{
				documentGroup.GET("", routes.GetDocuments)
				documentGroup.POST("", routes.CreateDocument)
				documentGroup.GET("/:id", routes.GetDocument)
				documentGroup.PUT("/:id", routes.UpdateDocument)
				documentGroup.DELETE("/:id", routes.DeleteDocument)
				documentGroup.GET("/:id/versions", routes.GetDocumentVersions)
				documentGroup.POST("/:id/versions/:version/restore", routes.RestoreDocumentVersion)
			}

			// AI routes
			aiGroup := protected.Group("/ai")
			{
				aiGroup.GET("/models", routes.GetAIModels)
				aiGroup.POST("/models", routes.CreateAIModel)
				aiGroup.POST("/requests", routes.ProcessAIRequest)
				aiGroup.GET("/requests/:requestId", routes.GetAIRequestStatus)
				aiGroup.GET("/ollama/models", routes.ListOllamaModels)
				aiGroup.POST("/ollama/models/:modelName/pull", routes.PullOllamaModel)
				aiGroup.GET("/usage", routes.GetAIUsage)
				aiGroup.GET("/hardware", routes.DetectHardware)
				aiGroup.GET("/recommendations", routes.GetRecommendedModels)
				aiGroup.POST("/setup", routes.SetupRecommendedModels)
				aiGroup.POST("/upgrades/check", routes.CheckForUpgrades)
				aiGroup.GET("/upgrades/history", routes.GetUpgradeHistory)
				aiGroup.POST("/upgrades/apply", routes.ApplyUpgrade)
			}

			// Contact routes
			contactGroup := protected.Group("/contacts")
			{
				contactGroup.GET("", routes.GetContacts)
				contactGroup.POST("", routes.CreateContact)
				contactGroup.GET("/:id", routes.GetContact)
				contactGroup.PUT("/:id", routes.UpdateContact)
				contactGroup.DELETE("/:id", routes.DeleteContact)
				contactGroup.GET("/search", routes.SearchContacts)
			}

			// Calendar routes
			calendarGroup := protected.Group("/calendars")
			{
				calendarGroup.GET("", routes.GetCalendars)
				calendarGroup.POST("", routes.CreateCalendar)
				calendarGroup.GET("/:id", routes.GetCalendar)
				calendarGroup.PUT("/:id", routes.UpdateCalendar)
				calendarGroup.DELETE("/:id", routes.DeleteCalendar)
				calendarGroup.GET("/:id/events", routes.GetCalendarEvents)
			}

			// Event routes
			eventGroup := protected.Group("/events")
			{
				eventGroup.GET("", routes.GetEvents)
				eventGroup.POST("", routes.CreateEvent)
				eventGroup.GET("/:id", routes.GetEvent)
				eventGroup.PUT("/:id", routes.UpdateEvent)
				eventGroup.DELETE("/:id", routes.DeleteEvent)
			}

			// Email routes
			emailGroup := protected.Group("/emails")
			{
				emailGroup.GET("", routes.GetEmails)
				emailGroup.POST("", routes.SendEmail)
				emailGroup.GET("/:id", routes.GetEmail)
				emailGroup.PUT("/:id/read", routes.MarkEmailAsRead)
				emailGroup.PUT("/:id/star", routes.StarEmail)
				emailGroup.PUT("/:id/move", routes.MoveEmailToFolder)
			}

			// Email account routes
			emailAccountGroup := protected.Group("/email-accounts")
			{
				emailAccountGroup.GET("", routes.GetEmailAccounts)
				emailAccountGroup.POST("", routes.CreateEmailAccount)
				emailAccountGroup.GET("/:id", routes.GetEmailAccount)
				emailAccountGroup.PUT("/:id", routes.UpdateEmailAccount)
				emailAccountGroup.DELETE("/:id", routes.DeleteEmailAccount)
				emailAccountGroup.POST("/:accountId/sync", routes.SyncEmails)
			}

			// Email stats
			protected.GET("/email-stats", routes.GetEmailStats)

			// Chat routes
			chatGroup := protected.Group("/chat")
			{
				chatGroup.GET("/rooms", routes.GetChatRooms)
				chatGroup.POST("/rooms", routes.CreateChatRoom)
				chatGroup.GET("/rooms/:id", routes.GetChatRoom)
				chatGroup.POST("/rooms/:id/participants", routes.AddRoomParticipants)
				chatGroup.DELETE("/rooms/:id/leave", routes.LeaveChatRoom)
				chatGroup.GET("/rooms/:id/messages", routes.GetMessages)
				chatGroup.POST("/rooms/:id/messages", routes.SendMessage)
				chatGroup.PUT("/rooms/:id/read", routes.MarkRoomAsRead)
				chatGroup.POST("/messages/:id/reactions", routes.AddReaction)
				chatGroup.DELETE("/messages/:id/reactions", routes.RemoveReaction)
				chatGroup.POST("/direct", routes.CreateDirectMessage)
			}

			// WebSocket route
			protected.GET("/ws", func(c *gin.Context) {
				if hub != nil {
					hub.HandleWebSocket(c)
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket not available"})
				}
			})

			// Monitoring routes
			monitoringGroup := protected.Group("/monitoring")
			{
				monitoringGroup.GET("/metrics", monitoringService.MetricsHandler())
				monitoringGroup.GET("/health", monitoringService.HealthHandler())
				monitoringGroup.GET("/performance", func(c *gin.Context) {
					report, err := monitoringService.GetPerformanceReport()
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, report)
				})
				monitoringGroup.GET("/alerts", func(c *gin.Context) {
					alerts := monitoringService.GetAlerts(50)
					c.JSON(http.StatusOK, gin.H{"alerts": alerts})
				})
			}

			// Math routes - TODO: Implement
			// mathGroup := protected.Group("/math")

			// Document export routes - TODO: Implement
			// exportGroup := protected.Group("/export")

			// Admin routes - TODO: Implement
			// admin := protected.Group("/admin")
			// admin.Use(services.NewAuthService(config.GetDatabase(), cfg.JWT.Secret, cacheService).RoleMiddleware("admin"))

			// Prometheus metrics endpoint
			api.GET("/metrics", monitoringService.PrometheusMetricsHandler())

			// TODO: Add file management routes
			// TODO: Add form routes
			// TODO: Add task routes
		}
	}

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting TPT Titan Backend on %s...", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
