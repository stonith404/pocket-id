package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/middleware"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"time"
)

func NewUserController(group *gin.RouterGroup, jwtAuthMiddleware *middleware.JwtAuthMiddleware, rateLimitMiddleware *middleware.RateLimitMiddleware, userService *service.UserService) {
	uc := UserController{
		UserService: userService,
	}

	group.GET("/users", jwtAuthMiddleware.Add(true), uc.listUsersHandler)
	group.GET("/users/me", jwtAuthMiddleware.Add(false), uc.getCurrentUserHandler)
	group.GET("/users/:id", jwtAuthMiddleware.Add(true), uc.getUserHandler)
	group.POST("/users", jwtAuthMiddleware.Add(true), uc.createUserHandler)
	group.PUT("/users/:id", jwtAuthMiddleware.Add(true), uc.updateUserHandler)
	group.PUT("/users/me", jwtAuthMiddleware.Add(false), uc.updateCurrentUserHandler)
	group.DELETE("/users/:id", jwtAuthMiddleware.Add(true), uc.deleteUserHandler)

	group.POST("/users/:id/one-time-access-token", jwtAuthMiddleware.Add(true), uc.createOneTimeAccessTokenHandler)
	group.POST("/one-time-access-token/:token", rateLimitMiddleware.Add(rate.Every(10*time.Second), 5), uc.exchangeOneTimeAccessTokenHandler)
	group.POST("/one-time-access-token/setup", uc.getSetupAccessTokenHandler)
}

type UserController struct {
	UserService *service.UserService
}

func (uc *UserController) listUsersHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	searchTerm := c.Query("search")

	users, pagination, err := uc.UserService.ListUsers(searchTerm, page, pageSize)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       users,
		"pagination": pagination,
	})
}

func (uc *UserController) getUserHandler(c *gin.Context) {
	user, err := uc.UserService.GetUser(c.Param("id"))
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) getCurrentUserHandler(c *gin.Context) {
	user, err := uc.UserService.GetUser(c.GetString("userID"))
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) deleteUserHandler(c *gin.Context) {
	if err := uc.UserService.DeleteUser(c.Param("id")); err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (uc *UserController) createUserHandler(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	if err := uc.UserService.CreateUser(&user); err != nil {
		if errors.Is(err, common.ErrEmailTaken) || errors.Is(err, common.ErrUsernameTaken) {
			utils.HandlerError(c, http.StatusConflict, err.Error())
		} else {
			utils.UnknownHandlerError(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (uc *UserController) updateUserHandler(c *gin.Context) {
	uc.updateUser(c, false)
}

func (uc *UserController) updateCurrentUserHandler(c *gin.Context) {
	uc.updateUser(c, true)
}

func (uc *UserController) createOneTimeAccessTokenHandler(c *gin.Context) {
	var input model.OneTimeAccessTokenCreateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	token, err := uc.UserService.CreateOneTimeAccessToken(input.UserID, input.ExpiresAt)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (uc *UserController) exchangeOneTimeAccessTokenHandler(c *gin.Context) {
	user, token, err := uc.UserService.ExchangeOneTimeAccessToken(c.Param("token"))
	if err != nil {
		if errors.Is(err, common.ErrTokenInvalidOrExpired) {
			utils.HandlerError(c, http.StatusUnauthorized, err.Error())
		} else {
			utils.UnknownHandlerError(c, err)
		}
		return
	}

	c.SetCookie("access_token", token, int(time.Hour.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) getSetupAccessTokenHandler(c *gin.Context) {
	user, token, err := uc.UserService.SetupInitialAdmin()
	if err != nil {
		if errors.Is(err, common.ErrSetupAlreadyCompleted) {
			utils.HandlerError(c, http.StatusBadRequest, err.Error())
		} else {
			utils.UnknownHandlerError(c, err)
		}
		return
	}

	c.SetCookie("access_token", token, int(time.Hour.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) updateUser(c *gin.Context, updateOwnUser bool) {
	var updatedUser model.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	var userID string
	if updateOwnUser {
		userID = c.GetString("userID")
	} else {
		userID = c.Param("id")
	}

	user, err := uc.UserService.UpdateUser(userID, updatedUser, updateOwnUser)
	if err != nil {
		if errors.Is(err, common.ErrEmailTaken) || errors.Is(err, common.ErrUsernameTaken) {
			utils.HandlerError(c, http.StatusConflict, err.Error())
		} else {
			utils.UnknownHandlerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
