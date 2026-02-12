package main

import (
	"log"
	"u_kom_be/cmd/server/routes"
	"u_kom_be/internal/config"
	"u_kom_be/internal/converter"
	"u_kom_be/internal/database"
	"u_kom_be/internal/handler"
	"u_kom_be/internal/middleware"
	"u_kom_be/internal/repository"
	"u_kom_be/internal/service"
	"u_kom_be/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Server holds the application dependencies
type Server struct {
	Config                    *config.Config
	Router                    *gin.Engine
	UserHandler               *handler.UserHandler
	AuthHandler               *handler.AuthHandler
	RoleHandler               *handler.RoleHandler
	PermissionHandler         *handler.PermissionHandler
	StudentHandler            *handler.StudentHandler
	ParentHandler             *handler.ParentHandler
	GuardianHandler           *handler.GuardianHandler
	EmployeeHandler           *handler.EmployeeHandler
	DashboardHandler          *handler.DashboardHandler
	AcademicYearHandler       *handler.AcademicYearHandler
	ClassroomHandler          *handler.ClassroomHandler
	SubjectHandler            *handler.SubjectHandler
	TeachingAssignmentHandler *handler.TeachingAssignmentHandler
	ScheduleHandler           *handler.ScheduleHandler
	AttendanceHandler         *handler.AttendanceHandler
	AuthService               service.AuthService
}

// NewServer creates a new server instance with all dependencies
func NewServer() *Server {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	baseURL := cfg.AppUrl

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
	guardianRepo := repository.NewGuardianRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	academicYearRepo := repository.NewAcademicYearRepository(db)
	classroomRepo := repository.NewClassroomRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	teachingAssignmentRepo := repository.NewTeachingAssignmentRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)

	// Initialize utils
	encryptionUtil, err := utils.NewEncryptionUtil(cfg.EncryptionKey)
	if err != nil {
		log.Fatal("Failed to create encryption util:", err)
	}

	// Initialize converters
	parentConverter := converter.NewParentConverter(encryptionUtil)
	guardianConverter := converter.NewGuardianConverter(encryptionUtil)
	studentConverter := converter.NewStudentConverter(encryptionUtil, parentConverter, baseURL)
	employeeConverter := converter.NewEmployeeConverter(encryptionUtil)

	// Initialize services
	userService := service.NewUserService(userRepo, roleRepo, permissionRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	parentService := service.NewParentService(
		parentRepo,
		userRepo,
		encryptionUtil,
		parentConverter,
	)
	guardianService := service.NewGuardianService(
		guardianRepo,
		userRepo,
		encryptionUtil,
		guardianConverter,
	)
	employeeService := service.NewEmployeeService(
		employeeRepo,
		userRepo,
		encryptionUtil,
		employeeConverter,
	)
	authService := service.NewAuthService(
		userRepo,
		cfg.JWTSecret,
		cfg.JWTRefreshSecret,
		cfg.JWTAccessTokenExpire,
		cfg.JWTRefreshTokenExpire,
	)
	studentService := service.NewStudentService(
		studentRepo,
		parentRepo,
		guardianRepo,
		userRepo,
		encryptionUtil,
		studentConverter,
	)
	dashboardService := service.NewDashboardService(dashboardRepo)
	academicYearService := service.NewAcademicYearService(academicYearRepo, db)
	classroomService := service.NewClassroomService(
		classroomRepo,
		academicYearRepo,
		employeeRepo,
		studentRepo,
		db,
	)
	subjectService := service.NewSubjectService(subjectRepo)
	teachingAssignmentService := service.NewTeachingAssignmentService(
		teachingAssignmentRepo,
		classroomRepo,
		subjectRepo,
		employeeRepo,
	)
	scheduleService := service.NewScheduleService(scheduleRepo, teachingAssignmentRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo, scheduleRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	studentHandler := handler.NewStudentHandler(studentService)
	parentHandler := handler.NewParentHandler(parentService)
	guardianHandler := handler.NewGuardianHandler(guardianService)
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	academicYearHandler := handler.NewAcademicYearHandler(academicYearService)
	classroomHandler := handler.NewClassroomHandler(classroomService)
	subjectHandler := handler.NewSubjectHandler(subjectService)
	teachingAssignmentHandler := handler.NewTeachingAssignmentHandler(teachingAssignmentService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)

	// Setup router with middleware
	router := setupRouter(cfg, authService)

	return &Server{
		Config:                    cfg,
		Router:                    router,
		UserHandler:               userHandler,
		AuthHandler:               authHandler,
		RoleHandler:               roleHandler,
		PermissionHandler:         permissionHandler,
		StudentHandler:            studentHandler,
		ParentHandler:             parentHandler,
		GuardianHandler:           guardianHandler,
		EmployeeHandler:           employeeHandler,
		DashboardHandler:          dashboardHandler,
		AcademicYearHandler:       academicYearHandler,
		ClassroomHandler:          classroomHandler,
		SubjectHandler:            subjectHandler,
		TeachingAssignmentHandler: teachingAssignmentHandler,
		ScheduleHandler:           scheduleHandler,
		AttendanceHandler:         attendanceHandler,
		AuthService:               authService,
	}
}

// setupRouter configures the router with middleware
func setupRouter(cfg *config.Config, authService service.AuthService) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Rate limiting middleware
	//if cfg.RateLimitEnabled {
	//	router.Use(middleware.RateLimitMiddleware(cfg))
	//}

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
		s.GuardianHandler,
		s.EmployeeHandler,
		s.DashboardHandler,
		s.AcademicYearHandler,
		s.ClassroomHandler,
		s.SubjectHandler,
		s.TeachingAssignmentHandler,
		s.ScheduleHandler,
		s.AttendanceHandler,
	)

	// Start server
	serverAddress := s.Config.ServerHost + ":" + s.Config.ServerPort
	log.Printf("Server starting on %s in %s mode", serverAddress, s.Config.ServerMode)

	return s.Router.Run(serverAddress)
}
