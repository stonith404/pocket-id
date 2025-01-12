package main

import "github.com/stonith404/pocket-id/backend/internal/bootstrap"

func main() {
	bootstrap.Bootstrap()
	// Uncomment the line below to only test the ldap functionality
	// ldap.GetLdapUser()
}
