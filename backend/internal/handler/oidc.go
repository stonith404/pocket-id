package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-rest-api-template/internal/common"
	"golang-rest-api-template/internal/common/middleware"
	"golang-rest-api-template/internal/model"
	"golang-rest-api-template/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

func RegisterOIDCRoutes(group *gin.RouterGroup) {
	group.POST("/oidc/authorize", middleware.JWTAuth(false), authorizeHandler)
	group.POST("/oidc/authorize/new-client", middleware.JWTAuth(false), authorizeNewClientHandler)
	group.POST("/oidc/token", createIDTokenHandler)

	group.GET("/oidc/clients", middleware.JWTAuth(true), listClientsHandler)
	group.POST("/oidc/clients", middleware.JWTAuth(true), createClientHandler)
	group.GET("/oidc/clients/:id", getClientHandler)
	group.PUT("/oidc/clients/:id", middleware.JWTAuth(true), updateClientHandler)
	group.DELETE("/oidc/clients/:id", middleware.JWTAuth(true), deleteClientHandler)

	group.POST("/oidc/clients/:id/secret", middleware.JWTAuth(true), createClientSecretHandler)

	group.GET("/oidc/clients/:id/logo", getClientLogoHandler)
	group.DELETE("/oidc/clients/:id/logo", deleteClientLogoHandler)
	group.POST("/oidc/clients/:id/logo", middleware.JWTAuth(true), middleware.LimitFileSize(2<<20), updateClientLogoHandler)
}

type AuthorizeRequest struct {
	ClientID string `json:"clientID" binding:"required"`
	Scope    string `json:"scope" binding:"required"`
	Nonce    string `json:"nonce"`
}

func authorizeHandler(c *gin.Context) {
	var parsedBody AuthorizeRequest
	if err := c.ShouldBindJSON(&parsedBody); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	var userAuthorizedOIDCClient model.UserAuthorizedOidcClient
	common.DB.First(&userAuthorizedOIDCClient, "client_id = ? AND user_id = ?", parsedBody.ClientID, c.GetString("userID"))

	// If the record isn't found or the scope is different return an error
	// The client will have to call the authorizeNewClientHandler
	if userAuthorizedOIDCClient.Scope != parsedBody.Scope {
		utils.HandlerError(c, http.StatusForbidden, "missing authorization")
		return
	}

	authorizationCode, err := createAuthorizationCode(parsedBody.ClientID, c.GetString("userID"), parsedBody.Scope, parsedBody.Nonce)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": authorizationCode})
}

// authorizeNewClientHandler authorizes a new client for the user
// a new client is a new client when the user has not authorized the client before
func authorizeNewClientHandler(c *gin.Context) {
	var parsedBody model.AuthorizeNewClientDto
	if err := c.ShouldBindJSON(&parsedBody); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userAuthorizedClient := model.UserAuthorizedOidcClient{
		UserID:   c.GetString("userID"),
		ClientID: parsedBody.ClientID,
		Scope:    parsedBody.Scope,
	}
	err := common.DB.Create(&userAuthorizedClient).Error

	if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
		err = common.DB.Model(&userAuthorizedClient).Update("scope", parsedBody.Scope).Error
	}

	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	authorizationCode, err := createAuthorizationCode(parsedBody.ClientID, c.GetString("userID"), parsedBody.Scope, parsedBody.Nonce)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": authorizationCode})

}

func createIDTokenHandler(c *gin.Context) {
	var body model.OidcIdTokenDto

	if err := c.ShouldBind(&body); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Currently only authorization_code grant type is supported
	if body.GrantType != "authorization_code" {
		utils.HandlerError(c, http.StatusBadRequest, "grant type not supported")
		return
	}

	clientID := body.ClientID
	clientSecret := body.ClientSecret

	// Client id and secret can also be passed over the Authorization header
	if clientID == "" || clientSecret == "" {
		var ok bool
		clientID, clientSecret, ok = c.Request.BasicAuth()
		if !ok {
			utils.HandlerError(c, http.StatusBadRequest, "Client id and secret not provided")
			return
		}
	}

	// Get the client
	var client model.OidcClient
	err := common.DB.First(&client, "id = ?", clientID, clientSecret).Error
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "OIDC OIDC client not found")
		return
	}

	// Check if client secret is correct
	err = bcrypt.CompareHashAndPassword([]byte(client.Secret), []byte(clientSecret))
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid client secret")
		return
	}

	var authorizationCodeMetaData model.OidcAuthorizationCode
	err = common.DB.Preload("User").First(&authorizationCodeMetaData, "code = ?", body.Code).Error
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid authorization code")
		return
	}

	// Check if the client id matches the client id in the authorization code and if the code has expired
	if authorizationCodeMetaData.ClientID != clientID && authorizationCodeMetaData.ExpiresAt.Before(time.Now()) {
		utils.HandlerError(c, http.StatusBadRequest, "invalid authorization code")
		return
	}

	idToken, e := common.GenerateIDToken(authorizationCodeMetaData.User, clientID, authorizationCodeMetaData.Scope, authorizationCodeMetaData.Nonce)
	if e != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	// Delete the authorization code after it has been used
	common.DB.Delete(&authorizationCodeMetaData)

	c.JSON(http.StatusOK, gin.H{"id_token": idToken})
}

func getClientHandler(c *gin.Context) {
	clientId := c.Param("id")

	var client model.OidcClient
	err := common.DB.First(&client, "id = ?", clientId).Error
	if err != nil {
		utils.HandlerError(c, http.StatusNotFound, "OIDC client not found")
		return
	}

	c.JSON(http.StatusOK, client)
}

func listClientsHandler(c *gin.Context) {
	var clients []model.OidcClient
	searchTerm := c.Query("search")

	query := common.DB.Model(&model.OidcClient{})

	if searchTerm != "" {
		searchPattern := "%" + searchTerm + "%"
		query = query.Where("name LIKE ?", searchPattern)
	}

	pagination, err := utils.Paginate(c, query, &clients)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       clients,
		"pagination": pagination,
	})
}

func createClientHandler(c *gin.Context) {
	var input model.OidcClientCreateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	client := model.OidcClient{
		Name:        input.Name,
		CallbackURL: input.CallbackURL,
		CreatedByID: c.GetString("userID"),
	}

	if err := common.DB.Create(&client).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, client)
}

func deleteClientHandler(c *gin.Context) {
	var client model.OidcClient
	if err := common.DB.First(&client, "id = ?", c.Param("id")).Error; err != nil {
		utils.HandlerError(c, http.StatusNotFound, "OIDC OIDC client not found")
		return
	}

	if err := common.DB.Delete(&client).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func updateClientHandler(c *gin.Context) {
	var input model.OidcClientCreateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	var client model.OidcClient
	if err := common.DB.First(&client, "id = ?", c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandlerError(c, http.StatusNotFound, "OIDC client not found")
			return
		}
		utils.UnknownHandlerError(c, err)
		return
	}

	client.Name = input.Name
	client.CallbackURL = input.CallbackURL

	if err := common.DB.Save(&client).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, client)
}

// createClientSecretHandler creates a new secret for the client and revokes the old one
func createClientSecretHandler(c *gin.Context) {
	var client model.OidcClient
	if err := common.DB.First(&client, "id = ?", c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandlerError(c, http.StatusNotFound, "OIDC client not found")
			return
		}
		utils.UnknownHandlerError(c, err)
		return
	}

	clientSecret, err := utils.GenerateRandomAlphanumericString(32)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	client.Secret = string(hashedSecret)
	if err := common.DB.Save(&client).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": clientSecret})
}

func getClientLogoHandler(c *gin.Context) {
	var client model.OidcClient
	if err := common.DB.First(&client, "id = ?", c.Param("id")).Error; err != nil {
		utils.HandlerError(c, http.StatusNotFound, "OIDC client not found")
		return
	}

	if client.ImageType == nil {
		utils.HandlerError(c, http.StatusNotFound, "image not found")
		return
	}

	imageType := *client.ImageType

	imagePath := fmt.Sprintf("%s/oidc-client-images/%s.%s", common.EnvConfig.UploadPath, client.ID, imageType)
	mimeType := utils.GetImageMimeType(imageType)

	c.Header("Content-Type", mimeType)
	c.File(imagePath)
}

func updateClientLogoHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	fileType := utils.GetFileExtension(file.Filename)
	if mimeType := utils.GetImageMimeType(fileType); mimeType == "" {
		utils.HandlerError(c, http.StatusBadRequest, "file type not supported")
		return
	}

	imagePath := fmt.Sprintf("%s/oidc-client-images/%s.%s", common.EnvConfig.UploadPath, c.Param("id"), fileType)
	err = c.SaveUploadedFile(file, imagePath)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	var client model.OidcClient
	if err := common.DB.First(&client, "id = ?", c.Param("id")).Error; err != nil {
		utils.HandlerError(c, http.StatusNotFound, "OIDC client not found")
		return
	}

	// Delete the old image if it has a different file type
	if client.ImageType != nil && fileType != *client.ImageType {
		oldImagePath := fmt.Sprintf("%s/oidc-client-images/%s.%s", common.EnvConfig.UploadPath, client.ID, *client.ImageType)
		if err := os.Remove(oldImagePath); err != nil {
			utils.UnknownHandlerError(c, err)
			return
		}
	}

	client.ImageType = &fileType
	if err := common.DB.Save(&client).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteClientLogoHandler(c *gin.Context) {
	var client model.OidcClient
	if err := common.DB.First(&client, "id = ?", c.Param("id")).Error; err != nil {
		utils.HandlerError(c, http.StatusNotFound, "OIDC client not found")
		return
	}

	if client.ImageType == nil {
		utils.HandlerError(c, http.StatusNotFound, "image not found")
		return
	}

	imagePath := fmt.Sprintf("%s/oidc-client-images/%s.%s", common.EnvConfig.UploadPath, client.ID, *client.ImageType)
	if err := os.Remove(imagePath); err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	client.ImageType = nil
	if err := common.DB.Save(&client).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func createAuthorizationCode(clientID string, userID string, scope string, nonce string) (string, error) {
	randomString, err := utils.GenerateRandomAlphanumericString(32)
	if err != nil {
		return "", err
	}

	oidcAuthorizationCode := model.OidcAuthorizationCode{
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Code:      randomString,
		ClientID:  clientID,
		UserID:    userID,
		Scope:     scope,
		Nonce:     nonce,
	}

	if err := common.DB.Create(&oidcAuthorizationCode).Error; err != nil {
		return "", err
	}

	return randomString, nil
}
