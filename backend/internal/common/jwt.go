package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang-rest-api-template/internal/model"
	"golang-rest-api-template/internal/utils"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

const (
	privateKeyPath = "data/keys/jwt_private_key.pem"
	publicKeyPath  = "data/keys/jwt_public_key.pem"
)

type accessTokenJWTClaims struct {
	jwt.RegisteredClaims
	IsAdmin bool `json:"isAdmin,omitempty"`
}

// GenerateIDToken generates an ID token for the given user, clientID, scope and nonce.
func GenerateIDToken(user model.User, clientID string, scope string, nonce string) (tokenString string, err error) {
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
	signedToken, err := token.SignedString(PrivateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateAccessToken generates an access token for the given user.
func GenerateAccessToken(user model.User) (tokenString string, err error) {
	sessionDurationInMinutes, _ := strconv.Atoi(DbConfig.SessionDuration.Value)
	claim := accessTokenJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(sessionDurationInMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  jwt.ClaimStrings{utils.GetHostFromURL(EnvConfig.AppURL)},
		},
		IsAdmin: user.IsAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	tokenString, err = token.SignedString(PrivateKey)
	return tokenString, err
}

// VerifyAccessToken verifies the given access token and returns the claims if the token is valid.
func VerifyAccessToken(tokenString string) (*accessTokenJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &accessTokenJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return PublicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("couldn't handle this token")
	}

	claims, isValid := token.Claims.(*accessTokenJWTClaims)
	if !isValid {
		return nil, errors.New("can't parse claims")
	}

	if !slices.Contains(claims.Audience, utils.GetHostFromURL(EnvConfig.AppURL)) {
		return nil, errors.New("audience doesn't match")
	}
	return claims, nil
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// GetJWK returns the JSON Web Key (JWK) for the public key.
func GetJWK() (JWK, error) {
	if PublicKey == nil {
		return JWK{}, errors.New("public key is not initialized")
	}

	// Create JWK from RSA public key
	jwk := JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: "1", // Key ID can be set to any identifier. Here it's statically set to "1"
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(PublicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(PublicKey.E)).Bytes()),
	}

	return jwk, nil
}

// generateKeys generates a new RSA key pair and saves the private and public keys to the data folder.
func generateKeys() {
	if err := os.MkdirAll(filepath.Dir(privateKeyPath), 0700); err != nil {
		log.Fatal("Failed to create directories for keys", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Failed to generate private key", err)
	}

	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		log.Fatal("Failed to create private key file", err)
	}
	defer privateKeyFile.Close()

	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	_, err = privateKeyFile.Write(privateKeyPEM)
	if err != nil {
		log.Fatal("Failed to write private key file", err)
	}

	publicKey := &privateKey.PublicKey
	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		log.Fatal("Failed to create public key file", err)
	}
	defer publicKeyFile.Close()

	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(publicKey),
		},
	)
	_, err = publicKeyFile.Write(publicKeyPEM)
	if err != nil {
		log.Fatal("Failed to write public key file", err)
	}
}

func init() {
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		generateKeys()
	}

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Can't read jwt private key", err)
	}
	PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Fatal("Can't parse jwt private key", err)
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal("Can't read jwt public key", err)
	}
	PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Fatal("Can't parse jwt public key", err)
	}
}
