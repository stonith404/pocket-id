package common

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	"golang-rest-api-template/internal/model"
	"log"
	"reflect"
)

type EnvConfigSchema struct {
	AppEnv     string `env:"APP_ENV"`
	AppURL     string `env:"PUBLIC_APP_URL"`
	DBPath     string `env:"DB_PATH"`
	UploadPath string `env:"UPLOAD_PATH"`
	Port       string `env:"BACKEND_PORT"`
	Host       string `env:"HOST"`
}

var EnvConfig = &EnvConfigSchema{
	AppEnv:     "production",
	DBPath:     "data/pocket-id.db",
	UploadPath: "data/uploads",
	AppURL:     "http://localhost",
	Port:       "8080",
	Host:       "localhost",
}

var DbConfig = NewDefaultDbConfig()

func NewDefaultDbConfig() model.ApplicationConfiguration {
	return model.ApplicationConfiguration{
		AppName: model.ApplicationConfigurationVariable{
			Key:      "appName",
			Type:     "string",
			IsPublic: true,
			Value:    "Pocket ID",
		},
		BackgroundImageType: model.ApplicationConfigurationVariable{
			Key:        "backgroundImageType",
			Type:       "string",
			IsInternal: true,
			Value:      "jpg",
		},
		LogoImageType: model.ApplicationConfigurationVariable{
			Key:        "logoImageType",
			Type:       "string",
			IsInternal: true,
			Value:      "svg",
		},
	}
}

// LoadDbConfigFromDb refreshes the database configuration by loading the current values
// from the database and updating the DbConfig struct.
func LoadDbConfigFromDb() error {
	dbConfigReflectValue := reflect.ValueOf(&DbConfig).Elem()

	for i := 0; i < dbConfigReflectValue.NumField(); i++ {
		dbConfigField := dbConfigReflectValue.Field(i)
		currentConfigVar := dbConfigField.Interface().(model.ApplicationConfigurationVariable)
		var storedConfigVar model.ApplicationConfigurationVariable
		if err := DB.First(&storedConfigVar, "key = ?", currentConfigVar.Key).Error; err != nil {
			return err
		}

		dbConfigField.Set(reflect.ValueOf(storedConfigVar))
	}

	return nil
}

// InitDbConfig creates the default configuration values in the database if they do not exist,
// updates existing configurations if they differ from the default, and deletes any configurations
// that are not in the default configuration.
func InitDbConfig() {
	// Reflect to get the underlying value of DbConfig and its default configuration
	dbConfigReflectValue := reflect.ValueOf(&DbConfig).Elem()
	defaultDbConfig := NewDefaultDbConfig()
	defaultConfigReflectValue := reflect.ValueOf(&defaultDbConfig).Elem()
	defaultKeys := make(map[string]struct{})

	// Iterate over the fields of DbConfig
	for i := 0; i < dbConfigReflectValue.NumField(); i++ {
		dbConfigField := dbConfigReflectValue.Field(i)
		currentConfigVar := dbConfigField.Interface().(model.ApplicationConfigurationVariable)
		defaultConfigVar := defaultConfigReflectValue.Field(i).Interface().(model.ApplicationConfigurationVariable)
		defaultKeys[currentConfigVar.Key] = struct{}{}

		var storedConfigVar model.ApplicationConfigurationVariable
		if err := DB.First(&storedConfigVar, "key = ?", currentConfigVar.Key).Error; err != nil {
			// If the configuration does not exist, create it
			if err := DB.Create(&defaultConfigVar).Error; err != nil {
				log.Fatalf("Failed to create default configuration: %v", err)
			}
			dbConfigField.Set(reflect.ValueOf(defaultConfigVar))
			continue
		}

		// Update existing configuration if it differs from the default
		if storedConfigVar.Type != defaultConfigVar.Type || storedConfigVar.IsPublic != defaultConfigVar.IsPublic || storedConfigVar.IsInternal != defaultConfigVar.IsInternal {
			storedConfigVar.Type = defaultConfigVar.Type
			storedConfigVar.IsPublic = defaultConfigVar.IsPublic
			storedConfigVar.IsInternal = defaultConfigVar.IsInternal
			if err := DB.Save(&storedConfigVar).Error; err != nil {
				log.Fatalf("Failed to update configuration: %v", err)
			}
		}

		// Set the value in DbConfig
		dbConfigField.Set(reflect.ValueOf(storedConfigVar))
	}

	// Delete any configurations not in the default keys
	var allConfigVars []model.ApplicationConfigurationVariable
	if err := DB.Find(&allConfigVars).Error; err != nil {
		log.Fatalf("Failed to retrieve existing configurations: %v", err)
	}

	for _, config := range allConfigVars {
		if _, exists := defaultKeys[config.Key]; !exists {
			if err := DB.Delete(&config).Error; err != nil {
				log.Fatalf("Failed to delete outdated configuration: %v", err)
			}
		}
	}
}

func init() {
	if err := env.ParseWithOptions(EnvConfig, env.Options{}); err != nil {
		log.Fatal(err)
	}
}
