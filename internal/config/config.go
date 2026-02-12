package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration

	// JWT
	JWTSecret             string
	JWTRefreshSecret      string
	JWTAccessTokenExpire  time.Duration
	JWTRefreshTokenExpire time.Duration

	// Encryption
	EncryptionKey string

	// Server
	AppUrl     string
	ServerPort string
	ServerHost string
	ServerMode string

	// CORS
	CORSAllowOrigins     string
	CORSAllowCredentials bool
	CORSAllowMethods     string
	CORSAllowHeaders     string

	// Logging
	LogLevel  string
	LogFormat string

	// Rate Limiting
	RateLimitEnabled    bool
	RateLimitRequests   int
	RateLimitTimeWindow time.Duration

	// Auto migration & seeding
	AutoMigrate bool
	AutoSeed    bool
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		// Database
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "gin_database"),
		DBMaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME", 300)) * time.Second,

		// JWT
		JWTSecret:             getEnv("JWT_SECRET", "very-secret-key-change-in-production"),
		JWTRefreshSecret:      getEnv("JWT_REFRESH_SECRET", "very-secret-refresh-key-change-in-production"),
		JWTAccessTokenExpire:  time.Duration(getEnvAsInt("JWT_ACCESS_TOKEN_EXPIRE", 15)) * time.Minute,
		JWTRefreshTokenExpire: time.Duration(getEnvAsInt("JWT_REFRESH_TOKEN_EXPIRE", 10080)) * time.Minute,

		// Encryption
		EncryptionKey: getEnv("ENCRYPTION_KEY", "default_32_byte_key_1234567890!@"),

		// Server
		AppUrl:     getEnv("APP_URL", "http://localhost:8080"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		ServerMode: getEnv("SERVER_MODE", "debug"),

		// CORS
		CORSAllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
		CORSAllowMethods:     getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,PATCH,OPTIONS"),
		CORSAllowHeaders:     getEnv("CORS_ALLOW_HEADERS", "Content-Type,Authorization,Accept,Origin,X-Requested-With"),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		// Rate Limiting
		RateLimitEnabled:    getEnvAsBool("RATE_LIMIT_ENABLED", false),
		RateLimitRequests:   getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitTimeWindow: time.Duration(getEnvAsInt("RATE_LIMIT_TIME_WINDOW", 3600)) * time.Second,

		// Auto migration settings
		AutoMigrate: getEnvAsBool("AUTO_MIGRATE", true),
		AutoSeed:    getEnvAsBool("AUTO_SEED", true),
	}
}

// Helper functions untuk parsing environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
