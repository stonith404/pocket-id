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
	regex := "^[a-z0-9_]*$"
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
}
