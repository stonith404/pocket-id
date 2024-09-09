package controller

import (
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"github.com/stonith404/pocket-id/backend/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
)

func NewAuditLogController(group *gin.RouterGroup, auditLogService *service.AuditLogService, jwtAuthMiddleware *middleware.JwtAuthMiddleware) {
	alc := AuditLogController{
		auditLogService: auditLogService,
	}

	group.GET("/audit-logs", jwtAuthMiddleware.Add(false), alc.listAuditLogsForUserHandler)
}

type AuditLogController struct {
	auditLogService *service.AuditLogService
}

func (alc *AuditLogController) listAuditLogsForUserHandler(c *gin.Context) {
	userID := c.GetString("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Fetch audit logs for the user
	logs, pagination, err := alc.auditLogService.ListAuditLogsForUser(userID, page, pageSize)
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	// Map the audit logs to DTOs
	var logsDtos []dto.AuditLogDto
	err = dto.MapStructList(logs, &logsDtos)
	if err != nil {
		utils.ControllerError(c, err)
		return
	}

	// Add device information to the logs
	for i, logsDto := range logsDtos {
		logsDto.Device = alc.auditLogService.DeviceStringFromUserAgent(logs[i].UserAgent)
		logsDtos[i] = logsDto
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       logsDtos,
		"pagination": pagination,
	})
}
