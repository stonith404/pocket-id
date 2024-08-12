package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang-rest-api-template/internal/common"
	"golang-rest-api-template/internal/common/middleware"
	"golang-rest-api-template/internal/model"
	"golang-rest-api-template/internal/utils"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func RegisterUserRoutes(group *gin.RouterGroup) {
	group.GET("/users", middleware.JWTAuth(true), listUsersHandler)
	group.GET("/users/me", middleware.JWTAuth(false), getCurrentUserHandler)
	group.GET("/users/:id", middleware.JWTAuth(true), getUserHandler)
	group.POST("/users", middleware.JWTAuth(true), createUserHandler)
	group.PUT("/users/:id", middleware.JWTAuth(true), updateUserHandler)
	group.PUT("/users/me", middleware.JWTAuth(false), updateCurrentUserHandler)
	group.DELETE("/users/:id", middleware.JWTAuth(true), deleteUserHandler)

	group.POST("/users/:id/one-time-access-token", middleware.JWTAuth(true), createOneTimeAccessTokenHandler)
	group.POST("/one-time-access-token/:token", middleware.RateLimiter(rate.Every(10*time.Second), 5), exchangeOneTimeAccessTokenHandler)
	group.POST("/one-time-access-token/setup", getSetupAccessTokenHandler)
}

func listUsersHandler(c *gin.Context) {
	var users []model.User
	searchTerm := c.Query("search")

	query := common.DB.Model(&model.User{})

	if searchTerm != "" {
		searchPattern := "%" + searchTerm + "%"
		query = query.Where("email LIKE ? OR first_name LIKE ? OR username LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	pagination, err := utils.Paginate(c, query, &users)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       users,
		"pagination": pagination,
	})
}

func getUserHandler(c *gin.Context) {
	var user model.User
	if err := common.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandlerError(c, http.StatusNotFound, "User not found")
			return
		}
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func getCurrentUserHandler(c *gin.Context) {
	var user model.User
	if err := common.DB.Where("id = ?", c.GetString("userID")).First(&user).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)

}

func deleteUserHandler(c *gin.Context) {
	var user model.User
	if err := common.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandlerError(c, http.StatusNotFound, "User not found")
			return
		}
		utils.UnknownHandlerError(c, err)
		return
	}

	if err := common.DB.Delete(&user).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func createUserHandler(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := common.DB.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if err := checkDuplicatedFields(user); err != nil {
				utils.HandlerError(c, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			utils.UnknownHandlerError(c, err)
			return
		}
	}

	c.JSON(http.StatusCreated, user)
}

func updateUserHandler(c *gin.Context) {
	updateUser(c, c.Param("id"))
}

func updateCurrentUserHandler(c *gin.Context) {
	updateUser(c, c.GetString("userID"))
}

func createOneTimeAccessTokenHandler(c *gin.Context) {
	var input model.OneTimeAccessTokenCreateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	randomString, err := utils.GenerateRandomAlphanumericString(16)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	oneTimeAccessToken := model.OneTimeAccessToken{
		UserID:    input.UserID,
		ExpiresAt: input.ExpiresAt,
		Token:     randomString,
	}

	if err := common.DB.Create(&oneTimeAccessToken).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": oneTimeAccessToken.Token})
}

func exchangeOneTimeAccessTokenHandler(c *gin.Context) {
	var oneTimeAccessToken model.OneTimeAccessToken
	if err := common.DB.Where("token = ? AND expires_at > ?", c.Param("token"), utils.FormatDateForDb(time.Now())).Preload("User").First(&oneTimeAccessToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandlerError(c, http.StatusForbidden, "Token is invalid or expired")
			return
		}
		utils.UnknownHandlerError(c, err)
		return
	}

	token, err := common.GenerateAccessToken(oneTimeAccessToken.User)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		log.Println(err)
		return
	}

	if err := common.DB.Delete(&oneTimeAccessToken).Error; err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.SetCookie("access_token", token, int(time.Hour.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, oneTimeAccessToken.User)
}

// getSetupAccessTokenHandler creates the initial admin user and returns an access token for the user
// This handler is only available if there are no users in the database
func getSetupAccessTokenHandler(c *gin.Context) {
	var userCount int64
	if err := common.DB.Model(&model.User{}).Count(&userCount).Error; err != nil {
		log.Fatal("failed to count users", err)
	}

	// If there are more than one user, we don't need to create the admin user
	if userCount > 1 {
		utils.HandlerError(c, http.StatusForbidden, "Setup already completed")
		return
	}

	var user = model.User{
		FirstName: "Admin",
		LastName:  "Admin",
		Username:  "admin",
		Email:     "admin@admin.com",
		IsAdmin:   true,
	}

	// Create the initial admin user if it doesn't exist
	if err := common.DB.Model(&model.User{}).Preload("Credentials").FirstOrCreate(&user).Error; err != nil {
		log.Fatal("failed to create admin user", err)
	}

	// If the user already has credentials, the setup is already completed
	if len(user.Credentials) > 0 {
		utils.HandlerError(c, http.StatusForbidden, "Setup already completed")
		return
	}

	token, err := common.GenerateAccessToken(user)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		log.Println(err)
		return
	}
	c.SetCookie("access_token", token, int(time.Hour.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context, userID string) {
	var user model.User
	if err := common.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandlerError(c, http.StatusNotFound, "User not found")
			return
		}
		utils.UnknownHandlerError(c, err)
		return
	}

	var updatedUser model.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := common.DB.Model(&user).Updates(&updatedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if err := checkDuplicatedFields(user); err != nil {
				utils.HandlerError(c, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			utils.UnknownHandlerError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, updatedUser)
}

func checkDuplicatedFields(user model.User) error {
	var existingUser model.User

	if common.DB.Where("id != ? AND email = ?", user.ID, user.Email).First(&existingUser).Error == nil {
		return errors.New("email is already taken")
	}

	if common.DB.Where("id != ? AND username = ?", user.ID, user.Username).First(&existingUser).Error == nil {
		return errors.New("username is already taken")
	}

	return nil
}
