package middleware

import (
	"github.com/stonith404/pocket-id/backend/internal/common"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct{}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

func (m *CorsMiddleware) Add() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{common.EnvConfig.AppURL},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
		MaxAge:       12 * time.Hour,
	})
}
