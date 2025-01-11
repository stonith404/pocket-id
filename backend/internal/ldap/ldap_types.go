package ldap

type LDAPUserSeachResult struct {
	DN                string   `ldap:"dn"`
	UserPrincipalName string   `ldap:"userPrincipalName"`
	Username          string   `ldap:"sAMAccountName"`
	Mail              string   `ldap:"mail"`
	MemberOf          []string `ldap:"memberOf"`
	GivenName         string   `ldap:"givenName"`
	LastName          string   `ldap:"sn"`
}
