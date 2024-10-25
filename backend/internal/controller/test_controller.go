package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"net/http"
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
		c.Error(err)
		return
	}

	if err := tc.TestService.ResetApplicationImages(); err != nil {
		c.Error(err)
		return
	}

	if err := tc.TestService.SeedDatabase(); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
