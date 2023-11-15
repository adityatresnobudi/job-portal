package db

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	log.Println(dsn)

	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{SlowThreshold: time.Second, LogLevel: logger.Info, Colorful: true})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger, TranslateError: true})
	if err != nil {
		return nil, err
	}
	return db, nil
}
