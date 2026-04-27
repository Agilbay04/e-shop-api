package config

import (
	"e-shop-api/internal/constant"
	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/pkg/querytracker"
	"e-shop-api/internal/pkg/util"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(constant.PostgresDSN,
		host, user, password, dbName, port)

	var db *gorm.DB
	var err error

	// Retries to connect to database
	err = util.AutoRetry(func() error {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            return err
        }
        
        // Verify database connection by pinging
        sqlDB, err := db.DB()
        if err != nil {
            return err
        }
        return sqlDB.Ping()

	})

	// Return error if failed connect to database after retries
	if err != nil {
		logger.L.Fatal("Failed connect to database:", zap.Error(err))
		panic(fmt.Sprintf("Failed connect to database: %v", err))
	}

	// Setup database pooling
	setupDatabasePooling(db)

	// Register slow query tracker for tracking query and slow query
	if err := querytracker.Register(db); err != nil {
		logger.L.Warn("Failed to register slow query plugin:", zap.Error(err))
	}

	// Success connect to database
	logger.L.Info("Connected to database!")
	return db
}

func setupDatabasePooling(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logger.L.Fatal("Failed to setup database pooling:", zap.Error(err))
		panic("Failed to setup database pooling!")
	}

	maxIdle := util.GetEnvInt("DB_MAX_IDLE_CONNS", constant.DBMaxIdleConns)
	maxOpen := util.GetEnvInt("DB_MAX_OPEN_CONNS", constant.DBMaxOpenConns)
	maxLifetimeMinutes := util.GetEnvTime("DB_CONN_MAX_LIFETIME", constant.DBConnMaxLifetime)
	maxIdleMinutes := util.GetEnvTime("DB_CONN_MAX_IDLETIME", constant.DBConnMaxIdleTime)

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(maxLifetimeMinutes)
	sqlDB.SetConnMaxIdleTime(maxIdleMinutes)
}
