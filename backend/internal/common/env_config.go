package common

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	"log"
)

type EnvConfigSchema struct {
	AppEnv             string `env:"APP_ENV"`
	AppURL             string `env:"PUBLIC_APP_URL"`
	DBPath             string `env:"DB_PATH"`
	UploadPath         string `env:"UPLOAD_PATH"`
	Port               string `env:"BACKEND_PORT"`
	Host               string `env:"HOST"`
	EmailTemplatesPath string `env:"EMAIL_TEMPLATES_PATH"`
}

var EnvConfig = &EnvConfigSchema{
	AppEnv:             "production",
	DBPath:             "data/pocket-id.db",
	UploadPath:         "data/uploads",
	AppURL:             "http://localhost",
	Port:               "8080",
	Host:               "localhost",
	EmailTemplatesPath: "./email-templates",
}

func init() {
	if err := env.ParseWithOptions(EnvConfig, env.Options{}); err != nil {
		log.Fatal(err)
	}
}
