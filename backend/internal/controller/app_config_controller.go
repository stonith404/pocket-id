package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"github.com/stonith404/pocket-id/backend/internal/middleware"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"net/http"
)

func NewAppConfigController(
	group *gin.RouterGroup,
	jwtAuthMiddleware *middleware.JwtAuthMiddleware,
	appConfigService *service.AppConfigService) {

	acc := &AppConfigController{
		appConfigService: appConfigService,
	}
	group.GET("/application-configuration", acc.listAppConfigHandler)
	group.GET("/application-configuration/all", jwtAuthMiddleware.Add(true), acc.listAllAppConfigHandler)
	group.PUT("/application-configuration", acc.updateAppConfigHandler)

	group.GET("/application-configuration/logo", acc.getLogoHandler)
	group.GET("/application-configuration/background-image", acc.getBackgroundImageHandler)
	group.GET("/application-configuration/favicon", acc.getFaviconHandler)
	group.PUT("/application-configuration/logo", jwtAuthMiddleware.Add(true), acc.updateLogoHandler)
	group.PUT("/application-configuration/favicon", jwtAuthMiddleware.Add(true), acc.updateFaviconHandler)
	group.PUT("/application-configuration/background-image", jwtAuthMiddleware.Add(true), acc.updateBackgroundImageHandler)
}

type AppConfigController struct {
	appConfigService *service.AppConfigService
}

func (acc *AppConfigController) listAppConfigHandler(c *gin.Context) {
	configuration, err := acc.appConfigService.ListAppConfig(false)
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	var configVariablesDto []dto.PublicAppConfigVariableDto
	if err := dto.MapStructList(configuration, &configVariablesDto); err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(200, configVariablesDto)
}

func (acc *AppConfigController) listAllAppConfigHandler(c *gin.Context) {
	configuration, err := acc.appConfigService.ListAppConfig(true)
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	var configVariablesDto []dto.AppConfigVariableDto
	if err := dto.MapStructList(configuration, &configVariablesDto); err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(200, configVariablesDto)
}

func (acc *AppConfigController) updateAppConfigHandler(c *gin.Context) {
	var input dto.AppConfigUpdateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ControllerError(c, err)
		return
	}

	savedConfigVariables, err := acc.appConfigService.UpdateAppConfig(input)
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	var configVariablesDto []dto.AppConfigVariableDto
	if err := dto.MapStructList(savedConfigVariables, &configVariablesDto); err != nil {
		utils.ControllerError(c, err)
		return
	}

	c.JSON(http.StatusOK, configVariablesDto)
}

func (acc *AppConfigController) getLogoHandler(c *gin.Context) {
	imageType := acc.appConfigService.DbConfig.LogoImageType.Value
	acc.getImage(c, "logo", imageType)
}

func (acc *AppConfigController) getFaviconHandler(c *gin.Context) {
	acc.getImage(c, "favicon", "ico")
}

func (acc *AppConfigController) getBackgroundImageHandler(c *gin.Context) {
	imageType := acc.appConfigService.DbConfig.BackgroundImageType.Value
	acc.getImage(c, "background", imageType)
}

func (acc *AppConfigController) updateLogoHandler(c *gin.Context) {
	imageType := acc.appConfigService.DbConfig.LogoImageType.Value
	acc.updateImage(c, "logo", imageType)
}

func (acc *AppConfigController) updateFaviconHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	fileType := utils.GetFileExtension(file.Filename)
	if fileType != "ico" {
		utils.CustomControllerError(c, http.StatusBadRequest, "File must be of type .ico")
		return
	}
	acc.updateImage(c, "favicon", "ico")
}

func (acc *AppConfigController) updateBackgroundImageHandler(c *gin.Context) {
	imageType := acc.appConfigService.DbConfig.BackgroundImageType.Value
	acc.updateImage(c, "background", imageType)
}

func (acc *AppConfigController) getImage(c *gin.Context, name string, imageType string) {
	imagePath := fmt.Sprintf("%s/application-images/%s.%s", common.EnvConfig.UploadPath, name, imageType)
	mimeType := utils.GetImageMimeType(imageType)

	c.Header("Content-Type", mimeType)
	c.File(imagePath)
}

func (acc *AppConfigController) updateImage(c *gin.Context, imageName string, oldImageType string) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	err = acc.appConfigService.UpdateImage(file, imageName, oldImageType)
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
