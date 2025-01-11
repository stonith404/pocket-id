package main

import (
	"github.com/stonith404/pocket-id/backend/internal/ldap"
)

func main() {
	// bootstrap.Bootstrap()
	ldap.GetLdapUser()
}
