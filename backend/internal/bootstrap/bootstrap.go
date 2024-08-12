package bootstrap

import (
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang-rest-api-template/internal/common"
	"golang-rest-api-template/internal/common/middleware"
	"golang-rest-api-template/internal/handler"
	"golang-rest-api-template/internal/job"
	"golang-rest-api-template/internal/utils"
	"golang.org/x/time/rate"
	"log"
	"os"
	"time"
)

func Bootstrap() {
	common.InitDatabase()
	common.InitDbConfig()
	initApplicationImages()
	job.RegisterJobs()
	initRouter()
}

func initRouter() {
	switch common.EnvConfig.AppEnv {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "development":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	}

	r := gin.Default()

	r.Use(gin.Logger())

	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(time.Second), 60))

	apiGroup := r.Group("/api")
	handler.RegisterRoutes(apiGroup)
	handler.RegisterOIDCRoutes(apiGroup)
	handler.RegisterUserRoutes(apiGroup)
	handler.RegisterConfigurationRoutes(apiGroup)
	if common.EnvConfig.AppEnv != "production" {
		handler.RegisterTestRoutes(apiGroup)
	}

	baseGroup := r.Group("/")
	handler.RegisterWellKnownRoutes(baseGroup)

	if err := r.Run(common.EnvConfig.Host + ":" + common.EnvConfig.Port); err != nil {
		log.Fatal(err)
	}

}

func initApplicationImages() {
	dirPath := common.EnvConfig.UploadPath + "/application-images"

	files, err := os.ReadDir(dirPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error reading directory: %v", err)
	}

	// Skip if files already exist
	if len(files) > 1 {
		return
	}

	// Copy files from source to destination
	err = utils.CopyDirectory("./images", dirPath)
	if err != nil {
		log.Fatalf("Error copying directory: %v", err)
	}
}
