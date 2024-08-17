package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/middleware"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"net/http"
	"strconv"
)

func NewOidcController(group *gin.RouterGroup, jwtAuthMiddleware *middleware.JwtAuthMiddleware, fileSizeLimitMiddleware *middleware.FileSizeLimitMiddleware, oidcService *service.OidcService) {
	oc := &OidcController{OidcService: oidcService}

	group.POST("/oidc/authorize", jwtAuthMiddleware.Add(false), oc.authorizeHandler)
	group.POST("/oidc/authorize/new-client", jwtAuthMiddleware.Add(false), oc.authorizeNewClientHandler)
	group.POST("/oidc/token", oc.createIDTokenHandler)

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
	OidcService *service.OidcService
}

func (oc *OidcController) authorizeHandler(c *gin.Context) {
	var parsedBody model.AuthorizeRequest
	if err := c.ShouldBindJSON(&parsedBody); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	code, err := oc.OidcService.Authorize(parsedBody, c.GetString("userID"))
	if err != nil {
		if errors.Is(err, common.ErrOidcMissingAuthorization) {
			utils.HandlerError(c, http.StatusForbidden, err.Error())
		} else {
			utils.UnknownHandlerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
}

func (oc *OidcController) authorizeNewClientHandler(c *gin.Context) {
	var parsedBody model.AuthorizeNewClientDto
	if err := c.ShouldBindJSON(&parsedBody); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	code, err := oc.OidcService.AuthorizeNewClient(parsedBody, c.GetString("userID"))
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
}

func (oc *OidcController) createIDTokenHandler(c *gin.Context) {
	var body model.OidcIdTokenDto

	if err := c.ShouldBind(&body); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	idToken, err := oc.OidcService.CreateIDToken(body)
	if err != nil {
		if errors.Is(err, common.ErrOidcGrantTypeNotSupported) ||
			errors.Is(err, common.ErrOidcMissingClientCredentials) ||
			errors.Is(err, common.ErrOidcClientSecretInvalid) ||
			errors.Is(err, common.ErrOidcInvalidAuthorizationCode) {
			utils.HandlerError(c, http.StatusBadRequest, err.Error())
		} else {
			utils.UnknownHandlerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"id_token": idToken})
}

func (oc *OidcController) getClientHandler(c *gin.Context) {
	clientId := c.Param("id")
	client, err := oc.OidcService.GetClient(clientId)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, client)
}

func (oc *OidcController) listClientsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	searchTerm := c.Query("search")

	clients, pagination, err := oc.OidcService.ListClients(searchTerm, page, pageSize)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       clients,
		"pagination": pagination,
	})
}

func (oc *OidcController) createClientHandler(c *gin.Context) {
	var input model.OidcClientCreateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	client, err := oc.OidcService.CreateClient(input, c.GetString("userID"))
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, client)
}

func (oc *OidcController) deleteClientHandler(c *gin.Context) {
	err := oc.OidcService.DeleteClient(c.Param("id"))
	if err != nil {
		utils.HandlerError(c, http.StatusNotFound, "OIDC client not found")
		return
	}

	c.Status(http.StatusNoContent)
}

func (oc *OidcController) updateClientHandler(c *gin.Context) {
	var input model.OidcClientCreateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	client, err := oc.OidcService.UpdateClient(c.Param("id"), input)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, client)
}

func (oc *OidcController) createClientSecretHandler(c *gin.Context) {
	secret, err := oc.OidcService.CreateClientSecret(c.Param("id"))
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": secret})
}

func (oc *OidcController) getClientLogoHandler(c *gin.Context) {
	imagePath, mimeType, err := oc.OidcService.GetClientLogo(c.Param("id"))
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Header("Content-Type", mimeType)
	c.File(imagePath)
}

func (oc *OidcController) updateClientLogoHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	err = oc.OidcService.UpdateClientLogo(c.Param("id"), file)
	if err != nil {
		if errors.Is(err, common.ErrFileTypeNotSupported) {
			utils.HandlerError(c, http.StatusBadRequest, err.Error())
		} else {
			utils.UnknownHandlerError(c, err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (oc *OidcController) deleteClientLogoHandler(c *gin.Context) {
	err := oc.OidcService.DeleteClientLogo(c.Param("id"))
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
