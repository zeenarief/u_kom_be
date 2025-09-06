package main

import (
	"belajar-golang/internal/config"
	"belajar-golang/internal/database"
	"belajar-golang/internal/handler"
	"belajar-golang/internal/middleware"
	"belajar-golang/internal/repository"
	"belajar-golang/internal/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.ServerMode)

	// Initialize database
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(
		userRepo,
		cfg.JWTSecret,
		cfg.JWTRefreshSecret,
		cfg.JWTAccessTokenExpire,
		cfg.JWTRefreshTokenExpire,
	)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	// Setup router
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery()) // Recovery middleware
	router.Use(middleware.CORSMiddleware())

	// Rate limiting middleware
	if cfg.RateLimitEnabled {
		// Pilihan 1: Simple (recommended untuk kebanyakan case)
		//router.Use(middleware.SimpleRateLimitMiddleware(cfg))

		// Atau Pilihan 2: Advanced (butuh install dependency)
		router.Use(middleware.RateLimitMiddleware(cfg))
	}

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", userHandler.CreateUser)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", authHandler.RefreshToken)
		public.GET("/health", func(c *gin.Context) {
			handler.SuccessResponse(c, "Server is healthy and running", gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		})
	}

	// Protected routes (require JWT)
	protected := public.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// User routes
		protected.GET("/users", userHandler.GetAllUsers)
		protected.POST("/users", userHandler.CreateUser)
		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
		protected.POST("/users/:id/change-password", userHandler.ChangePassword)

		// Auth routes
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/profile", userHandler.GetProfile)
	}

	// Start server
	serverAddress := cfg.ServerHost + ":" + cfg.ServerPort
	log.Printf("Server starting on %s in %s mode", serverAddress, cfg.ServerMode)

	if err := router.Run(serverAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
