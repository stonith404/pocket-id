package dto

type PublicOidcClientDto struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	HasLogo bool   `json:"hasLogo"`
}

type OidcClientDto struct {
	PublicOidcClientDto
	CallbackURLs []string `json:"callbackURLs"`
	IsPublic     bool     `json:"isPublic"`
	PkceEnabled  bool     `json:"pkceEnabled"`
	CreatedBy    UserDto  `json:"createdBy"`
}

type OidcClientCreateDto struct {
	Name         string   `json:"name" binding:"required,max=50"`
	CallbackURLs []string `json:"callbackURLs" binding:"required,urlList"`
	IsPublic     bool     `json:"isPublic"`
	PkceEnabled  bool     `json:"pkceEnabled"`
}

type AuthorizeOidcClientRequestDto struct {
	ClientID            string `json:"clientID" binding:"required"`
	Scope               string `json:"scope" binding:"required"`
	CallbackURL         string `json:"callbackURL"`
	Nonce               string `json:"nonce"`
	CodeChallenge       string `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
}

type AuthorizeOidcClientResponseDto struct {
	Code        string `json:"code"`
	CallbackURL string `json:"callbackURL"`
}

type OidcCreateTokensDto struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code" binding:"required"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	CodeVerifier string `form:"code_verifier"`
}
