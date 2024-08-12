package job

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"golang-rest-api-template/internal/common"
	"golang-rest-api-template/internal/model"
	"golang-rest-api-template/internal/utils"
	"log"
	"time"
)

func RegisterJobs() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create a new scheduler: %s", err)
	}

	registerJob(scheduler, "ClearWebauthnSessions", "0 3 * * *", clearWebauthnSessions)
	registerJob(scheduler, "ClearOneTimeAccessTokens", "0 3 * * *", clearOneTimeAccessTokens)
	registerJob(scheduler, "ClearOidcAuthorizationCodes", "0 3 * * *", clearOidcAuthorizationCodes)

	scheduler.Start()
}

func registerJob(scheduler gocron.Scheduler, name string, interval string, job func() error) {

	_, err := scheduler.NewJob(
		gocron.CronJob(interval, false),
		gocron.NewTask(job),
		gocron.WithEventListeners(
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				log.Printf("Job %q run successfully", name)
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				log.Printf("Job %q failed with error: %v", name, err)
			}),
		),
	)

	if err != nil {
		log.Fatalf("Failed to register job %q: %v", name, err)
	}
}

func clearWebauthnSessions() error {
	return common.DB.Delete(&model.WebauthnSession{}, "expires_at < ?", utils.FormatDateForDb(time.Now())).Error
}

func clearOneTimeAccessTokens() error {
	return common.DB.Debug().Delete(&model.OneTimeAccessToken{}, "expires_at < ?", utils.FormatDateForDb(time.Now())).Error
}

func clearOidcAuthorizationCodes() error {
	return common.DB.Delete(&model.OidcAuthorizationCode{}, "expires_at < ?", utils.FormatDateForDb(time.Now())).Error

}
