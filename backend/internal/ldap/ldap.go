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

func printGreen(text string) string {
	return fmt.Sprintf("\033[92m%s\033[0m", text)
}

func MergeLdapUsers(result LDAPUserSeachResult) {
	userObject := model.User{
		Username:  result.Username,
		Email:     result.Mail,
		FirstName: result.GivenName,
		LastName:  result.LastName,
		IsAdmin:   false,
	}
	fmt.Printf("First Name: %s\n", printGreen(userObject.FirstName))
	fmt.Printf("Last Name: %s\n", printGreen(userObject.LastName))
	fmt.Printf("Email: %s\n", printGreen(userObject.Email))
	fmt.Printf("Username: %s\n", printGreen(userObject.Username))
	fmt.Println("Admin Status:", userObject.IsAdmin)
}

func GetLdapUsers() LDAPUserSeachResult {
	client := ldapInit()
	// user := username
	baseDN := common.EnvConfig.LDAPSearchBase
	filter := "(objectClass=person)"

	//TODO Make options in UI to configure what options should be synced etc etc, as this depends on what ldap backend is being used.
	searchAttrs := []string{
		"sAMAccountName",
		"mail",
		"memberOf",
		"userPrincipalName",
		"givenName",
		"sn",
		"cn",
		"uid",
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
			MergeLdapUsers(userResult)
			// This temp username is just for my testing, until we can build out a full Web Config UI for this.
			// tempUsername := ""
			// if userResult.Username == "" {
			// 	tempUsername = userResult.UID
			// }
			// // fmt.Println("\nUser Attributes:")
			// // fmt.Printf("Full Name: %s\n", printGreen(userResult.GivenName+" "+userResult.LastName))
			// // fmt.Printf("Email: %s\n", printGreen(userResult.Mail))
			// // fmt.Printf("Username: %s\n", printGreen(tempUsername))
			// // fmt.Printf("DN: %s\n", printGreen(userResult.DN))
		}

		return userResult

	} else {
		fmt.Println("No Users Found")
		panic(1)
	}
}
