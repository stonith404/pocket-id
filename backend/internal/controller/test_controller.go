package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"github.com/stonith404/pocket-id/backend/internal/utils"
)

func NewTestController(group *gin.RouterGroup, testService *service.TestService) {
	testController := &TestController{TestService: testService}

	group.POST("/test/reset", testController.resetAndSeedHandler)
}

type TestController struct {
	TestService *service.TestService
}

func (tc *TestController) resetAndSeedHandler(c *gin.Context) {
	if err := tc.TestService.ResetDatabase(); err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	if err := tc.TestService.ResetApplicationImages(); err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	if err := tc.TestService.SeedDatabase(); err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Database reset and seeded"})
}
