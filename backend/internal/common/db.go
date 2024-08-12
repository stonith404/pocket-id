package common

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	connectDatabase()
	sqlDb, err := DB.DB()
	if err != nil {
		log.Fatal("failed to get sql db", err)
	}
	driver, err := sqlite3.WithInstance(sqlDb, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("failed to create migration instance", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to run migrations", err)
	}
}

func connectDatabase() {
	var database *gorm.DB
	var err error

	dbPath := EnvConfig.DBPath
	if EnvConfig.AppEnv == "test" {
		dbPath = "file::memory:?cache=shared"
	}

	for i := 1; i <= 3; i++ {
		database, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
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

	DB = database
}

func getLogger() logger.Interface {
	isProduction := EnvConfig.AppEnv == "production"

	var logLevel logger.LogLevel
	if isProduction {
		logLevel = logger.Error
	} else {
		logLevel = logger.Info
	}

	// Create the GORM logger
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
