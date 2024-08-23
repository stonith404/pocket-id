package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"github.com/stonith404/pocket-id/backend/internal/middleware"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

func NewOidcController(group *gin.RouterGroup, jwtAuthMiddleware *middleware.JwtAuthMiddleware, fileSizeLimitMiddleware *middleware.FileSizeLimitMiddleware, oidcService *service.OidcService, jwtService *service.JwtService) {
	oc := &OidcController{oidcService: oidcService, jwtService: jwtService}

	group.POST("/oidc/authorize", jwtAuthMiddleware.Add(false), oc.authorizeHandler)
	group.POST("/oidc/authorize/new-client", jwtAuthMiddleware.Add(false), oc.authorizeNewClientHandler)
	group.POST("/oidc/token", oc.createIDTokenHandler)
	group.GET("/oidc/userinfo", oc.userInfoHandler)

	group.GET("/oidc/clients", jwtAuthMiddleware.Add(true), oc.listClientsHandler)
	group.POST("/oidc/clients", jwtAuthMiddleware.Add(true), oc.createClientHandler)
	group.GET("/oidc/clients/:id", oc.getClientHandler)
	group.PUT("/oidc/clients/:id", jwtAuthMiddleware.Add(true), oc.updateClientHandler)
	group.DELETE("/oidc/clients/:id", jwtAuthMiddleware.Add(true), oc.deleteClientHandler)

	group.POST("/oidc/clients/:id/secret", jwtAuthMiddleware.Add(true), oc.createClientSecretHandler)

	group.GET("/oidc/clients/:id/logo", oc.getClientLogoHandler)
	group.DELETE("/oidc/clients/:id/logo", oc.deleteClientLogoHandler)
	group.POST("/oidc/clients/:id/logo", jwtAuthMiddleware.Add(true), fileSizeLimitMiddleware.Add(2<<20), oc.updateClientLogoHandler)
}

type OidcController struct {
	oidcService *service.OidcService
	jwtService  *service.JwtService
}

func (oc *OidcController) authorizeHandler(c *gin.Context) {
	var input dto.AuthorizeOidcClientDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ControllerError(c, err)
		return
	}

	code, err := oc.oidcService.Authorize(input, c.GetString("userID"))
	if err != nil {
		if errors.Is(err, common.ErrOidcMissingAuthorization) {
			utils.CustomControllerError(c, http.StatusForbidden, err.Error())
		} else {
			utils.ControllerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
}

func (oc *OidcController) authorizeNewClientHandler(c *gin.Context) {
	var input dto.AuthorizeOidcClientDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ControllerError(c, err)
		return
	}

	code, err := oc.oidcService.AuthorizeNewClient(input, c.GetString("userID"))
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
}

func (oc *OidcController) createIDTokenHandler(c *gin.Context) {
	var input dto.OidcIdTokenDto

	if err := c.ShouldBind(&input); err != nil {
		utils.ControllerError(c, err)
		return
	}

	clientID := input.ClientID
	clientSecret := input.ClientSecret

	// Client id and secret can also be passed over the Authorization header
	if clientID == "" || clientSecret == "" {
		var ok bool
		clientID, clientSecret, ok = c.Request.BasicAuth()
		if !ok {
			utils.CustomControllerError(c, http.StatusBadRequest, "Client id and secret not provided")
			return
		}
	}

	idToken, accessToken, err := oc.oidcService.CreateTokens(input.Code, input.GrantType, clientID, clientSecret)
	if err != nil {
		if errors.Is(err, common.ErrOidcGrantTypeNotSupported) ||
			errors.Is(err, common.ErrOidcMissingClientCredentials) ||
			errors.Is(err, common.ErrOidcClientSecretInvalid) ||
			errors.Is(err, common.ErrOidcInvalidAuthorizationCode) {
			utils.CustomControllerError(c, http.StatusBadRequest, err.Error())
		} else {
			utils.ControllerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"id_token": idToken, "access_token": accessToken, "token_type": "Bearer"})
}

func (oc *OidcController) userInfoHandler(c *gin.Context) {
	token := strings.Split(c.GetHeader("Authorization"), " ")[1]
	jwtClaims, err := oc.jwtService.VerifyOauthAccessToken(token)
	if err != nil {
		utils.CustomControllerError(c, http.StatusUnauthorized, common.ErrTokenInvalidOrExpired.Error())
		return
	}
	userID := jwtClaims.Subject
	clientId := jwtClaims.Audience[0]
	claims, err := oc.oidcService.GetUserClaimsForClient(userID, clientId)
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(http.StatusOK, claims)
}

func (oc *OidcController) getClientHandler(c *gin.Context) {
	clientId := c.Param("id")
	client, err := oc.oidcService.GetClient(clientId)
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	// Return a different DTO based on the user's role
	if c.GetBool("userIsAdmin") {
		clientDto := dto.OidcClientDto{}
		err = dto.MapStruct(client, &clientDto)
		if err == nil {
			c.JSON(http.StatusOK, clientDto)
			return
		}
	} else {
		clientDto := dto.PublicOidcClientDto{}
		err = dto.MapStruct(client, &clientDto)
		if err == nil {
			c.JSON(http.StatusOK, clientDto)
			return
		}
	}

	utils.ControllerError(c, err)
}

func (oc *OidcController) listClientsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	searchTerm := c.Query("search")

	clients, pagination, err := oc.oidcService.ListClients(searchTerm, page, pageSize)
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	var clientsDto []dto.OidcClientDto
	if err := dto.MapStructList(clients, &clientsDto); err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       clientsDto,
		"pagination": pagination,
	})
}

func (oc *OidcController) createClientHandler(c *gin.Context) {
	var input dto.OidcClientCreateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ControllerError(c, err)
		return
	}

	client, err := oc.oidcService.CreateClient(input, c.GetString("userID"))
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	var clientDto dto.OidcClientDto
	if err := dto.MapStruct(client, &clientDto); err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, clientDto)
}

func (oc *OidcController) deleteClientHandler(c *gin.Context) {
	err := oc.oidcService.DeleteClient(c.Param("id"))
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (oc *OidcController) updateClientHandler(c *gin.Context) {
	var input dto.OidcClientCreateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ControllerError(c, err)
		return
	}

	client, err := oc.oidcService.UpdateClient(c.Param("id"), input)
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	var clientDto dto.OidcClientDto
	if err := dto.MapStruct(client, &clientDto); err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(http.StatusOK, clientDto)
}

func (oc *OidcController) createClientSecretHandler(c *gin.Context) {
	secret, err := oc.oidcService.CreateClientSecret(c.Param("id"))
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": secret})
}

func (oc *OidcController) getClientLogoHandler(c *gin.Context) {
	imagePath, mimeType, err := oc.oidcService.GetClientLogo(c.Param("id"))
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.Header("Content-Type", mimeType)
	c.File(imagePath)
}

func (oc *OidcController) updateClientLogoHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	err = oc.oidcService.UpdateClientLogo(c.Param("id"), file)
	if err != nil {
		if errors.Is(err, common.ErrFileTypeNotSupported) {
			utils.CustomControllerError(c, http.StatusBadRequest, err.Error())
		} else {
			utils.ControllerError(c, err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (oc *OidcController) deleteClientLogoHandler(c *gin.Context) {
	err := oc.oidcService.DeleteClientLogo(c.Param("id"))
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
