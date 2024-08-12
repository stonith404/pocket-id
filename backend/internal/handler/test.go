package handler

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"golang-rest-api-template/internal/common"
	"golang-rest-api-template/internal/model"
	"golang-rest-api-template/internal/utils"
	"gorm.io/gorm"
)

func RegisterTestRoutes(group *gin.RouterGroup) {
	group.POST("/test/reset", resetAndSeedHandler)
}

func resetAndSeedHandler(c *gin.Context) {
	if err := resetDatabase(); err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	if err := resetApplicationImages(); err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	if err := seedDatabase(); err != nil {
		utils.UnknownHandlerError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Database reset and seeded"})
}

// seedDatabase seeds the database with initial data and uses a transaction to ensure atomicity.
func seedDatabase() error {
	return common.DB.Transaction(func(tx *gorm.DB) error {
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

		oidcClients := []model.OidcClient{
			{
				Base: model.Base{
					ID: "3654a746-35d4-4321-ac61-0bdcff2b4055",
				},
				Name:        "Nextcloud",
				Secret:      "$2a$10$9dypwot8nGuCjT6wQWWpJOckZfRprhe2EkwpKizxS/fpVHrOLEJHC", // w2mUeZISmEvIDMEDvpY0PnxQIpj1m3zY
				CallbackURL: "http://nextcloud/auth/callback",
				ImageType:   utils.StringPointer("png"),
				CreatedByID: users[0].ID,
			},
			{
				Base: model.Base{
					ID: "606c7782-f2b1-49e5-8ea9-26eb1b06d018",
				},
				Name:        "Immich",
				Secret:      "$2a$10$Ak.FP8riD1ssy2AGGbG.gOpnp/rBpymd74j0nxNMtW0GG1Lb4gzxe", // PYjrE9u4v9GVqXKi52eur0eb2Ci4kc0x
				CallbackURL: "http://immich/auth/callback",
				CreatedByID: users[0].ID,
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

		webauthnCredentials := []model.WebauthnCredential{
			{
				Name:            "Passkey 1",
				CredentialID:    "test-credential-1",
				PublicKey:       getCborPublicKey("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwcOo5KV169KR67QEHrcYkeXE3CCxv2BgwnSq4VYTQxyLtdmKxegexa8JdwFKhKXa2BMI9xaN15BoL6wSCRFJhg=="),
				AttestationType: "none",
				Transport:       model.AuthenticatorTransportList{protocol.Internal},
				UserID:          users[0].ID,
			},
			{
				Name:            "Passkey 2",
				CredentialID:    "test-credential-2",
				PublicKey:       getCborPublicKey("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAESq/wR8QbBu3dKnpaw/v0mDxFFDwnJ/L5XHSg2tAmq5x1BpSMmIr3+DxCbybVvGRmWGh8kKhy7SMnK91M6rFHTA=="),
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

// resetDatabase resets the database by deleting all rows from each table.
func resetDatabase() error {
	err := common.DB.Transaction(func(tx *gorm.DB) error {
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
	common.InitDbConfig()
	return nil
}

// resetApplicationImages resets the application images by removing existing images and replacing them with the default ones
func resetApplicationImages() error {

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
func getCborPublicKey(base64PublicKey string) []byte {
	decodedKey, err := base64.StdEncoding.DecodeString(base64PublicKey)
	if err != nil {
		log.Fatalf("Failed to decode base64 key: %v", err)
	}

	pubKey, err := x509.ParsePKIXPublicKey(decodedKey)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Not an ECDSA public key")
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
		log.Fatalf("Failed to encode CBOR: %v", err)
	}

	return cborPublicKey
}
