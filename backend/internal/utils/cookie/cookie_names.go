package cookie

import (
	"github.com/stonith404/pocket-id/backend/internal/common"
	"strings"
)

var AccessTokenCookieName = "__Host-access_token"
var SessionIdCookieName = "__Host-session"

func init() {
	if strings.HasPrefix(common.EnvConfig.AppURL, "http://") {
		AccessTokenCookieName = "access_token"
		SessionIdCookieName = "session"
	}
}
