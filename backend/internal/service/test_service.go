package service

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/fxamacker/cbor/v2"
	"log"
	"os"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"gorm.io/gorm"
)

type TestService struct {
	db               *gorm.DB
	appConfigService *AppConfigService
}

func NewTestService(db *gorm.DB, appConfigService *AppConfigService) *TestService {
	return &TestService{db: db, appConfigService: appConfigService}
}

func (s *TestService) SeedDatabase() error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		users := []model.User{
			{
				Base: model.Base{
					ID: "f4b89dc2-62fb-46bf-9f5f-c34f4eafe93e",
				},
				Username:  "tim",
				Email:     "tim.cook@test.com",
				FirstName: "Tim",
				LastName:  "Cook",
				IsAdmin:   true,
			},
			{
				Base: model.Base{
					ID: "1cd19686-f9a6-43f4-a41f-14a0bf5b4036",
				},
				Username:  "craig",
				Email:     "craig.federighi@test.com",
				FirstName: "Craig",
				LastName:  "Federighi",
				IsAdmin:   false,
			},
		}
		for _, user := range users {
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
		}

		userGroups := []model.UserGroup{
			{
				Base: model.Base{
					ID: "4110f814-56f1-4b28-8998-752b69bc97c0e",
				},
				Name:         "developers",
				FriendlyName: "Developers",
				Users:        []model.User{users[0], users[1]},
			},
			{
				Base: model.Base{
					ID: "adab18bf-f89d-4087-9ee1-70ff15b48211",
				},
				Name:         "designers",
				FriendlyName: "Designers",
				Users:        []model.User{users[0]},
			},
		}
		for _, group := range userGroups {
			if err := tx.Create(&group).Error; err != nil {
				return err
			}
		}

		oidcClients := []model.OidcClient{
			{
				Base: model.Base{
					ID: "3654a746-35d4-4321-ac61-0bdcff2b4055",
				},
				Name:         "Nextcloud",
				Secret:       "$2a$10$9dypwot8nGuCjT6wQWWpJOckZfRprhe2EkwpKizxS/fpVHrOLEJHC", // w2mUeZISmEvIDMEDvpY0PnxQIpj1m3zY
				CallbackURLs: model.CallbackURLs{"http://nextcloud/auth/callback"},
				ImageType:    utils.StringPointer("png"),
				CreatedByID:  users[0].ID,
			},
			{
				Base: model.Base{
					ID: "606c7782-f2b1-49e5-8ea9-26eb1b06d018",
				},
				Name:         "Immich",
				Secret:       "$2a$10$Ak.FP8riD1ssy2AGGbG.gOpnp/rBpymd74j0nxNMtW0GG1Lb4gzxe", // PYjrE9u4v9GVqXKi52eur0eb2Ci4kc0x
				CallbackURLs: model.CallbackURLs{"http://immich/auth/callback"},
				CreatedByID:  users[0].ID,
			},
		}
		for _, client := range oidcClients {
			if err := tx.Create(&client).Error; err != nil {
				return err
			}
		}

		authCode := model.OidcAuthorizationCode{
			Code:      "auth-code",
			Scope:     "openid profile",
			Nonce:     "nonce",
			ExpiresAt: time.Now().Add(1 * time.Hour),
			UserID:    users[0].ID,
			ClientID:  oidcClients[0].ID,
		}
		if err := tx.Create(&authCode).Error; err != nil {
			return err
		}

		accessToken := model.OneTimeAccessToken{
			Token:     "one-time-token",
			ExpiresAt: time.Now().Add(1 * time.Hour),
			UserID:    users[0].ID,
		}
		if err := tx.Create(&accessToken).Error; err != nil {
			return err
		}

		userAuthorizedClient := model.UserAuthorizedOidcClient{
			Scope:    "openid profile email",
			UserID:   users[0].ID,
			ClientID: oidcClients[0].ID,
		}
		if err := tx.Create(&userAuthorizedClient).Error; err != nil {
			return err
		}

		publicKey1, err := getCborPublicKey("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwcOo5KV169KR67QEHrcYkeXE3CCxv2BgwnSq4VYTQxyLtdmKxegexa8JdwFKhKXa2BMI9xaN15BoL6wSCRFJhg==")
		publicKey2, err := getCborPublicKey("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAESq/wR8QbBu3dKnpaw/v0mDxFFDwnJ/L5XHSg2tAmq5x1BpSMmIr3+DxCbybVvGRmWGh8kKhy7SMnK91M6rFHTA==")
		if err != nil {
			return err
		}
		webauthnCredentials := []model.WebauthnCredential{
			{
				Name:            "Passkey 1",
				CredentialID:    "test-credential-1",
				PublicKey:       publicKey1,
				AttestationType: "none",
				Transport:       model.AuthenticatorTransportList{protocol.Internal},
				UserID:          users[0].ID,
			},
			{
				Name:            "Passkey 2",
				CredentialID:    "test-credential-2",
				PublicKey:       publicKey2,
				AttestationType: "none",
				Transport:       model.AuthenticatorTransportList{protocol.Internal},
				UserID:          users[0].ID,
			},
		}
		for _, credential := range webauthnCredentials {
			if err := tx.Create(&credential).Error; err != nil {
				return err
			}
		}

		webauthnSession := model.WebauthnSession{
			Challenge:        "challenge",
			ExpiresAt:        time.Now().Add(1 * time.Hour),
			UserVerification: "preferred",
		}
		if err := tx.Create(&webauthnSession).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *TestService) ResetDatabase() error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var tables []string
		if err := tx.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name != 'schema_migrations';").Scan(&tables).Error; err != nil {
			return err
		}

		for _, table := range tables {
			if err := tx.Exec("DELETE FROM " + table).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = s.appConfigService.InitDbConfig()
	return err
}

func (s *TestService) ResetApplicationImages() error {
	if err := os.RemoveAll(common.EnvConfig.UploadPath); err != nil {
		log.Printf("Error removing directory: %v", err)
		return err
	}

	if err := utils.CopyDirectory("./images", common.EnvConfig.UploadPath+"/application-images"); err != nil {
		log.Printf("Error copying directory: %v", err)
		return err
	}

	return nil
}

// getCborPublicKey decodes a Base64 encoded public key and returns the CBOR encoded COSE key
func getCborPublicKey(base64PublicKey string) ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(base64PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key: %w", err)
	}
	pubKey, err := x509.ParsePKIXPublicKey(decodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	coseKey := map[int]interface{}{
		1:  2,                     // Key type: EC2
		3:  -7,                    // Algorithm: ECDSA with SHA-256
		-1: 1,                     // Curve: P-256
		-2: ecdsaPubKey.X.Bytes(), // X coordinate
		-3: ecdsaPubKey.Y.Bytes(), // Y coordinate
	}

	cborPublicKey, err := cbor.Marshal(coseKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal COSE key: %w", err)
	}

	return cborPublicKey, nil
}
