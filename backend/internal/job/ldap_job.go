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

	registerJob(scheduler, "SyncLdapUsers", "*/5 * * * *", jobs.ldapUserSyncJob)
	registerJob(scheduler, "SyncLdapGroups", "*/5 * * * *", jobs.ldapGroupSyncJob)
	scheduler.Start()
}

func (j *LdapJobs) ldapUserSyncJob() error {
	return j.ldapService.GetLdapUsers()
}

func (j *LdapJobs) ldapGroupSyncJob() error {
	return j.ldapService.GetLdapGroups()
}
