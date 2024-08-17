package model

import (
	"gorm.io/gorm"
	"time"
)

type UserAuthorizedOidcClient struct {
	Scope  string
	UserID string `json:"userId" gorm:"primary_key;"`

	ClientID string `json:"clientId" gorm:"primary_key;"`
	Client   OidcClient
}

type OidcClient struct {
	Base

	Name        string  `json:"name"`
	Secret      string  `json:"-"`
	CallbackURL string  `json:"callbackURL"`
	ImageType   *string `json:"-"`
	HasLogo     bool    `gorm:"-" json:"hasLogo"`

	CreatedByID string
	CreatedBy   User
}

func (c *OidcClient) AfterFind(_ *gorm.DB) (err error) {
	// Compute HasLogo field
	c.HasLogo = c.ImageType != nil && *c.ImageType != ""
	return nil
}

type OidcAuthorizationCode struct {
	Base

	Code      string
	Scope     string
	Nonce     string
	ExpiresAt time.Time

	UserID string
	User   User

	ClientID string
}

type OidcClientCreateDto struct {
	Name        string `json:"name" binding:"required"`
	CallbackURL string `json:"callbackURL" binding:"required"`
}

type AuthorizeNewClientDto struct {
	ClientID string `json:"clientID" binding:"required"`
	Scope    string `json:"scope" binding:"required"`
	Nonce    string `json:"nonce"`
}

type OidcIdTokenDto struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code" binding:"required"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
}

type AuthorizeRequest struct {
	ClientID string `json:"clientID" binding:"required"`
	Scope    string `json:"scope" binding:"required"`
	Nonce    string `json:"nonce"`
}
