package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/middleware"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"net/http"
)

func NewApplicationConfigurationController(
	group *gin.RouterGroup,
	jwtAuthMiddleware *middleware.JwtAuthMiddleware,
	appConfigService *service.AppConfigService) {

	acc := &ApplicationConfigurationController{
		appConfigService: appConfigService,
	}
	group.GET("/application-configuration", acc.listApplicationConfigurationHandler)
	group.GET("/application-configuration/all", jwtAuthMiddleware.Add(true), acc.listAllApplicationConfigurationHandler)
	group.PUT("/application-configuration", acc.updateApplicationConfigurationHandler)

	group.GET("/application-configuration/logo", acc.getLogoHandler)
	group.GET("/application-configuration/background-image", acc.getBackgroundImageHandler)
	group.GET("/application-configuration/favicon", acc.getFaviconHandler)
	group.PUT("/application-configuration/logo", jwtAuthMiddleware.Add(true), acc.updateLogoHandler)
	group.PUT("/application-configuration/favicon", jwtAuthMiddleware.Add(true), acc.updateFaviconHandler)
	group.PUT("/application-configuration/background-image", jwtAuthMiddleware.Add(true), acc.updateBackgroundImageHandler)
}

type ApplicationConfigurationController struct {
	appConfigService *service.AppConfigService
}

func (acc *ApplicationConfigurationController) listApplicationConfigurationHandler(c *gin.Context) {
	configuration, err := acc.appConfigService.ListApplicationConfiguration(false)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(200, configuration)
}

func (acc *ApplicationConfigurationController) listAllApplicationConfigurationHandler(c *gin.Context) {
	configuration, err := acc.appConfigService.ListApplicationConfiguration(true)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(200, configuration)
}

func (acc *ApplicationConfigurationController) updateApplicationConfigurationHandler(c *gin.Context) {
	var input model.AppConfigUpdateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	savedConfigVariables, err := acc.appConfigService.UpdateApplicationConfiguration(input)
	if err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, savedConfigVariables)
}

func (acc *ApplicationConfigurationController) getLogoHandler(c *gin.Context) {
	imageType := acc.appConfigService.DbConfig.LogoImageType.Value
	acc.getImage(c, "logo", imageType)
}

func (acc *ApplicationConfigurationController) getFaviconHandler(c *gin.Context) {
	acc.getImage(c, "favicon", "ico")
}

func (acc *ApplicationConfigurationController) getBackgroundImageHandler(c *gin.Context) {
	imageType := acc.appConfigService.DbConfig.BackgroundImageType.Value
	acc.getImage(c, "background", imageType)
}

func (acc *ApplicationConfigurationController) updateLogoHandler(c *gin.Context) {
	imageType := acc.appConfigService.DbConfig.LogoImageType.Value
	acc.updateImage(c, "logo", imageType)
}

func (acc *ApplicationConfigurationController) updateFaviconHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	fileType := utils.GetFileExtension(file.Filename)
	if fileType != "ico" {
		utils.HandlerError(c, http.StatusBadRequest, "File must be of type .ico")
		return
	}
	acc.updateImage(c, "favicon", "ico")
}

func (acc *ApplicationConfigurationController) updateBackgroundImageHandler(c *gin.Context) {
	imageType := acc.appConfigService.DbConfig.BackgroundImageType.Value
	acc.updateImage(c, "background", imageType)
}

func (acc *ApplicationConfigurationController) getImage(c *gin.Context, name string, imageType string) {
	imagePath := fmt.Sprintf("%s/application-images/%s.%s", common.EnvConfig.UploadPath, name, imageType)
	mimeType := utils.GetImageMimeType(imageType)

	c.Header("Content-Type", mimeType)
	c.File(imagePath)
}

func (acc *ApplicationConfigurationController) updateImage(c *gin.Context, imageName string, oldImageType string) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.HandlerError(c, http.StatusBadRequest, common.ErrInvalidBody.Error())
		return
	}

	err = acc.appConfigService.UpdateImage(file, imageName, oldImageType)
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
