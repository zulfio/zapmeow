package database

import (
	"zapmeow/pkg/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Database interface {
	RunMigrate(dst ...interface{}) error
	Client() *gorm.DB
}

type database struct {
	client *gorm.DB
}

func NewDatabase(databasePath string) *database {
	client, err := gorm.Open(sqlite.Open(databasePath+"?_busy_timeout=5000"), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Fatal("Error creating gorm database. ", err)
	}

	// Get the underlying *sql.DB object
	sqlDB, err := client.DB()
	if err != nil {
		logger.Fatal("Error getting underlying database connection. ", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(1) // This is crucial for SQLite to prevent "database is locked" errors
	sqlDB.SetMaxIdleConns(1)

	return &database{
		client: client,
	}
}

func (d *database) RunMigrate(dst ...interface{}) error {
	return d.client.AutoMigrate(dst...)
}

func (d *database) Client() *gorm.DB {
	return d.client
}