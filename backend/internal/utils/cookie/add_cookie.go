package cookie

import (
	"github.com/gin-gonic/gin"
)

func AddAccessTokenCookie(c *gin.Context, maxAgeInSeconds int, token string) {
	c.SetCookie(AccessTokenCookieName, token, maxAgeInSeconds, "/", "", true, true)
}

func AddSessionIdCookie(c *gin.Context, maxAgeInSeconds int, sessionID string) {
	c.SetCookie(SessionIdCookieName, sessionID, maxAgeInSeconds, "/", "", true, true)
}
