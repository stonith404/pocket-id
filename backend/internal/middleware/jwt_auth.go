package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"net/http"
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
		token, err := c.Cookie("access_token")
		if err != nil {
			authorizationHeaderSplitted := strings.Split(c.GetHeader("Authorization"), " ")
			if len(authorizationHeaderSplitted) == 2 {
				token = authorizationHeaderSplitted[1]
			} else if m.ignoreUnauthenticated {
				c.Next()
				return
			} else {
				utils.CustomControllerError(c, http.StatusUnauthorized, "You're not signed in")
				c.Abort()
				return
			}
		}

		claims, err := m.jwtService.VerifyAccessToken(token)
		if err != nil && m.ignoreUnauthenticated {
			c.Next()
			return
		} else if err != nil {
			utils.CustomControllerError(c, http.StatusUnauthorized, "You're not signed in")
			c.Abort()
			return
		}

		// Check if the user is an admin
		if adminOnly && !claims.IsAdmin {
			utils.CustomControllerError(c, http.StatusForbidden, "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Set("userID", claims.Subject)
		c.Set("userIsAdmin", claims.IsAdmin)
		c.Next()
	}
}
