package handler

import (
	"github.com/gin-gonic/gin"
	"golang-rest-api-template/internal/common"
	"golang-rest-api-template/internal/utils"
	"net/http"
)

func RegisterWellKnownRoutes(group *gin.RouterGroup) {
	group.GET("/.well-known/jwks.json", jwks)
	group.GET("/.well-known/openid-configuration", openIDConfiguration)
}

func jwks(c *gin.Context) {
	jwk, err := common.GetJWK()
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"keys": []interface{}{jwk}})
}

func openIDConfiguration(c *gin.Context) {
	appUrl := common.EnvConfig.AppURL
	config := map[string]interface{}{
		"issuer":                                appUrl,
		"authorization_endpoint":                appUrl + "/authorize",
		"token_endpoint":                        appUrl + "/api/oidc/token",
		"jwks_uri":                              appUrl + "/.well-known/jwks.json",
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"claims_supported":                      []string{"sub", "given_name", "family_name", "email", "preferred_username"},
		"response_types_supported":              []string{"code", "id_token"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	}
	c.JSON(http.StatusOK, config)
}
