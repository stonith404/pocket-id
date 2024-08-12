package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"golang-rest-api-template/internal/common"
	"golang-rest-api-template/internal/common/middleware"
	"golang-rest-api-template/internal/model"
	"golang-rest-api-template/internal/utils"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

func RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/webauthn/register/start", middleware.JWTAuth(false), beginRegistrationHandler)
	group.POST("/webauthn/register/finish", middleware.JWTAuth(false), verifyRegistrationHandler)

	group.GET("/webauthn/login/start", beginLoginHandler)
	group.POST("/webauthn/login/finish", middleware.RateLimiter(rate.Every(10*time.Second), 5), verifyLoginHandler)

	group.POST("/webauthn/logout", middleware.JWTAuth(false), logoutHandler)

	group.GET("/webauthn/credentials", middleware.JWTAuth(false), listCredentialsHandler)
	group.PATCH("/webauthn/credentials/:id", middleware.JWTAuth(false), updateCredentialHandler)
	group.DELETE("/webauthn/credentials/:id", middleware.JWTAuth(false), deleteCredentialHandler)
}

func beginRegistrationHandler(c *gin.Context) {
	var user model.User
	err := common.DB.Preload("Credentials").Find(&user, "id = ?", c.GetString("userID")).Error
	if err != nil {
		utils.UnknownHandlerError(c, err)
		log.Println(err)
		return
	}

	options, session, err := common.WebAuthn.BeginRegistration(&user, webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired), webauthn.WithExclusions(user.WebAuthnCredentialDescriptors()))
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	// Save the webauthn session so we can retrieve it in the verifyRegistrationHandler
	sessionToStore := &model.WebauthnSession{
		ExpiresAt:        session.Expires,
		Challenge:        session.Challenge,
		UserVerification: string(session.UserVerification),
	}

	if err = common.DB.Create(&sessionToStore).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.SetCookie("session_id", sessionToStore.ID, int(common.WebAuthn.Config.Timeouts.Registration.Timeout.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, options.Response)
}

func verifyRegistrationHandler(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "Session ID missing")
		return
	}

	// Retrieve the session that was previously created by the beginRegistrationHandler
	var storedSession model.WebauthnSession
	err = common.DB.First(&storedSession, "id = ?", sessionID).Error

	session := webauthn.SessionData{
		Challenge: storedSession.Challenge,
		Expires:   storedSession.ExpiresAt,
		UserID:    []byte(c.GetString("userID")),
	}

	var user model.User
	err = common.DB.Find(&user, "id = ?", c.GetString("userID")).Error
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	credential, err := common.WebAuthn.FinishRegistration(&user, session, c.Request)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	credentialToStore := model.WebauthnCredential{
		Name:            "New Passkey",
		CredentialID:    string(credential.ID),
		AttestationType: credential.AttestationType,
		PublicKey:       credential.PublicKey,
		Transport:       credential.Transport,
		UserID:          user.ID,
	}
	if err := common.DB.Create(&credentialToStore).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, credentialToStore)
}

func beginLoginHandler(c *gin.Context) {
	options, session, err := common.WebAuthn.BeginDiscoverableLogin()
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	// Save the webauthn session so we can retrieve it in the verifyLoginHandler
	sessionToStore := &model.WebauthnSession{
		ExpiresAt:        session.Expires,
		Challenge:        session.Challenge,
		UserVerification: string(session.UserVerification),
	}

	if err = common.DB.Create(&sessionToStore).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.SetCookie("session_id", sessionToStore.ID, int(common.WebAuthn.Config.Timeouts.Registration.Timeout.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, options.Response)
}

func verifyLoginHandler(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "Session ID missing")
		return
	}

	credentialAssertionData, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "Invalid body")
		return
	}

	// Retrieve the session that was previously created by the beginLoginHandler
	var storedSession model.WebauthnSession
	if err := common.DB.First(&storedSession, "id = ?", sessionID).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	session := webauthn.SessionData{
		Challenge: storedSession.Challenge,
		Expires:   storedSession.ExpiresAt,
	}

	var user *model.User
	_, err = common.WebAuthn.ValidateDiscoverableLogin(func(_, userHandle []byte) (webauthn.User, error) {
		if err := common.DB.Preload("Credentials").First(&user, "id = ?", string(userHandle)).Error; err != nil {
			return nil, err
		}
		return user, nil
	}, session, credentialAssertionData)

	if err != nil {
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			utils.HandlerError(c, http.StatusBadRequest, "no user with this passkey exists")
		} else {
			utils.UnknownHandlerError(c, err)
		}
		return
	}

	err = common.DB.Find(&user, "id = ?", c.GetString("userID")).Error
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	token, err := common.GenerateAccessToken(*user)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.SetCookie("access_token", token, int(time.Hour.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, user)
}

func listCredentialsHandler(c *gin.Context) {
	var credentials []model.WebauthnCredential
	if err := common.DB.Find(&credentials, "user_id = ?", c.GetString("userID")).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, credentials)
}

func deleteCredentialHandler(c *gin.Context) {
	var passkeyCount int64
	if err := common.DB.Model(&model.WebauthnCredential{}).Where("user_id = ?", c.GetString("userID")).Count(&passkeyCount).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	if passkeyCount == 1 {
		utils.HandlerError(c, http.StatusBadRequest, "You must have at least one passkey")
		return
	}

	var credential model.WebauthnCredential
	if err := common.DB.First(&credential, "id = ? AND user_id = ?", c.Param("id"), c.GetString("userID")).Error; err != nil {
		utils.HandlerError(c, http.StatusNotFound, "Credential not found")
		return
	}

	if err := common.DB.Delete(&credential).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func updateCredentialHandler(c *gin.Context) {
	var credential model.WebauthnCredential
	if err := common.DB.Where("id = ? AND user_id = ?", c.Param("id"), c.GetString("userID")).First(&credential).Error; err != nil {
		utils.HandlerError(c, http.StatusNotFound, "Credential not found")
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	credential.Name = input.Name

	if err := common.DB.Save(&credential).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func logoutHandler(c *gin.Context) {
	c.SetCookie("access_token", "", 0, "/", "", false, true)
	c.Status(http.StatusNoContent)
}
