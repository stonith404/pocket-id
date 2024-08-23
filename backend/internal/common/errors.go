package common

import "errors"

var (
	ErrUsernameTaken                = errors.New("username is already taken")
	ErrEmailTaken                   = errors.New("email is already taken")
	ErrSetupAlreadyCompleted        = errors.New("setup already completed")
	ErrTokenInvalidOrExpired        = errors.New("token is invalid or expired")
	ErrOidcMissingAuthorization     = errors.New("missing authorization")
	ErrOidcGrantTypeNotSupported    = errors.New("grant type not supported")
	ErrOidcMissingClientCredentials = errors.New("client id or secret not provided")
	ErrOidcClientSecretInvalid      = errors.New("invalid client secret")
	ErrOidcInvalidAuthorizationCode = errors.New("invalid authorization code")
	ErrOidcInvalidCallbackURL       = errors.New("invalid callback URL")
	ErrFileTypeNotSupported         = errors.New("file type not supported")
	ErrInvalidCredentials           = errors.New("no user found with provided credentials")
)
