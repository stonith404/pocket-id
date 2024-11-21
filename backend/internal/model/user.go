package model

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stonith404/pocket-id/backend/internal/model/types"
)

type User struct {
	Base

	Username  string
	Email     string
	FirstName string
	LastName  string
	IsAdmin   bool

	CustomClaims []CustomClaim
	UserGroups   []UserGroup `gorm:"many2many:user_groups_users;"`
	Credentials  []WebauthnCredential
}

func (u User) WebAuthnID() []byte { return []byte(u.ID) }

func (u User) WebAuthnName() string { return u.Username }

func (u User) WebAuthnDisplayName() string { return u.FirstName + " " + u.LastName }

func (u User) WebAuthnIcon() string { return "" }

func (u User) WebAuthnCredentials() []webauthn.Credential {
	credentials := make([]webauthn.Credential, len(u.Credentials))

	for i, credential := range u.Credentials {
		credentials[i] = webauthn.Credential{
			ID:              []byte(credential.CredentialID),
			AttestationType: credential.AttestationType,
			PublicKey:       credential.PublicKey,
			Transport:       credential.Transport,
			Flags: webauthn.CredentialFlags{
				BackupState:    credential.BackupState,
				BackupEligible: credential.BackupEligible,
			},
		}

	}
	return credentials
}

func (u User) WebAuthnCredentialDescriptors() (descriptors []protocol.CredentialDescriptor) {
	credentials := u.WebAuthnCredentials()

	descriptors = make([]protocol.CredentialDescriptor, len(credentials))

	for i, credential := range credentials {
		descriptors[i] = credential.Descriptor()
	}

	return descriptors
}

func (u User) FullName() string { return u.FirstName + " " + u.LastName }

type OneTimeAccessToken struct {
	Base
	Token     string
	ExpiresAt datatype.DateTime

	UserID string
	User   User
}
