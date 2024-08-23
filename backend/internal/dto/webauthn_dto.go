package dto

import (
	"github.com/go-webauthn/webauthn/protocol"
	"time"
)

type WebauthnCredentialDto struct {
	ID              string                            `json:"id"`
	Name            string                            `json:"name"`
	CredentialID    string                            `json:"credentialID"`
	AttestationType string                            `json:"attestationType"`
	Transport       []protocol.AuthenticatorTransport `json:"transport"`

	BackupEligible bool `json:"backupEligible"`
	BackupState    bool `json:"backupState"`

	CreatedAt time.Time `json:"createdAt"`
}

type WebauthnCredentialUpdateDto struct {
	Name string `json:"name" binding:"required,min=1,max=30"`
}
