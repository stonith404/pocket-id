package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils/cookie"
	"strings"
)

type JwtAuthMiddleware struct {
	jwtService            *service.JwtService
	ignoreUnauthenticated bool
}

func NewJwtAuthMiddleware(jwtService *service.JwtService, ignoreUnauthenticated bool) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{jwtService: jwtService, ignoreUnauthenticated: ignoreUnauthenticated}
}

func (m *JwtAuthMiddleware) Add(adminOnly bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the cookie or the Authorization header
		token, err := c.Cookie(cookie.AccessTokenCookieName)
		if err != nil {
			authorizationHeaderSplitted := strings.Split(c.GetHeader("Authorization"), " ")
			if len(authorizationHeaderSplitted) == 2 {
				token = authorizationHeaderSplitted[1]
			} else if m.ignoreUnauthenticated {
				c.Next()
				return
			} else {
				c.Error(&common.NotSignedInError{})
				c.Abort()
				return
			}
		}

		claims, err := m.jwtService.VerifyAccessToken(token)
		if err != nil && m.ignoreUnauthenticated {
			c.Next()
			return
		} else if err != nil {
			c.Error(&common.NotSignedInError{})
			c.Abort()
			return
		}

		// Check if the user is an admin
		if adminOnly && !claims.IsAdmin {
			c.Error(&common.MissingPermissionError{})
			c.Abort()
			return
		}

		c.Set("userID", claims.Subject)
		c.Set("userIsAdmin", claims.IsAdmin)
		c.Next()
	}
}
