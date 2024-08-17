package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	privateKeyPath = "data/keys/jwt_private_key.pem"
	publicKeyPath  = "data/keys/jwt_public_key.pem"
)

type JwtService struct {
	publicKey        *rsa.PublicKey
	privateKey       *rsa.PrivateKey
	appConfigService *AppConfigService
}

func NewJwtService(appConfigService *AppConfigService) *JwtService {
	service := &JwtService{
		appConfigService: appConfigService,
	}

	// Ensure keys are generated or loaded
	if err := service.loadOrGenerateKeys(); err != nil {
		log.Fatalf("Failed to initialize jwt service: %v", err)
	}

	return service
}

type AccessTokenJWTClaims struct {
	jwt.RegisteredClaims
	IsAdmin bool `json:"isAdmin,omitempty"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// loadOrGenerateKeys loads RSA keys from the given paths or generates them if they do not exist.
func (s *JwtService) loadOrGenerateKeys() error {
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		if err := s.generateKeys(); err != nil {
			return err
		}
	}

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return errors.New("can't read jwt private key: " + err.Error())
	}
	s.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return errors.New("can't parse jwt private key: " + err.Error())
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return errors.New("can't read jwt public key: " + err.Error())
	}
	s.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return errors.New("can't parse jwt public key: " + err.Error())
	}

	return nil
}

func (s *JwtService) GenerateIDToken(user model.User, clientID string, scope string, nonce string) (string, error) {
	profileClaims := map[string]interface{}{
		"given_name":         user.FirstName,
		"family_name":        user.LastName,
		"email":              user.Email,
		"preferred_username": user.Username,
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"aud": clientID,
		"exp": jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
	}

	if nonce != "" {
		claims["nonce"] = nonce
	}
	if strings.Contains(scope, "profile") {
		for k, v := range profileClaims {
			claims[k] = v
		}
	}
	if strings.Contains(scope, "email") {
		claims["email"] = user.Email
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *JwtService) GenerateAccessToken(user model.User) (string, error) {
	sessionDurationInMinutes, _ := strconv.Atoi(s.appConfigService.DbConfig.SessionDuration.Value)
	claim := AccessTokenJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(sessionDurationInMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  jwt.ClaimStrings{utils.GetHostFromURL(common.EnvConfig.AppURL)},
		},
		IsAdmin: user.IsAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	return token.SignedString(s.privateKey)
}

func (s *JwtService) VerifyAccessToken(tokenString string) (*AccessTokenJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("couldn't handle this token")
	}

	claims, isValid := token.Claims.(*AccessTokenJWTClaims)
	if !isValid {
		return nil, errors.New("can't parse claims")
	}

	if !slices.Contains(claims.Audience, utils.GetHostFromURL(common.EnvConfig.AppURL)) {
		return nil, errors.New("audience doesn't match")
	}
	return claims, nil
}

// GetJWK returns the JSON Web Key (JWK) for the public key.
func (s *JwtService) GetJWK() (JWK, error) {
	if s.publicKey == nil {
		return JWK{}, errors.New("public key is not initialized")
	}

	jwk := JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: "1",
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(s.publicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(s.publicKey.E)).Bytes()),
	}

	return jwk, nil
}

// generateKeys generates a new RSA key pair and saves them to the specified paths.
func (s *JwtService) generateKeys() error {
	if err := os.MkdirAll(filepath.Dir(privateKeyPath), 0700); err != nil {
		return errors.New("failed to create directories for keys: " + err.Error())
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.New("failed to generate private key: " + err.Error())
	}
	s.privateKey = privateKey

	if err := s.savePEMKey(privateKeyPath, x509.MarshalPKCS1PrivateKey(privateKey), "RSA PRIVATE KEY"); err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey
	s.publicKey = publicKey

	if err := s.savePEMKey(publicKeyPath, x509.MarshalPKCS1PublicKey(publicKey), "RSA PUBLIC KEY"); err != nil {
		return err
	}

	return nil
}

// savePEMKey saves a PEM encoded key to a file.
func (s *JwtService) savePEMKey(path string, keyBytes []byte, keyType string) error {
	keyFile, err := os.Create(path)
	if err != nil {
		return errors.New("failed to create key file: " + err.Error())
	}
	defer keyFile.Close()

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  keyType,
		Bytes: keyBytes,
	})

	if _, err := keyFile.Write(keyPEM); err != nil {
		return errors.New("failed to write key file: " + err.Error())
	}

	return nil
}

// loadKeys loads RSA keys from the given paths.
func (s *JwtService) loadKeys() error {
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		if err := s.generateKeys(); err != nil {
			return err
		}
	}

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("can't read jwt private key: %w", err)
	}
	s.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return fmt.Errorf("can't parse jwt private key: %w", err)
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("can't read jwt public key: %w", err)
	}
	s.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return fmt.Errorf("can't parse jwt public key: %w", err)
	}

	return nil
}
