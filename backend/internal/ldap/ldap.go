package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/model"
)

func ldapInit() *ldap.Conn {
	// Setup AD Connection
	ldapURL := common.EnvConfig.LDAPServer + ":" + common.EnvConfig.LDAPPort
	client, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: common.EnvConfig.LDAPTLSVerify}))
	if err != nil {
		//TODO Handle Errors Better
		panic(err)
	}
	// defer client.Close()

	// Bind as Service Account
	err = client.Bind(common.EnvConfig.LDAPBindUser, common.EnvConfig.LDAPBindPassword)
	if err != nil {
		//TODO Handle Errors Better
		panic(err)
	}
	return client
}

func GetLdapGroups() []model.UserGroup {
	client := ldapInit()
	baseDN := common.EnvConfig.LDAPSearchBase
	filter := "(objectClass=groupOfUniqueNames)"

	searchAttrs := []string{
		common.EnvConfig.LDAPGroupAttribute,
		"member",
	}

	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, searchAttrs, []ldap.Control{})
	result, err := client.Search(searchReq)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to query LDAP: %w", err))
	}

	if len(result.Entries) >= 1 {

		var ldapGroups []model.UserGroup
		for _, value := range result.Entries {
			group := model.UserGroup{
				Name: value.GetAttributeValue(common.EnvConfig.LDAPGroupAttribute),
			}
			ldapGroups = append(ldapGroups, group)
		}

		client.Close()

		return ldapGroups
	} else {
		fmt.Println("No Groups Found")
		panic(1)
	}

}

func GetLdapUsers() []model.User {
	client := ldapInit()
	// user := username
	baseDN := common.EnvConfig.LDAPSearchBase
	filter := "(objectClass=person)"

	//TODO Make options in UI to configure what options should be synced etc etc, as this depends on what ldap backend is being used.
	searchAttrs := []string{
		"mail",
		"memberOf",
		"givenName",
		"sn",
		"cn",
		common.EnvConfig.LDAPUsernameAttribute, // Search for the Username Attribute supplied by the user.
	}

	// Filters must start and finish with ()!
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, searchAttrs, []ldap.Control{})

	result, err := client.Search(searchReq)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to query LDAP: %w", err))
	}

	if len(result.Entries) >= 1 {

		var ldapUsers []model.User
		for _, value := range result.Entries {
			user := model.User{
				Username:  value.GetAttributeValue(common.EnvConfig.LDAPUsernameAttribute),
				Email:     value.GetAttributeValue("mail"),
				FirstName: value.GetAttributeValue("givenName"),
				LastName:  value.GetAttributeValue("sn"),
			}
			ldapUsers = append(ldapUsers, user)
		}

		for _, user := range ldapUsers {
			fmt.Printf("Username: %s\n", user.Username)
			fmt.Printf("First Name: %s\n", user.FirstName)
			fmt.Printf("Last Name: %s\n", user.LastName)
			fmt.Printf("Email: %s\n", user.Email)
			fmt.Printf("Admin: %t\n", user.IsAdmin)
		}

		client.Close()
		return ldapUsers

	} else {
		fmt.Println("No Users Found")
		panic(1)
	}

}
