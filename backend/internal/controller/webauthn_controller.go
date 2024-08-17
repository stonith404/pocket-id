package controller

import (
	"errors"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/stonith404/pocket-id/backend/internal/middleware"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"golang.org/x/time/rate"
)

func NewWebauthnController(group *gin.RouterGroup, jwtAuthMiddleware *middleware.JwtAuthMiddleware, rateLimitMiddleware *middleware.RateLimitMiddleware, webauthnService *service.WebAuthnService, jwtService *service.JwtService) {
	wc := &WebauthnController{webAuthnService: webauthnService, jwtService: jwtService}
	group.GET("/webauthn/register/start", jwtAuthMiddleware.Add(false), wc.beginRegistrationHandler)
	group.POST("/webauthn/register/finish", jwtAuthMiddleware.Add(false), wc.verifyRegistrationHandler)

	group.GET("/webauthn/login/start", wc.beginLoginHandler)
	group.POST("/webauthn/login/finish", rateLimitMiddleware.Add(rate.Every(10*time.Second), 5), wc.verifyLoginHandler)

	group.POST("/webauthn/logout", jwtAuthMiddleware.Add(false), wc.logoutHandler)

	group.GET("/webauthn/credentials", jwtAuthMiddleware.Add(false), wc.listCredentialsHandler)
	group.PATCH("/webauthn/credentials/:id", jwtAuthMiddleware.Add(false), wc.updateCredentialHandler)
	group.DELETE("/webauthn/credentials/:id", jwtAuthMiddleware.Add(false), wc.deleteCredentialHandler)
}

type WebauthnController struct {
	webAuthnService *service.WebAuthnService
	jwtService      *service.JwtService
}

func (wc *WebauthnController) beginRegistrationHandler(c *gin.Context) {
	userID := c.GetString("userID")
	options, err := wc.webAuthnService.BeginRegistration(userID)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		log.Println(err)
		return
	}

	c.SetCookie("session_id", options.SessionID, int(options.Timeout.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, options.Response)
}

func (wc *WebauthnController) verifyRegistrationHandler(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "Session ID missing")
		return
	}

	userID := c.GetString("userID")
	credential, err := wc.webAuthnService.VerifyRegistration(sessionID, userID, c.Request)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, credential)
}

func (wc *WebauthnController) beginLoginHandler(c *gin.Context) {
	options, err := wc.webAuthnService.BeginLogin()
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.SetCookie("session_id", options.SessionID, int(options.Timeout.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, options.Response)
}

func (wc *WebauthnController) verifyLoginHandler(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "Session ID missing")
		return
	}

	credentialAssertionData, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	userID := c.GetString("userID")
	user, err := wc.webAuthnService.VerifyLogin(sessionID, userID, credentialAssertionData)
	if err != nil {
		if errors.Is(err, common.ErrInvalidCredentials) {
			utils.HandlerError(c, http.StatusUnauthorized, err.Error())
		} else {
			utils.UnknownHandlerError(c, err)
		}
		return
	}

	token, err := wc.jwtService.GenerateAccessToken(*user)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.SetCookie("access_token", token, int(time.Hour.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, user)
}

func (wc *WebauthnController) listCredentialsHandler(c *gin.Context) {
	userID := c.GetString("userID")
	credentials, err := wc.webAuthnService.ListCredentials(userID)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, credentials)
}

func (wc *WebauthnController) deleteCredentialHandler(c *gin.Context) {
	userID := c.GetString("userID")
	credentialID := c.Param("id")

	err := wc.webAuthnService.DeleteCredential(userID, credentialID)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (wc *WebauthnController) updateCredentialHandler(c *gin.Context) {
	userID := c.GetString("userID")
	credentialID := c.Param("id")

	var input model.WebauthnCredentialUpdateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	err := wc.webAuthnService.UpdateCredential(userID, credentialID, input.Name)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (wc *WebauthnController) logoutHandler(c *gin.Context) {
	c.SetCookie("access_token", "", 0, "/", "", false, true)
	c.Status(http.StatusNoContent)
}
