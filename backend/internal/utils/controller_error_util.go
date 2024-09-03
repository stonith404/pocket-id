package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

import (
	"fmt"
)

func ControllerError(c *gin.Context, err error) {
	// Check for record not found errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		CustomControllerError(c, http.StatusNotFound, "Record not found")
		return
	}

	// Check for validation errors
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		message := handleValidationError(validationErrors)
		CustomControllerError(c, http.StatusBadRequest, message)
		return

	}

	log.Println(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
}

func handleValidationError(validationErrors validator.ValidationErrors) string {
	var errorMessages []string

	for _, ve := range validationErrors {
		fieldName := ve.Field()
		var errorMessage string
		switch ve.Tag() {
		case "required":
			errorMessage = fmt.Sprintf("%s is required", fieldName)
		case "email":
			errorMessage = fmt.Sprintf("%s must be a valid email address", fieldName)
		case "username":
			errorMessage = fmt.Sprintf("%s must only contain lowercase letters, numbers, underscores, dots, hyphens, and '@' symbols and not start or end with a special character", fieldName)
		case "url":
			errorMessage = fmt.Sprintf("%s must be a valid URL", fieldName)
		case "min":
			errorMessage = fmt.Sprintf("%s must be at least %s characters long", fieldName, ve.Param())
		case "max":
			errorMessage = fmt.Sprintf("%s must be at most %s characters long", fieldName, ve.Param())
		case "urlList":
			errorMessage = fmt.Sprintf("%s must be a list of valid URLs", fieldName)
		default:
			errorMessage = fmt.Sprintf("%s is invalid", fieldName)
		}

		errorMessages = append(errorMessages, errorMessage)
	}

	// Join all the error messages into a single string
	combinedErrors := strings.Join(errorMessages, ", ")

	return combinedErrors
}

func CustomControllerError(c *gin.Context, statusCode int, message string) {
	// Capitalize the first letter of the message
	message = strings.ToUpper(message[:1]) + message[1:]
	c.JSON(statusCode, gin.H{"error": message})
}
