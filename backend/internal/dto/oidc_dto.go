package dto

type PublicOidcClientDto struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	HasLogo bool   `json:"hasLogo"`
}

type OidcClientDto struct {
	PublicOidcClientDto
	CallbackURLs []string `json:"callbackURLs"`
	CreatedBy    UserDto  `json:"createdBy"`
}

type OidcClientCreateDto struct {
	Name         string   `json:"name" binding:"required,max=50"`
	CallbackURLs []string `json:"callbackURLs" binding:"required,urlList"`
}

type AuthorizeOidcClientRequestDto struct {
	ClientID    string `json:"clientID" binding:"required"`
	Scope       string `json:"scope" binding:"required"`
	CallbackURL string `json:"callbackURL"`
	Nonce       string `json:"nonce"`
}

type AuthorizeOidcClientResponseDto struct {
	Code        string `json:"code"`
	CallbackURL string `json:"callbackURL"`
}

type OidcIdTokenDto struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code" binding:"required"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
}
