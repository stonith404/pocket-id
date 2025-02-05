package cookie

import (
	"strings"

	"github.com/pocket-id/pocket-id/backend/internal/common"
)

var AccessTokenCookieName = "__Host-access_token"
var SessionIdCookieName = "__Host-session"

func init() {
	if strings.HasPrefix(common.EnvConfig.AppURL, "http://") {
		AccessTokenCookieName = "access_token"
		SessionIdCookieName = "session"
	}
}
