package db

import (
	"fmt"
	"log"
	"os"
	"pulse/internal/config"
	"pulse/internal/db/models"
	"time"

	charm_log "github.com/charmbracelet/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	// Configure GORM logger
	logLevel := logger.Info
	if cfg.Environment == "production" {
		logLevel = logger.Error
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Log queries slower than 1 second
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true, // Enable prepared statement cache
	})
	if err != nil {
		charm_log.Fatalf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		charm_log.Fatalf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	// Enable UUID extension
	// db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	db.Exec("DEALLOCATE ALL")
	// Run migrations
	if err := autoMigrate(db); err != nil {
		charm_log.Fatalf("failed to run migrations: %v", err)
	}
	charm_log.Info("Database connection established & Migrations run")

	return db
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Chain{},
		&models.TransactionCategory{},
		&models.Wallet{},
		&models.Subwallet{},
		&models.Token{},
		&models.SubwalletToken{},
		&models.HistoricalPrice{},
		&models.Transaction{},
		&models.Snapshot{},
		&models.TokenSnapshot{},
		&models.Alert{},           // depends on User and Token
		&models.PortfolioMetric{}, // depends on User
	)
}
