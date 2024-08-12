package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-rest-api-template/internal/common"
	"golang-rest-api-template/internal/common/middleware"
	"golang-rest-api-template/internal/model"
	"golang-rest-api-template/internal/utils"
	"gorm.io/gorm"
	"net/http"
	"os"
	"reflect"
)

func RegisterConfigurationRoutes(group *gin.RouterGroup) {
	group.GET("/application-configuration", listApplicationConfigurationHandler)
	group.PUT("/application-configuration", updateApplicationConfigurationHandler)

	group.GET("/application-configuration/logo", getLogoHandler)
	group.GET("/application-configuration/background-image", getBackgroundImageHandler)
	group.GET("/application-configuration/favicon", getFaviconHandler)
	group.PUT("/application-configuration/logo", middleware.JWTAuth(true), updateLogoHandler)
	group.PUT("/application-configuration/favicon", middleware.JWTAuth(true), updateFaviconHandler)
	group.PUT("/application-configuration/background-image", middleware.JWTAuth(true), updateBackgroundImageHandler)
}

func listApplicationConfigurationHandler(c *gin.Context) {
	// Return also the private configuration variables if the user is admin and showAll is true
	showAll := c.GetBool("userIsAdmin") && c.DefaultQuery("showAll", "false") == "true"

	var configuration []model.ApplicationConfigurationVariable
	var err error

	if showAll {
		err = common.DB.Find(&configuration).Error
	} else {
		err = common.DB.Find(&configuration, "is_public = true").Error
	}

	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(200, configuration)
}

func updateApplicationConfigurationHandler(c *gin.Context) {
	var input model.ApplicationConfigurationUpdateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	savedConfigVariables := make([]model.ApplicationConfigurationVariable, 10)

	tx := common.DB.Begin()
	rt := reflect.ValueOf(input).Type()
	rv := reflect.ValueOf(input)

	// Loop over the input struct fields and update the related configuration variables
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		key := field.Tag.Get("json")
		value := rv.FieldByName(field.Name).String()

		// Get the existing configuration variable from the db
		var applicationConfigurationVariable model.ApplicationConfigurationVariable
		if err := tx.First(&applicationConfigurationVariable, "key = ? AND is_internal = false", key).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.HandlerError(c, http.StatusNotFound, fmt.Sprintf("Invalid configuration variable '%s'", value))
			} else {
				utils.UnknownHandlerError(c, err)
			}
			return
		}

		// Update the value of the existing configuration variable and save it
		applicationConfigurationVariable.Value = value
		if err := tx.Save(&applicationConfigurationVariable).Error; err != nil {
			tx.Rollback()
			utils.UnknownHandlerError(c, err)
			return
		}

		savedConfigVariables[i] = applicationConfigurationVariable
	}

	tx.Commit()

	if err := common.LoadDbConfigFromDb(); err != nil {
		utils.UnknownHandlerError(c, err)
	}

	c.JSON(http.StatusOK, savedConfigVariables)

}

func getLogoHandler(c *gin.Context) {
	imagType := common.DbConfig.LogoImageType.Value
	getImage(c, "logo", imagType)
}

func getFaviconHandler(c *gin.Context) {
	getImage(c, "favicon", "ico")
}

func getBackgroundImageHandler(c *gin.Context) {
	imageType := common.DbConfig.BackgroundImageType.Value
	getImage(c, "background", imageType)
}

func updateLogoHandler(c *gin.Context) {
	imageType := common.DbConfig.LogoImageType.Value
	updateImage(c, "logo", imageType)
}

func updateFaviconHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	fileType := utils.GetFileExtension(file.Filename)
	if fileType != "ico" {
		utils.HandlerError(c, http.StatusBadRequest, "File must be of type .ico")
		return
	}
	updateImage(c, "favicon", "ico")
}

func updateBackgroundImageHandler(c *gin.Context) {
	imagType := common.DbConfig.BackgroundImageType.Value
	updateImage(c, "background", imagType)
}

func getImage(c *gin.Context, name string, imageType string) {
	imagePath := fmt.Sprintf("%s/application-images/%s.%s", common.EnvConfig.UploadPath, name, imageType)
	mimeType := utils.GetImageMimeType(imageType)

	c.Header("Content-Type", mimeType)
	c.File(imagePath)
}

func updateImage(c *gin.Context, imageName string, oldImageType string) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	fileType := utils.GetFileExtension(file.Filename)
	if mimeType := utils.GetImageMimeType(fileType); mimeType == "" {
		utils.HandlerError(c, http.StatusBadRequest, "File type not supported")
		return
	}

	// Delete the old image if it has a different file type
	if fileType != oldImageType {
		oldImagePath := fmt.Sprintf("%s/application-images/%s.%s", common.EnvConfig.UploadPath, imageName, oldImageType)
		if err := os.Remove(oldImagePath); err != nil {
			utils.UnknownHandlerError(c, err)
			return
		}
	}

	imagePath := fmt.Sprintf("%s/application-images/%s.%s", common.EnvConfig.UploadPath, imageName, fileType)
	err = c.SaveUploadedFile(file, imagePath)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	// Update the file type in the database
	key := fmt.Sprintf("%sImageType", imageName)
	err = common.DB.Model(&model.ApplicationConfigurationVariable{}).Where("key = ?", key).Update("value", fileType).Error
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	if err := common.LoadDbConfigFromDb(); err != nil {
		utils.UnknownHandlerError(c, err)
	}

	c.Status(http.StatusNoContent)
}
