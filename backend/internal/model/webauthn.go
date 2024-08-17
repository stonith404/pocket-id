package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/go-webauthn/webauthn/protocol"
	"time"
)

type WebauthnSession struct {
	Base

	Challenge        string
	ExpiresAt        time.Time
	UserVerification string
}

type WebauthnCredential struct {
	Base

	Name            string                     `json:"name"`
	CredentialID    string                     `json:"credentialID"`
	PublicKey       []byte                     `json:"-"`
	AttestationType string                     `json:"attestationType"`
	Transport       AuthenticatorTransportList `json:"-"`

	BackupEligible bool `json:"backupEligible"`
	BackupState    bool `json:"backupState"`

	UserID string
}

type PublicKeyCredentialCreationOptions struct {
	Response  protocol.PublicKeyCredentialCreationOptions `json:"response"`
	SessionID string                                      `json:"session_id"`
	Timeout   time.Duration                               `json:"timeout"`
}

type PublicKeyCredentialRequestOptions struct {
	Response  protocol.PublicKeyCredentialRequestOptions `json:"response"`
	SessionID string                                     `json:"session_id"`
	Timeout   time.Duration                              `json:"timeout"`
}

type AuthenticatorTransportList []protocol.AuthenticatorTransport

// Scan and Value methods for GORM to handle the custom type
func (atl *AuthenticatorTransportList) Scan(value interface{}) error {

	if v, ok := value.([]byte); ok {
		return json.Unmarshal(v, atl)
	} else {
		return errors.New("type assertion to []byte failed")
	}
}

func (atl AuthenticatorTransportList) Value() (driver.Value, error) {
	return json.Marshal(atl)
}

type WebauthnCredentialUpdateDto struct {
	Name string `json:"name"`
}
