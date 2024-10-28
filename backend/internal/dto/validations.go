package dto

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"net/url"
	"regexp"
)

var validateUrlList validator.Func = func(fl validator.FieldLevel) bool {
	urls := fl.Field().Interface().([]string)
	for _, u := range urls {
		_, err := url.ParseRequestURI(u)
		if err != nil {
			return false
		}
	}
	return true
}

var validateUsername validator.Func = func(fl validator.FieldLevel) bool {
	// [a-zA-Z0-9]      : The username must start with an alphanumeric character
	// [a-zA-Z0-9_.@-]* : The rest of the username can contain alphanumeric characters, dots, underscores, hyphens, and "@" symbols
	// [a-zA-Z0-9]$     : The username must end with an alphanumeric character
	regex := "^[a-zA-Z0-9][a-zA-Z0-9_.@-]*[a-zA-Z0-9]$"
	matched, _ := regexp.MatchString(regex, fl.Field().String())
	return matched
}

var validateUserGroupName validator.Func = func(fl validator.FieldLevel) bool {
	// The string can only contain lowercase letters, numbers, and underscores
	regex := "^[a-z0-9_]*$"
	matched, _ := regexp.MatchString(regex, fl.Field().String())
	return matched
}

var validateClaimKey validator.Func = func(fl validator.FieldLevel) bool {
	// The string can only contain letters and numbers
	regex := "^[A-Za-z0-9]*$"
	matched, _ := regexp.MatchString(regex, fl.Field().String())
	return matched
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("urlList", validateUrlList); err != nil {
			log.Fatalf("Failed to register custom validation: %v", err)
		}
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("username", validateUsername); err != nil {
			log.Fatalf("Failed to register custom validation: %v", err)
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("userGroupName", validateUserGroupName); err != nil {
			log.Fatalf("Failed to register custom validation: %v", err)
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("claimKey", validateClaimKey); err != nil {
			log.Fatalf("Failed to register custom validation: %v", err)
		}
	}
}
