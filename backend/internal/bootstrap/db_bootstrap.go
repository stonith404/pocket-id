package bootstrap

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func newDatabase() (db *gorm.DB) {
	db, err := connectDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	sqlDb, err := db.DB()
	sqlDb.SetMaxOpenConns(1)
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	driver, err := sqlite3.WithInstance(sqlDb, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("failed to create migration instance: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	return db
}

func connectDatabase() (db *gorm.DB, err error) {
	dbPath := common.EnvConfig.DBPath

	// Use in-memory database for testing
	if common.EnvConfig.AppEnv == "test" {
		dbPath = "file::memory:?cache=shared"
	}

	for i := 1; i <= 3; i++ {
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			TranslateError: true,
			Logger:         getLogger(),
		})
		if err == nil {
			break
		} else {
			log.Printf("Attempt %d: Failed to initialize database. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}

	return db, err
}

func getLogger() logger.Interface {
	isProduction := common.EnvConfig.AppEnv == "production"

	var logLevel logger.LogLevel
	if isProduction {
		logLevel = logger.Error
	} else {
		logLevel = logger.Info
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: isProduction,
			ParameterizedQueries:      isProduction,
			Colorful:                  !isProduction,
		},
	)
}
