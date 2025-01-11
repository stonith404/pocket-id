package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/stonith404/pocket-id/backend/internal/common"
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

func printGreen(text string) string {
	return fmt.Sprintf("\033[92m%s\033[0m", text)
}

func GetLdapUser() LDAPUserSeachResult {
	client := ldapInit()
	// user := username
	baseDN := common.EnvConfig.LDAPSearchBase
	filter := "(objectClass=person)"

	searchAttrs := []string{
		"sAMAccountName",
		"mail",
		"memberOf",
		"userPrincipalName",
		"givenName",
		"sn",
	}

	// Filters must start and finish with ()!
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, searchAttrs, []ldap.Control{})

	result, err := client.Search(searchReq)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to query LDAP: %w", err))
	}

	userResult := LDAPUserSeachResult{}

	if len(result.Entries) >= 1 {

		if err := result.Entries[0].Unmarshal(&userResult); err != nil {
			panic(err)
		}

		for _, value := range result.Entries {
			if err := value.Unmarshal(&userResult); err != nil {
				panic(err)
			}
			fmt.Println("\nUser Attributes:")
			fmt.Printf("Full Name: %s\n", printGreen(userResult.GivenName+" "+userResult.LastName))
			fmt.Printf("Email: %s\n", printGreen(userResult.Mail))
			fmt.Printf("Username: %s\n", printGreen(userResult.Username))
			fmt.Printf("DN: %s\n", printGreen(userResult.DN))
		}

		return userResult

	} else {
		fmt.Println("No Users Found")
		panic(1)
	}
}
