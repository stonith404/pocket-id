package common

import (
	"fmt"
	"net/http"
)

type AppError interface {
	error
	HttpStatusCode() int
}

// Custom error types for various conditions

type AlreadyInUseError struct {
	Property string
}

func (e *AlreadyInUseError) Error() string {
	return fmt.Sprintf("%s is already in use", e.Property)
}
func (e *AlreadyInUseError) HttpStatusCode() int { return 400 }

type SetupAlreadyCompletedError struct{}

func (e *SetupAlreadyCompletedError) Error() string       { return "setup already completed" }
func (e *SetupAlreadyCompletedError) HttpStatusCode() int { return 400 }

type TokenInvalidOrExpiredError struct{}

func (e *TokenInvalidOrExpiredError) Error() string       { return "token is invalid or expired" }
func (e *TokenInvalidOrExpiredError) HttpStatusCode() int { return 400 }

type OidcMissingAuthorizationError struct{}

func (e *OidcMissingAuthorizationError) Error() string       { return "missing authorization" }
func (e *OidcMissingAuthorizationError) HttpStatusCode() int { return http.StatusForbidden }

type OidcGrantTypeNotSupportedError struct{}

func (e *OidcGrantTypeNotSupportedError) Error() string       { return "grant type not supported" }
func (e *OidcGrantTypeNotSupportedError) HttpStatusCode() int { return 400 }

type OidcMissingClientCredentialsError struct{}

func (e *OidcMissingClientCredentialsError) Error() string       { return "client id or secret not provided" }
func (e *OidcMissingClientCredentialsError) HttpStatusCode() int { return 400 }

type OidcClientSecretInvalidError struct{}

func (e *OidcClientSecretInvalidError) Error() string       { return "invalid client secret" }
func (e *OidcClientSecretInvalidError) HttpStatusCode() int { return 400 }

type OidcInvalidAuthorizationCodeError struct{}

func (e *OidcInvalidAuthorizationCodeError) Error() string       { return "invalid authorization code" }
func (e *OidcInvalidAuthorizationCodeError) HttpStatusCode() int { return 400 }

type OidcInvalidCallbackURLError struct{}

func (e *OidcInvalidCallbackURLError) Error() string       { return "invalid callback URL" }
func (e *OidcInvalidCallbackURLError) HttpStatusCode() int { return 400 }

type FileTypeNotSupportedError struct{}

func (e *FileTypeNotSupportedError) Error() string       { return "file type not supported" }
func (e *FileTypeNotSupportedError) HttpStatusCode() int { return 400 }

type InvalidCredentialsError struct{}

func (e *InvalidCredentialsError) Error() string       { return "no user found with provided credentials" }
func (e *InvalidCredentialsError) HttpStatusCode() int { return 400 }

type FileTooLargeError struct {
	MaxSize string
}

func (e *FileTooLargeError) Error() string {
	return fmt.Sprintf("The file can't be larger than %s", e.MaxSize)
}
func (e *FileTooLargeError) HttpStatusCode() int { return http.StatusRequestEntityTooLarge }

type NotSignedInError struct{}

func (e *NotSignedInError) Error() string       { return "You are not signed in" }
func (e *NotSignedInError) HttpStatusCode() int { return http.StatusUnauthorized }

type MissingPermissionError struct{}

func (e *MissingPermissionError) Error() string {
	return "You don't have permission to perform this action"
}
func (e *MissingPermissionError) HttpStatusCode() int { return http.StatusForbidden }

type TooManyRequestsError struct{}

func (e *TooManyRequestsError) Error() string {
	return "Too many requests. Please wait a while before trying again."
}
func (e *TooManyRequestsError) HttpStatusCode() int { return http.StatusTooManyRequests }

type ClientIdOrSecretNotProvidedError struct{}

func (e *ClientIdOrSecretNotProvidedError) Error() string {
	return "Client id and secret not provided"
}

func (e *ClientIdOrSecretNotProvidedError) HttpStatusCode() int { return http.StatusBadRequest }

type WrongFileTypeError struct {
	ExpectedFileType string
}

func (e *WrongFileTypeError) Error() string {
	return fmt.Sprintf("File must be of type %s", e.ExpectedFileType)
}

func (e *WrongFileTypeError) HttpStatusCode() int { return http.StatusBadRequest }

type MissingSessionIdError struct{}

func (e *MissingSessionIdError) Error() string {
	return "Missing session id"
}

func (e *MissingSessionIdError) HttpStatusCode() int { return http.StatusBadRequest }

type ReservedClaimError struct {
	Key string
}

func (e *ReservedClaimError) Error() string {
	return fmt.Sprintf("Claim %s is reserved and can't be used", e.Key)
}

func (e *ReservedClaimError) HttpStatusCode() int { return http.StatusBadRequest }

type DuplicateClaimError struct {
	Key string
}

func (e *DuplicateClaimError) Error() string {
	return fmt.Sprintf("Claim %s is already defined", e.Key)
}

func (e *DuplicateClaimError) HttpStatusCode() int { return http.StatusBadRequest }
