package main

import (
	"flag"
	"log"
	"os"
	"u_kom_be/internal/config"
	"u_kom_be/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	// Command line flags
	migrate := flag.Bool("migrate", false, "Run database migrations only")
	seed := flag.Bool("seed", false, "Run database seeding only")
	devMigrate := flag.Bool("dev-migrate", false, "Run DEV database migrations (GORM AutoMigrate)")
	migrateSql := flag.Bool("migrate-sql", false, "Run SQL migrations from /migrations folder")
	flag.Parse()

	if *migrateSql {
		runSqlMigrationsOnly()
		return
	}

	if *migrate {
		runMigrationsOnly()
		return
	}

	if *seed {
		runSeedingOnly()
		return
	}

	if *devMigrate {
		runDevMigrationsOnly()
		return
	}

	// Create and start a server
	server := NewServer()

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Fungsi untuk menjalankan migrasi SQL
func runSqlMigrationsOnly() {
	log.Println("Running SQL database migrations...")

	cfg := config.LoadConfig()

	// Panggil fungsi migrasi baru kita
	if err := database.RunSQLMigrations(cfg); err != nil {
		log.Fatal("Failed to run SQL migrations:", err)
	}

	log.Println("SQL Migrations completed successfully")
	os.Exit(0)
}

func runDevMigrationsOnly() {
	log.Println("Running DEV database migrations (GORM AutoMigrate)...")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.AutoMigrateForDev(db); err != nil {
		log.Fatal("Failed to run DEV migrations:", err)
	}

	log.Println("DEV Migrations completed successfully")
	os.Exit(0)
}

func runMigrationsOnly() {
	log.Println("Running database migrations only...")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Migrations completed successfully")
	os.Exit(0)
}

func runSeedingOnly() {
	log.Println("Running database seeding only...")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.SeedData(db); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	log.Println("Seeding completed successfully")
	os.Exit(0)
}
