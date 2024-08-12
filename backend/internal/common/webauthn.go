package common

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"golang-rest-api-template/internal/utils"
	"log"
	"time"
)

var (
	WebAuthn *webauthn.WebAuthn
	err      error
)

func init() {
	config := &webauthn.Config{
		RPDisplayName: DbConfig.AppName.Value,
		RPID:          utils.GetHostFromURL(EnvConfig.AppURL),
		RPOrigins:     []string{EnvConfig.AppURL},
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    time.Second * 60,
				TimeoutUVD: time.Second * 60,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    time.Second * 60,
				TimeoutUVD: time.Second * 60,
			},
		},
	}

	if WebAuthn, err = webauthn.New(config); err != nil {
		log.Fatal(err)
	}
}
