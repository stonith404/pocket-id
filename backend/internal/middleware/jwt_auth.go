package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"net/http"
	"strings"
)

type JwtAuthMiddleware struct {
	jwtService *service.JwtService
}

func NewJwtAuthMiddleware(jwtService *service.JwtService) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{jwtService: jwtService}
}

func (m *JwtAuthMiddleware) Add(adminOnly bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the cookie or the Authorization header
		token, err := c.Cookie("access_token")
		if err != nil {
			authorizationHeaderSplitted := strings.Split(c.GetHeader("Authorization"), " ")
			if len(authorizationHeaderSplitted) == 2 {
				token = authorizationHeaderSplitted[1]
			} else {
				utils.HandlerError(c, http.StatusUnauthorized, "You're not signed in")
				c.Abort()
				return
			}
		}

		claims, err := m.jwtService.VerifyAccessToken(token)
		if err != nil {
			utils.HandlerError(c, http.StatusUnauthorized, "You're not signed in")
			c.Abort()
			return
		}

		// Check if the user is an admin
		if adminOnly && !claims.IsAdmin {
			utils.HandlerError(c, http.StatusForbidden, "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Set("userID", claims.Subject)
		c.Set("userIsAdmin", claims.IsAdmin)
		c.Next()
	}
}
