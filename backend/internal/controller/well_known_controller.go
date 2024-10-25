package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"net/http"
)

func NewWellKnownController(group *gin.RouterGroup, jwtService *service.JwtService) {
	wkc := &WellKnownController{jwtService: jwtService}
	group.GET("/.well-known/jwks.json", wkc.jwksHandler)
	group.GET("/.well-known/openid-configuration", wkc.openIDConfigurationHandler)
}

type WellKnownController struct {
	jwtService *service.JwtService
}

func (wkc *WellKnownController) jwksHandler(c *gin.Context) {
	jwk, err := wkc.jwtService.GetJWK()
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"keys": []interface{}{jwk}})
}

func (wkc *WellKnownController) openIDConfigurationHandler(c *gin.Context) {
	appUrl := common.EnvConfig.AppURL
	config := map[string]interface{}{
		"issuer":                                appUrl,
		"authorization_endpoint":                appUrl + "/authorize",
		"token_endpoint":                        appUrl + "/api/oidc/token",
		"userinfo_endpoint":                     appUrl + "/api/oidc/userinfo",
		"jwks_uri":                              appUrl + "/.well-known/jwks.json",
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"claims_supported":                      []string{"sub", "given_name", "family_name", "name", "email", "email_verified", "preferred_username"},
		"response_types_supported":              []string{"code", "id_token"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	}
	c.JSON(http.StatusOK, config)
}
