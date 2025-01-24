package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddAccessTokenCookie(c *gin.Context, sessionDurationInMinutes string, token string) {
	sessionDurationInMinutesParsed, _ := strconv.Atoi(sessionDurationInMinutes)
	maxAge := sessionDurationInMinutesParsed * 60
	c.SetCookie("__Host-access_token", token, maxAge, "/", "", true, true)
}
