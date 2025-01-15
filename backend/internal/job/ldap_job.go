package job

import (
	"log"

	"github.com/go-co-op/gocron/v2"
	"github.com/stonith404/pocket-id/backend/internal/service"
)

type LdapJobs struct {
	ldapService *service.LdapService
}

func RegisterLdapJobs(ls *service.LdapService) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create a new scheduler: %s", err)
	}

	jobs := &LdapJobs{ldapService: ls}

	registerJob(scheduler, "ClearWebauthnSessions", "*/5 * * * *", jobs.ldapSyncJob)
	scheduler.Start()
}

func (j *LdapJobs) ldapSyncJob() error {
	return j.ldapService.GetLdapUsers()
}
