package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"github.com/stonith404/pocket-id/backend/internal/model"
	datatype "github.com/stonith404/pocket-id/backend/internal/model/types"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"slices"
	"strings"
	"time"
)

type OidcService struct {
	db                 *gorm.DB
	jwtService         *JwtService
	appConfigService   *AppConfigService
	auditLogService    *AuditLogService
	customClaimService *CustomClaimService
}

func NewOidcService(db *gorm.DB, jwtService *JwtService, appConfigService *AppConfigService, auditLogService *AuditLogService, customClaimService *CustomClaimService) *OidcService {
	return &OidcService{
		db:                 db,
		jwtService:         jwtService,
		appConfigService:   appConfigService,
		auditLogService:    auditLogService,
		customClaimService: customClaimService,
	}
}

func (s *OidcService) Authorize(input dto.AuthorizeOidcClientRequestDto, userID, ipAddress, userAgent string) (string, string, error) {
	var userAuthorizedOIDCClient model.UserAuthorizedOidcClient
	s.db.Preload("Client").First(&userAuthorizedOIDCClient, "client_id = ? AND user_id = ?", input.ClientID, userID)

	if userAuthorizedOIDCClient.Client.IsPublic && input.CodeChallenge == "" {
		return "", "", &common.OidcMissingCodeChallengeError{}
	}

	if userAuthorizedOIDCClient.Scope != input.Scope {
		return "", "", &common.OidcMissingAuthorizationError{}
	}

	callbackURL, err := s.getCallbackURL(userAuthorizedOIDCClient.Client, input.CallbackURL)
	if err != nil {
		return "", "", err
	}

	code, err := s.createAuthorizationCode(input.ClientID, userID, input.Scope, input.Nonce, input.CodeChallenge, input.CodeChallengeMethod)
	if err != nil {
		return "", "", err
	}

	s.auditLogService.Create(model.AuditLogEventClientAuthorization, ipAddress, userAgent, userID, model.AuditLogData{"clientName": userAuthorizedOIDCClient.Client.Name})

	return code, callbackURL, nil
}

func (s *OidcService) AuthorizeNewClient(input dto.AuthorizeOidcClientRequestDto, userID, ipAddress, userAgent string) (string, string, error) {
	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", input.ClientID).Error; err != nil {
		return "", "", err
	}

	if client.IsPublic && input.CodeChallenge == "" {
		return "", "", &common.OidcMissingCodeChallengeError{}
	}

	callbackURL, err := s.getCallbackURL(client, input.CallbackURL)
	if err != nil {
		return "", "", err
	}

	userAuthorizedClient := model.UserAuthorizedOidcClient{
		UserID:   userID,
		ClientID: input.ClientID,
		Scope:    input.Scope,
	}

	if err := s.db.Create(&userAuthorizedClient).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = s.db.Model(&userAuthorizedClient).Update("scope", input.Scope).Error
		} else {
			return "", "", err
		}
	}

	code, err := s.createAuthorizationCode(input.ClientID, userID, input.Scope, input.Nonce, input.CodeChallenge, input.CodeChallengeMethod)
	if err != nil {
		return "", "", err
	}

	s.auditLogService.Create(model.AuditLogEventNewClientAuthorization, ipAddress, userAgent, userID, model.AuditLogData{"clientName": client.Name})

	return code, callbackURL, nil
}

func (s *OidcService) CreateTokens(code, grantType, clientID, clientSecret, codeVerifier string) (string, string, error) {
	if grantType != "authorization_code" {
		return "", "", &common.OidcGrantTypeNotSupportedError{}
	}

	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return "", "", err
	}

	// Verify the client secret if the client is not public
	if !client.IsPublic {
		if clientID == "" || clientSecret == "" {
			return "", "", &common.OidcMissingClientCredentialsError{}
		}

		err := bcrypt.CompareHashAndPassword([]byte(client.Secret), []byte(clientSecret))
		if err != nil {
			return "", "", &common.OidcClientSecretInvalidError{}
		}
	}

	var authorizationCodeMetaData model.OidcAuthorizationCode
	err := s.db.Preload("User").First(&authorizationCodeMetaData, "code = ?", code).Error
	if err != nil {
		return "", "", &common.OidcInvalidAuthorizationCodeError{}
	}

	// If the client is public, the code verifier must match the code challenge
	if client.IsPublic {
		if !s.validateCodeVerifier(codeVerifier, *authorizationCodeMetaData.CodeChallenge, *authorizationCodeMetaData.CodeChallengeMethodSha256) {
			return "", "", &common.OidcInvalidCodeVerifierError{}
		}
	}

	if authorizationCodeMetaData.ClientID != clientID && authorizationCodeMetaData.ExpiresAt.ToTime().Before(time.Now()) {
		return "", "", &common.OidcInvalidAuthorizationCodeError{}
	}

	userClaims, err := s.GetUserClaimsForClient(authorizationCodeMetaData.UserID, clientID)
	if err != nil {
		return "", "", err
	}

	idToken, err := s.jwtService.GenerateIDToken(userClaims, clientID, authorizationCodeMetaData.Nonce)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.jwtService.GenerateOauthAccessToken(authorizationCodeMetaData.User, clientID)

	s.db.Delete(&authorizationCodeMetaData)

	return idToken, accessToken, nil
}

func (s *OidcService) GetClient(clientID string) (model.OidcClient, error) {
	var client model.OidcClient
	if err := s.db.Preload("CreatedBy").First(&client, "id = ?", clientID).Error; err != nil {
		return model.OidcClient{}, err
	}
	return client, nil
}

func (s *OidcService) ListClients(searchTerm string, page int, pageSize int) ([]model.OidcClient, utils.PaginationResponse, error) {
	var clients []model.OidcClient

	query := s.db.Preload("CreatedBy").Model(&model.OidcClient{})
	if searchTerm != "" {
		searchPattern := "%" + searchTerm + "%"
		query = query.Where("name LIKE ?", searchPattern)
	}

	pagination, err := utils.Paginate(page, pageSize, query, &clients)
	if err != nil {
		return nil, utils.PaginationResponse{}, err
	}

	return clients, pagination, nil
}

func (s *OidcService) CreateClient(input dto.OidcClientCreateDto, userID string) (model.OidcClient, error) {
	client := model.OidcClient{
		Name:         input.Name,
		CallbackURLs: input.CallbackURLs,
		CreatedByID:  userID,
	}

	if err := s.db.Create(&client).Error; err != nil {
		return model.OidcClient{}, err
	}

	return client, nil
}

func (s *OidcService) UpdateClient(clientID string, input dto.OidcClientCreateDto) (model.OidcClient, error) {
	var client model.OidcClient
	if err := s.db.Preload("CreatedBy").First(&client, "id = ?", clientID).Error; err != nil {
		return model.OidcClient{}, err
	}

	client.Name = input.Name
	client.CallbackURLs = input.CallbackURLs
	client.IsPublic = input.IsPublic

	if err := s.db.Save(&client).Error; err != nil {
		return model.OidcClient{}, err
	}

	return client, nil
}

func (s *OidcService) DeleteClient(clientID string) error {
	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&client).Error; err != nil {
		return err
	}

	return nil
}

func (s *OidcService) CreateClientSecret(clientID string) (string, error) {
	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return "", err
	}

	clientSecret, err := utils.GenerateRandomAlphanumericString(32)
	if err != nil {
		return "", err
	}

	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	client.Secret = string(hashedSecret)
	if err := s.db.Save(&client).Error; err != nil {
		return "", err
	}

	return clientSecret, nil
}

func (s *OidcService) GetClientLogo(clientID string) (string, string, error) {
	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return "", "", err
	}

	if client.ImageType == nil {
		return "", "", errors.New("image not found")
	}

	imageType := *client.ImageType
	imagePath := fmt.Sprintf("%s/oidc-client-images/%s.%s", common.EnvConfig.UploadPath, client.ID, imageType)
	mimeType := utils.GetImageMimeType(imageType)

	return imagePath, mimeType, nil
}

func (s *OidcService) UpdateClientLogo(clientID string, file *multipart.FileHeader) error {
	fileType := utils.GetFileExtension(file.Filename)
	if mimeType := utils.GetImageMimeType(fileType); mimeType == "" {
		return &common.FileTypeNotSupportedError{}
	}

	imagePath := fmt.Sprintf("%s/oidc-client-images/%s.%s", common.EnvConfig.UploadPath, clientID, fileType)
	if err := utils.SaveFile(file, imagePath); err != nil {
		return err
	}

	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return err
	}

	if client.ImageType != nil && fileType != *client.ImageType {
		oldImagePath := fmt.Sprintf("%s/oidc-client-images/%s.%s", common.EnvConfig.UploadPath, client.ID, *client.ImageType)
		if err := os.Remove(oldImagePath); err != nil {
			return err
		}
	}

	client.ImageType = &fileType
	if err := s.db.Save(&client).Error; err != nil {
		return err
	}

	return nil
}

func (s *OidcService) DeleteClientLogo(clientID string) error {
	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return err
	}

	if client.ImageType == nil {
		return errors.New("image not found")
	}

	imagePath := fmt.Sprintf("%s/oidc-client-images/%s.%s", common.EnvConfig.UploadPath, client.ID, *client.ImageType)
	if err := os.Remove(imagePath); err != nil {
		return err
	}

	client.ImageType = nil
	if err := s.db.Save(&client).Error; err != nil {
		return err
	}

	return nil
}

func (s *OidcService) GetUserClaimsForClient(userID string, clientID string) (map[string]interface{}, error) {
	var authorizedOidcClient model.UserAuthorizedOidcClient
	if err := s.db.Preload("User.UserGroups").First(&authorizedOidcClient, "user_id = ? AND client_id = ?", userID, clientID).Error; err != nil {
		return nil, err
	}

	user := authorizedOidcClient.User
	scope := authorizedOidcClient.Scope

	claims := map[string]interface{}{
		"sub": user.ID,
	}

	if strings.Contains(scope, "email") {
		claims["email"] = user.Email
		claims["email_verified"] = s.appConfigService.DbConfig.EmailsVerified.Value == "true"
	}

	if strings.Contains(scope, "groups") {
		userGroups := make([]string, len(user.UserGroups))
		for i, group := range user.UserGroups {
			userGroups[i] = group.Name
		}
		claims["groups"] = userGroups
	}

	profileClaims := map[string]interface{}{
		"given_name":         user.FirstName,
		"family_name":        user.LastName,
		"name":               user.FirstName + " " + user.LastName,
		"preferred_username": user.Username,
	}

	if strings.Contains(scope, "profile") {
		// Add profile claims
		for k, v := range profileClaims {
			claims[k] = v
		}

		// Add custom claims
		customClaims, err := s.customClaimService.GetCustomClaimsForUserWithUserGroups(userID)
		if err != nil {
			return nil, err
		}

		for _, customClaim := range customClaims {
			claims[customClaim.Key] = customClaim.Value
		}
	}
	if strings.Contains(scope, "email") {
		claims["email"] = user.Email
	}

	return claims, nil
}

func (s *OidcService) createAuthorizationCode(clientID string, userID string, scope string, nonce string, codeChallenge string, codeChallengeMethod string) (string, error) {
	randomString, err := utils.GenerateRandomAlphanumericString(32)
	if err != nil {
		return "", err
	}

	codeChallengeMethodSha256 := strings.ToUpper(codeChallengeMethod) == "S256"

	oidcAuthorizationCode := model.OidcAuthorizationCode{
		ExpiresAt:                 datatype.DateTime(time.Now().Add(15 * time.Minute)),
		Code:                      randomString,
		ClientID:                  clientID,
		UserID:                    userID,
		Scope:                     scope,
		Nonce:                     nonce,
		CodeChallenge:             &codeChallenge,
		CodeChallengeMethodSha256: &codeChallengeMethodSha256,
	}

	if err := s.db.Create(&oidcAuthorizationCode).Error; err != nil {
		return "", err
	}

	return randomString, nil
}

func (s *OidcService) validateCodeVerifier(codeVerifier, codeChallenge string, codeChallengeMethodSha256 bool) bool {
	if !codeChallengeMethodSha256 {
		return codeVerifier == codeChallenge
	}

	// Compute SHA-256 hash of the codeVerifier
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	codeVerifierHash := h.Sum(nil)

	// Base64 URL encode the verifier hash
	encodedVerifierHash := base64.RawURLEncoding.EncodeToString(codeVerifierHash)

	return encodedVerifierHash == codeChallenge
}

func (s *OidcService) getCallbackURL(client model.OidcClient, inputCallbackURL string) (callbackURL string, err error) {
	if inputCallbackURL == "" {
		return client.CallbackURLs[0], nil
	}
	if slices.Contains(client.CallbackURLs, inputCallbackURL) {
		return inputCallbackURL, nil
	}

	return "", &common.OidcInvalidCallbackURLError{}
}
