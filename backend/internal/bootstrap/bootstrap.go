package bootstrap

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stonith404/pocket-id/backend/internal/job"
	"github.com/stonith404/pocket-id/backend/internal/service"
)

func Bootstrap() {
	db := newDatabase()
	appConfigService := service.NewAppConfigService(db)

	initApplicationImages()
	job.RegisterJobs(db)
	initRouter(db, appConfigService)
}
