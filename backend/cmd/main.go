package main

import (
	"github.com/stonith404/pocket-id/backend/internal/ldap"
)

func main() {
	// bootstrap.Bootstrap()
	// this is for testing only so its easier to debug
	ldap.GetLdapUser()
}
