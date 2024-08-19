package service

import (
	"errors"
	"fmt"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type OidcService struct {
	db         *gorm.DB
	jwtService *JwtService
}

func NewOidcService(db *gorm.DB, jwtService *JwtService) *OidcService {
	return &OidcService{
		db:         db,
		jwtService: jwtService,
	}
}

func (s *OidcService) Authorize(req model.AuthorizeRequest, userID string) (string, error) {
	var userAuthorizedOIDCClient model.UserAuthorizedOidcClient
	s.db.First(&userAuthorizedOIDCClient, "client_id = ? AND user_id = ?", req.ClientID, userID)

	if userAuthorizedOIDCClient.Scope != req.Scope {
		return "", common.ErrOidcMissingAuthorization
	}

	return s.createAuthorizationCode(req.ClientID, userID, req.Scope, req.Nonce)
}

func (s *OidcService) AuthorizeNewClient(req model.AuthorizeNewClientDto, userID string) (string, error) {
	userAuthorizedClient := model.UserAuthorizedOidcClient{
		UserID:   userID,
		ClientID: req.ClientID,
		Scope:    req.Scope,
	}

	if err := s.db.Create(&userAuthorizedClient).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = s.db.Model(&userAuthorizedClient).Update("scope", req.Scope).Error
		} else {
			return "", err
		}
	}

	return s.createAuthorizationCode(req.ClientID, userID, req.Scope, req.Nonce)
}

func (s *OidcService) CreateTokens(code, grantType, clientID, clientSecret string) (string, string, error) {
	if grantType != "authorization_code" {
		return "", "", common.ErrOidcGrantTypeNotSupported
	}

	if clientID == "" || clientSecret == "" {
		return "", "", common.ErrOidcMissingClientCredentials
	}

	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return "", "", err
	}

	err := bcrypt.CompareHashAndPassword([]byte(client.Secret), []byte(clientSecret))
	if err != nil {
		return "", "", common.ErrOidcClientSecretInvalid
	}

	var authorizationCodeMetaData model.OidcAuthorizationCode
	err = s.db.Preload("User").First(&authorizationCodeMetaData, "code = ?", code).Error
	if err != nil {
		return "", "", common.ErrOidcInvalidAuthorizationCode
	}

	if authorizationCodeMetaData.ClientID != clientID && authorizationCodeMetaData.ExpiresAt.Before(time.Now()) {
		return "", "", common.ErrOidcInvalidAuthorizationCode
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

func (s *OidcService) GetClient(clientID string) (*model.OidcClient, error) {
	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (s *OidcService) ListClients(searchTerm string, page int, pageSize int) ([]model.OidcClient, utils.PaginationResponse, error) {
	var clients []model.OidcClient

	query := s.db.Model(&model.OidcClient{})
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

func (s *OidcService) CreateClient(input model.OidcClientCreateDto, userID string) (*model.OidcClient, error) {
	client := model.OidcClient{
		Name:        input.Name,
		CallbackURL: input.CallbackURL,
		CreatedByID: userID,
	}

	if err := s.db.Create(&client).Error; err != nil {
		return nil, err
	}

	return &client, nil
}

func (s *OidcService) UpdateClient(clientID string, input model.OidcClientCreateDto) (*model.OidcClient, error) {
	var client model.OidcClient
	if err := s.db.First(&client, "id = ?", clientID).Error; err != nil {
		return nil, err
	}

	client.Name = input.Name
	client.CallbackURL = input.CallbackURL

	if err := s.db.Save(&client).Error; err != nil {
		return nil, err
	}

	return &client, nil
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
		return common.ErrFileTypeNotSupported
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
	if err := s.db.Preload("User").First(&authorizedOidcClient, "user_id = ? AND client_id = ?", userID, clientID).Error; err != nil {
		return nil, err
	}

	user := authorizedOidcClient.User
	scope := authorizedOidcClient.Scope

	claims := map[string]interface{}{
		"sub": user.ID,
	}

	if strings.Contains(scope, "email") {
		claims["email"] = user.Email
	}

	profileClaims := map[string]interface{}{
		"given_name":         user.FirstName,
		"family_name":        user.LastName,
		"preferred_username": user.Username,
	}

	if strings.Contains(scope, "profile") {
		for k, v := range profileClaims {
			claims[k] = v
		}
	}
	if strings.Contains(scope, "email") {
		claims["email"] = user.Email
	}

	return claims, nil
}

func (s *OidcService) createAuthorizationCode(clientID string, userID string, scope string, nonce string) (string, error) {
	randomString, err := utils.GenerateRandomAlphanumericString(32)
	if err != nil {
		return "", err
	}

	oidcAuthorizationCode := model.OidcAuthorizationCode{
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Code:      randomString,
		ClientID:  clientID,
		UserID:    userID,
		Scope:     scope,
		Nonce:     nonce,
	}

	if err := s.db.Create(&oidcAuthorizationCode).Error; err != nil {
		return "", err
	}

	return randomString, nil
}
