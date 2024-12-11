package bootstrap

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	sqliteMigrate "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"gorm.io/driver/postgres"
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
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// Choose the correct driver for the database provider
	var driver database.Driver
	switch common.EnvConfig.DbProvider {
	case common.DbProviderSqlite:
		driver, err = sqliteMigrate.WithInstance(sqlDb, &sqliteMigrate.Config{})
	case common.DbProviderPostgres:
		driver, err = postgresMigrate.WithInstance(sqlDb, &postgresMigrate.Config{})
	default:
		log.Fatalf("unsupported database provider: %s", common.EnvConfig.DbProvider)
	}
	if err != nil {
		log.Fatalf("failed to create migration driver: %v", err)
	}

	// Run migrations
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/"+string(common.EnvConfig.DbProvider),
		"pocket-id", driver,
	)
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
	var dialector gorm.Dialector

	// Choose the correct database provider
	switch common.EnvConfig.DbProvider {
	case common.DbProviderSqlite:
		dialector = sqlite.Open(common.EnvConfig.SqliteDBPath)
	case common.DbProviderPostgres:
		dialector = postgres.Open(common.EnvConfig.PostgresConnectionString)
	default:
		return nil, fmt.Errorf("unsupported database provider: %s", common.EnvConfig.DbProvider)
	}

	for i := 1; i <= 3; i++ {
		db, err = gorm.Open(dialector, &gorm.Config{
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
