package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func UnknownHandlerError(c *gin.Context, err error) {
	log.Println(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
}

func HandlerError(c *gin.Context, statusCode int, message string) {
	// Capitalize the first letter of the message
	message = strings.ToUpper(message[:1]) + message[1:]
	c.JSON(statusCode, gin.H{"error": message})
}
