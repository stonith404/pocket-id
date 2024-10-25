package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"github.com/stonith404/pocket-id/backend/internal/middleware"
	"github.com/stonith404/pocket-id/backend/internal/service"
	"net/http"
)

func NewCustomClaimController(group *gin.RouterGroup, jwtAuthMiddleware *middleware.JwtAuthMiddleware, customClaimService *service.CustomClaimService) {
	wkc := &CustomClaimController{customClaimService: customClaimService}
	group.GET("/custom-claims/suggestions", jwtAuthMiddleware.Add(true), wkc.getSuggestionsHandler)
	group.PUT("/custom-claims/user/:userId", jwtAuthMiddleware.Add(true), wkc.updateUserCustomClaimsHandler)
}

type CustomClaimController struct {
	customClaimService *service.CustomClaimService
}

func (ccc *CustomClaimController) getSuggestionsHandler(c *gin.Context) {
	claims, err := ccc.customClaimService.GetSuggestions()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, claims)
}

func (ccc *CustomClaimController) updateUserCustomClaimsHandler(c *gin.Context) {
	var input []dto.CustomClaimCreateDto

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	userId := c.Param("userId")
	claims, err := ccc.customClaimService.UpdateUserCustomClaims(userId, input)
	if err != nil {
		c.Error(err)
		return
	}

	var customClaimsDto []dto.CustomClaimDto
	if err := dto.MapStructList(claims, &customClaimsDto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, customClaimsDto)
}
