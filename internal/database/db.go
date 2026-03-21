package database

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get underlying DB")
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(25)                 // Max open connections
	sqlDB.SetMaxIdleConns(5)                  // Max idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Max connection lifetime

	DB = db
}
