package main

import (
	"belajar-golang/cmd/server/routes"
	"belajar-golang/internal/config"
	"belajar-golang/internal/converter"
	"belajar-golang/internal/database"
	"belajar-golang/internal/handler"
	"belajar-golang/internal/middleware"
	"belajar-golang/internal/repository"
	"belajar-golang/internal/service"
	"belajar-golang/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Server holds the application dependencies
type Server struct {
	Config            *config.Config
	Router            *gin.Engine
	UserHandler       *handler.UserHandler
	AuthHandler       *handler.AuthHandler
	RoleHandler       *handler.RoleHandler
	PermissionHandler *handler.PermissionHandler
	StudentHandler    *handler.StudentHandler
	ParentHandler     *handler.ParentHandler
	AuthService       service.AuthService
}

// NewServer creates a new server instance with all dependencies
func NewServer() *Server {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode dari config
	gin.SetMode(cfg.ServerMode)

	// Initialize database
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repository
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	parentRepo := repository.NewParentRepository(db)

	// Initialize utils
	encryptionUtil, err := utils.NewEncryptionUtil(cfg.EncryptionKey)
	if err != nil {
		log.Fatal("Failed to create encryption util:", err)
	}

	// Initialize converters
	studentConverter := converter.NewStudentConverter(encryptionUtil)
	parentConverter := converter.NewParentConverter(encryptionUtil)

	// Initialize services
	userService := service.NewUserService(userRepo, roleRepo, permissionRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	studentService := service.NewStudentService(studentRepo, encryptionUtil, studentConverter)
	parentService := service.NewParentService(parentRepo, encryptionUtil, parentConverter)
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
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	studentHandler := handler.NewStudentHandler(studentService)
	parentHandler := handler.NewParentHandler(parentService)

	// Setup router with middleware
	router := setupRouter(cfg, authService)

	return &Server{
		Config:            cfg,
		Router:            router,
		UserHandler:       userHandler,
		AuthHandler:       authHandler,
		RoleHandler:       roleHandler,
		PermissionHandler: permissionHandler,
		StudentHandler:    studentHandler,
		ParentHandler:     parentHandler,
		AuthService:       authService,
	}
}

// setupRouter configures the router with middleware
func setupRouter(cfg *config.Config, authService service.AuthService) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Rate limiting middleware
	if cfg.RateLimitEnabled {
		router.Use(middleware.RateLimitMiddleware(cfg))
	}

	return router
}

// Start runs the HTTP server
func (s *Server) Start() error {
	// Setup routes
	routes.SetupRoutes(
		s.Router,
		s.AuthHandler,
		s.UserHandler,
		s.AuthService,
		s.RoleHandler,
		s.PermissionHandler,
		s.StudentHandler,
		s.ParentHandler,
	)

	// Start server
	serverAddress := s.Config.ServerHost + ":" + s.Config.ServerPort
	log.Printf("Server starting on %s in %s mode", serverAddress, s.Config.ServerMode)

	return s.Router.Run(serverAddress)
}
