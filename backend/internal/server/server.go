package server

import (
	"fmt"
	"log"
	"net/http"
	"tpt-titan-simple/backend/config"
	"tpt-titan-simple/backend/middleware"
	"tpt-titan-simple/backend/routes"
	"tpt-titan-simple/backend/routes/auth"
	"tpt-titan-simple/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server represents the TPT Titan server
type Server struct {
	config   *config.Config
	router   *gin.Engine
	database *gorm.DB
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// Initialize sets up all server components
func (s *Server) Initialize() error {
	// Set Gin mode
	gin.SetMode(s.config.Server.Mode)

	// Connect to database
	if err := config.ConnectDatabase(&s.config.Database); err != nil {
		return err
	}
	s.database = config.GetDatabase()

	// Initialize cache service if enabled (optional for small teams)
	var cacheService *services.CacheService
	if s.config.Redis.Enabled {
		redisURL := fmt.Sprintf("redis://%s:%s", s.config.Redis.Host, s.config.Redis.Port)
		var err error
		cacheService, err = services.NewCacheService(redisURL)
		if err != nil {
			log.Printf("Failed to connect to Redis, continuing without cache: %v", err)
			cacheService = nil
		}
	} else {
		log.Println("Redis is disabled (REDIS_ENABLED=false). Running without cache.")
	}

	// Validate required configuration
	if err := s.validateConfig(); err != nil {
		return err
	}

	// Initialize auth package with JWT secret
	auth.InitAuth(s.config.JWT.Secret)

	// Initialize auth service
	authService := services.NewAuthService(s.database, s.config.JWT.Secret, cacheService)

	// Initialize AI service
	routes.InitAIService(&s.config.AI)

	// Initialize Speech service
	routes.InitSpeechService(&s.config.Speech)

	// Initialize Workflow service
	routes.InitWorkflowService()

	// Initialize Workflow AI service
	routes.InitWorkflowAIService(routes.GetAIService())

	// Initialize contact service
	routes.InitContactService(s.database)

	// Initialize calendar service
	routes.InitCalendarService(s.database)

	// Initialize email service
	routes.InitEmailService(s.database)

	// Initialize chat service
	routes.InitChatService(s.database)

	// Initialize voice service
	routes.InitVoiceService()

	// Initialize database optimizer
	dbOptimizer := services.NewDatabaseOptimizer(s.database)
	dbOptimizer.OptimizeConnectionPool()

	// Initialize monitoring service
	monitoringService := services.NewMonitoringService(cacheService, dbOptimizer)

	// Initialize WebSocket hub
	hub := routes.InitWebSocketHub(routes.GetChatService())

	// Initialize security middleware
	securityMiddleware := middleware.NewSecurityMiddleware()

	// Initialize Gin router
	s.router = gin.Default()

	// Global middleware - applied to all routes
	s.router.Use(securityMiddleware.RequestIDMiddleware())
	s.router.Use(securityMiddleware.CORSMiddleware())
	s.router.Use(securityMiddleware.SecurityHeadersMiddleware())
	s.router.Use(securityMiddleware.RateLimitMiddleware())
	s.router.Use(securityMiddleware.InputValidationMiddleware())
	s.router.Use(securityMiddleware.SQLInjectionProtectionMiddleware())
	s.router.Use(securityMiddleware.AuditMiddleware())

	// Add CORS middleware
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())

	// Setup routes
	s.setupRoutes(authService, monitoringService, hub)

	return nil
}

// setupRoutes configures all the API routes
func (s *Server) setupRoutes(authService *services.AuthService, monitoringService *services.MonitoringService, hub *routes.WebSocketHub) {
	// Basic routes
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "TPT Titan Backend API",
			"version": "0.1.0",
		})
	})

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"database": "connected",
		})
	})

	// API routes
	api := s.router.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("db", s.database)
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
		protected.Use(auth.AuthMiddleware(s.database))
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

			// Basic Form routes
			formGroup := protected.Group("/forms")
			{
				formGroup.GET("", routes.GetForms)
				formGroup.POST("", routes.CreateForm)
				formGroup.GET("/:id", routes.GetForm)
				formGroup.PUT("/:id", routes.UpdateForm)
				formGroup.DELETE("/:id", routes.DeleteForm)
				formGroup.GET("/:id/responses", routes.GetFormResponses)
				formGroup.POST("/:id/responses", routes.SubmitFormResponse)
			}

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

				// AI processing routes
				documentGroup.POST("/upload", routes.UploadDocumentWithAI)
				documentGroup.POST("/:id/process", routes.ProcessDocumentWithAI)
				documentGroup.GET("/:id/analysis", routes.GetDocumentAnalysis)
				documentGroup.GET("/:id/status", routes.GetDocumentProcessingStatus)
				documentGroup.GET("/:id/analyses", routes.GetDocumentAnalyses)
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

			// Speech routes (TTS/STT)
			speechGroup := protected.Group("/speech")
			{
				speechGroup.GET("/models", routes.GetSpeechModels)
				speechGroup.POST("/models", routes.CreateSpeechModel)
				speechGroup.POST("/tts", routes.TextToSpeech)
				speechGroup.POST("/stt", routes.SpeechToText)
				speechGroup.GET("/requests/:requestId", routes.GetSpeechRequestStatus)
				speechGroup.GET("/settings", routes.GetSpeechSettings)
				speechGroup.PUT("/settings", routes.UpdateSpeechSettings)
				speechGroup.GET("/history", routes.GetSpeechHistory)
			}

			// AI Settings routes
			settingsGroup := protected.Group("/settings")
			{
				settingsGroup.GET("/ai", routes.GetAISettings)
				settingsGroup.PUT("/ai", routes.UpdateAISettings)
				settingsGroup.GET("/speech", routes.GetSpeechSettings)
				settingsGroup.PUT("/speech", routes.UpdateSpeechSettings)
				settingsGroup.GET("/ai/usage", routes.GetAIUsageStats)
				settingsGroup.POST("/ai/test-key", routes.TestAPIKey)
			}

			// Workflow automation routes
			workflowGroup := protected.Group("/workflows")
			{
				workflowGroup.GET("", routes.GetWorkflows)
				workflowGroup.POST("", routes.CreateWorkflow)
				workflowGroup.GET("/:id", routes.GetWorkflow)
				workflowGroup.PUT("/:id", routes.UpdateWorkflow)
				workflowGroup.DELETE("/:id", routes.DeleteWorkflow)
				workflowGroup.POST("/:id/execute", routes.ExecuteWorkflow)
				workflowGroup.GET("/:id/executions", routes.GetWorkflowExecutions)
				workflowGroup.PUT("/:id/nodes", routes.UpdateWorkflowNodes)
				workflowGroup.PUT("/:id/connections", routes.UpdateWorkflowConnections)
			}

			// Workflow execution routes
			protected.GET("/executions/:executionId", routes.GetWorkflowExecution)

			// Workflow templates
			templateGroup := protected.Group("/workflow-templates")
			{
				templateGroup.GET("", routes.GetWorkflowTemplates)
				templateGroup.POST("/:templateId/create", routes.CreateWorkflowFromTemplate)
			}

			// AI Workflow Intelligence routes
			aiWorkflowGroup := protected.Group("/ai/workflows")
			{
				aiWorkflowGroup.GET("/insights", routes.GetWorkflowInsights)
				aiWorkflowGroup.GET("/usage-analysis", routes.AnalyzeWorkflowUsage)
				aiWorkflowGroup.GET("/template-recommendations", routes.GetSmartTemplateRecommendations)
				aiWorkflowGroup.GET("/:workflowId/optimization", routes.OptimizeWorkflow)
				aiWorkflowGroup.GET("/predictions", routes.PredictWorkflows)
			}

			// Integration connectors
			protected.GET("/connectors", routes.GetIntegrationConnectors)

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

			// Voice notes and annotations routes
			voiceGroup := protected.Group("/voice")
			{
				// Voice notes
				voiceGroup.GET("/notes", routes.GetVoiceNotes)
				voiceGroup.POST("/notes", routes.CreateVoiceNote)
				voiceGroup.GET("/notes/:id", routes.GetVoiceNote)
				voiceGroup.PUT("/notes/:id", routes.UpdateVoiceNote)
				voiceGroup.DELETE("/notes/:id", routes.DeleteVoiceNote)

				// Voice annotations
				voiceGroup.GET("/annotations", routes.GetVoiceAnnotations)
				voiceGroup.POST("/annotations", routes.CreateVoiceAnnotation)
				voiceGroup.GET("/annotations/:id", routes.GetVoiceAnnotation)
				voiceGroup.DELETE("/annotations/:id", routes.DeleteVoiceAnnotation)
			}

			// Database table editor routes (spreadsheet-like database editing)
			databaseGroup := protected.Group("/database")
			{
				databaseGroup.GET("/tables", routes.GetDatabaseTables)
				databaseGroup.GET("/tables/:table/info", routes.GetTableInfo)
				databaseGroup.GET("/tables/:table/data", routes.GetTableData)
				databaseGroup.PUT("/tables/:table/records/:id", routes.UpdateTableRecord)
				databaseGroup.POST("/tables/:table/records", routes.CreateTableRecord)
				databaseGroup.DELETE("/tables/:table/records/:id", routes.DeleteTableRecord)
			}

			// TODO: Add file management routes
			// TODO: Add form routes
			// TODO: Add task routes
		}
	}

	// Prometheus metrics endpoint
	api.GET("/metrics", monitoringService.PrometheusMetricsHandler())
}

// Start starts the server
func (s *Server) Start() error {
	serverAddr := s.config.Server.Host + ":" + s.config.Server.Port
	log.Printf("Starting TPT Titan Backend on %s...", serverAddr)
	return s.router.Run(serverAddr)
}

// GetRouter returns the Gin router (useful for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// GetDatabase returns the database instance
func (s *Server) GetDatabase() *gorm.DB {
	return s.database
}

// Close closes the server and cleans up resources
func (s *Server) Close() {
	config.CloseDatabase()
}

// validateConfig checks that all required configuration is set
func (s *Server) validateConfig() error {
	if s.config.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required. Please set a secure secret key")
	}
	
	if len(s.config.JWT.Secret) < 32 {
		log.Println("WARNING: JWT_SECRET should be at least 32 characters for security")
	}
	
	return nil
}
