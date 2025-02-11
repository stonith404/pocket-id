package common

import (
	"log"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type DbProvider string

const (
	DbProviderSqlite   DbProvider = "sqlite"
	DbProviderPostgres DbProvider = "postgres"
)

type EnvConfigSchema struct {
	AppEnv                   string     `env:"APP_ENV"`
	AppURL                   string     `env:"PUBLIC_APP_URL"`
	DbProvider               DbProvider `env:"DB_PROVIDER"`
	SqliteDBPath             string     `env:"SQLITE_DB_PATH"`
	PostgresConnectionString string     `env:"POSTGRES_CONNECTION_STRING"`
	UploadPath               string     `env:"UPLOAD_PATH"`
	Port                     string     `env:"BACKEND_PORT"`
	Host                     string     `env:"HOST"`
	MaxMindLicenseKey        string     `env:"MAXMIND_LICENSE_KEY"`
	GeoLiteDBPath            string     `env:"GEOLITE_DB_PATH"`
}

var EnvConfig = &EnvConfigSchema{
	AppEnv:                   "production",
	DbProvider:               "sqlite",
	SqliteDBPath:             "data/pocket-id.db",
	PostgresConnectionString: "",
	UploadPath:               "data/uploads",
	AppURL:                   "http://localhost",
	Port:                     "8080",
	Host:                     "0.0.0.0",
	MaxMindLicenseKey:        "",
	GeoLiteDBPath:            "data/GeoLite2-City.mmdb",
}

func init() {
	if err := env.ParseWithOptions(EnvConfig, env.Options{}); err != nil {
		log.Fatal(err)
	}
	// Validate the environment variables
	if EnvConfig.DbProvider != DbProviderSqlite && EnvConfig.DbProvider != DbProviderPostgres {
		log.Fatal("Invalid DB_PROVIDER value. Must be 'sqlite' or 'postgres'")
	}

	if EnvConfig.DbProvider == DbProviderPostgres && EnvConfig.PostgresConnectionString == "" {
		log.Fatal("Missing POSTGRES_CONNECTION_STRING environment variable")
	}

	if EnvConfig.DbProvider == DbProviderSqlite && EnvConfig.SqliteDBPath == "" {
		log.Fatal("Missing SQLITE_DB_PATH environment variable")
	}
}
