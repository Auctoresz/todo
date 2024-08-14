package config

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func OpenDB() *sql.DB {
	var db *sql.DB
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbOptions := os.Getenv("DB_OPTIONS")

	// Check if environment variables are set
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Println("One or more environment variables are not set")
		return nil
	}

	// Create the data source name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbOptions)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
