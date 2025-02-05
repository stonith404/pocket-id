package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pocket-id/pocket-id/backend/internal/common"
	"github.com/pocket-id/pocket-id/backend/internal/model"
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
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
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

func (s *JwtService) GenerateAccessToken(user model.User) (string, error) {
	sessionDurationInMinutes, _ := strconv.Atoi(s.appConfigService.DbConfig.SessionDuration.Value)
	claim := AccessTokenJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(sessionDurationInMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  jwt.ClaimStrings{common.EnvConfig.AppURL},
		},
		IsAdmin: user.IsAdmin,
	}

	kid, err := s.generateKeyID(s.publicKey)
	if err != nil {
		return "", errors.New("failed to generate key ID: " + err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	token.Header["kid"] = kid

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

	if !slices.Contains(claims.Audience, common.EnvConfig.AppURL) {
		return nil, errors.New("audience doesn't match")
	}
	return claims, nil
}

func (s *JwtService) GenerateIDToken(userClaims map[string]interface{}, clientID string, nonce string) (string, error) {
	claims := jwt.MapClaims{
		"aud": clientID,
		"exp": jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
		"iss": common.EnvConfig.AppURL,
	}

	for k, v := range userClaims {
		claims[k] = v
	}

	if nonce != "" {
		claims["nonce"] = nonce
	}

	kid, err := s.generateKeyID(s.publicKey)
	if err != nil {
		return "", errors.New("failed to generate key ID: " + err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	return token.SignedString(s.privateKey)
}

func (s *JwtService) GenerateOauthAccessToken(user model.User, clientID string) (string, error) {
	claim := jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Audience:  jwt.ClaimStrings{clientID},
		Issuer:    common.EnvConfig.AppURL,
	}

	kid, err := s.generateKeyID(s.publicKey)
	if err != nil {
		return "", errors.New("failed to generate key ID: " + err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	token.Header["kid"] = kid

	return token.SignedString(s.privateKey)
}

func (s *JwtService) VerifyOauthAccessToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("couldn't handle this token")
	}

	claims, isValid := token.Claims.(*jwt.RegisteredClaims)
	if !isValid {
		return nil, errors.New("can't parse claims")
	}

	return claims, nil
}

// GetJWK returns the JSON Web Key (JWK) for the public key.
func (s *JwtService) GetJWK() (JWK, error) {
	if s.publicKey == nil {
		return JWK{}, errors.New("public key is not initialized")
	}

	kid, err := s.generateKeyID(s.publicKey)
	if err != nil {
		return JWK{}, err
	}

	jwk := JWK{
		Kid: kid,
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(s.publicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(s.publicKey.E)).Bytes()),
	}

	return jwk, nil
}

// GenerateKeyID generates a Key ID for the public key using the first 8 bytes of the SHA-256 hash of the public key.
func (s *JwtService) generateKeyID(publicKey *rsa.PublicKey) (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", errors.New("failed to marshal public key: " + err.Error())
	}

	// Compute SHA-256 hash of the public key
	hash := sha256.New()
	hash.Write(pubASN1)
	hashed := hash.Sum(nil)

	// Truncate the hash to the first 8 bytes for a shorter Key ID
	shortHash := hashed[:8]

	// Return Base64 encoded truncated hash as Key ID
	return base64.RawURLEncoding.EncodeToString(shortHash), nil
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
