package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mustaphalimar/asmar/internal/app"
	"github.com/mustaphalimar/asmar/internal/db"
	"github.com/mustaphalimar/asmar/internal/env"
)

const version = "0.0.1"

// dbConfig holds database connection settings
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Database configuration
	dbCfg := dbConfig{
		addr:         env.GetString("DATABASE_URL", "postgresql://postgres:admin@localhost/asmar?sslmode=disable"),
		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	// Connect to the database
	sqlDB, err := db.New(dbCfg.addr, dbCfg.maxOpenConns, dbCfg.maxIdleConns, dbCfg.maxIdleTime)
	if err != nil {
		log.Panic(err.Error())
	}
	defer sqlDB.Close()
	log.Println("Database connection pool established.")

	// Application configuration
	appConfig := app.Config{
		Addr: env.GetString("ADDR", ":8080"),
		Env:  env.GetString("ENV", "development"),
	}

	// Create application instance
	application := app.NewApplication(appConfig, sqlDB)

	// Start the HTTP server (defined in api.go)
	if err := serve(application); err != nil {
		log.Fatal(err)
	}
}
