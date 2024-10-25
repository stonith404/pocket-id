package dto

import "time"

type UserDto struct {
	ID           string           `json:"id"`
	Username     string           `json:"username"`
	Email        string           `json:"email" `
	FirstName    string           `json:"firstName"`
	LastName     string           `json:"lastName"`
	IsAdmin      bool             `json:"isAdmin"`
	CustomClaims []CustomClaimDto `json:"customClaims"`
}

type UserCreateDto struct {
	Username  string `json:"username" binding:"required,username,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName" binding:"required,min=3,max=30"`
	LastName  string `json:"lastName" binding:"required,min=3,max=30"`
	IsAdmin   bool   `json:"isAdmin" binding:"required"`
}

type OneTimeAccessTokenCreateDto struct {
	UserID    string    `json:"userId" binding:"required"`
	ExpiresAt time.Time `json:"expiresAt" binding:"required"`
}
