package config

import (
	"e-shop-api/internal/pkg/util"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbName, port)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	setupDatabasePooling(db)

	fmt.Println("Connected to database!")
	return db
}

func setupDatabasePooling(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to setup database pooling!")
	}

	maxIdle := util.GetEnvInt("DB_MAX_IDLE_CONNS", 10)
	maxOpen := util.GetEnvInt("DB_MAX_OPEN_CONNS", 100)
	maxLifetimeMinutes := util.GetEnvInt("DB_CONN_MAX_LIFETIME", 60)
	maxIdleMinutes := util.GetEnvInt("DB_CONN_MAX_IDLETIME", 15)

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetimeMinutes) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(maxIdleMinutes) * time.Minute)
}
