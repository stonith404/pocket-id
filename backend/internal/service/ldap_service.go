package service

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"gorm.io/gorm"
)

type LdapService struct {
	db           *gorm.DB
	userService  *UserService
	groupService *UserGroupService
}

func NewLdapService(db *gorm.DB, userService *UserService, groupService *UserGroupService) *LdapService {
	return &LdapService{db: db, userService: userService, groupService: groupService}
}

func ldapInit() *ldap.Conn {
	// Setup AD Connection
	ldapURL := common.EnvConfig.LDAPServer + ":" + common.EnvConfig.LDAPPort
	client, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: common.EnvConfig.LDAPTLSVerify}))
	if err != nil {
		//TODO Handle Errors Better
		panic(err)
	}

	// Bind as Service Account
	err = client.Bind(common.EnvConfig.LDAPBindUser, common.EnvConfig.LDAPBindPassword)
	if err != nil {
		//TODO Handle Errors Better
		panic(err)
	}
	return client
}

func (s *LdapService) GetLdapGroups() error {
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

	var groupError error

	if len(result.Entries) >= 1 {

		for _, value := range result.Entries {

			syncGroup := dto.UserGroupCreateDto{
				Name:         value.GetAttributeValue(common.EnvConfig.LDAPGroupAttribute),
				FriendlyName: value.GetAttributeValue(common.EnvConfig.LDAPGroupAttribute),
			}

			_, groupError = s.groupService.Create(syncGroup)

		}

		client.Close()
		return groupError
	} else {
		fmt.Println("No Groups Found")
		return groupError
	}

}

func (s *LdapService) GetLdapUsers() error {
	client := ldapInit()
	baseDN := common.EnvConfig.LDAPSearchBase
	filter := "(objectClass=person)"

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

	var userError error

	if len(result.Entries) >= 1 {

		for _, value := range result.Entries {

			newUser := dto.UserCreateDto{
				Username:  value.GetAttributeValue(common.EnvConfig.LDAPUsernameAttribute),
				Email:     value.GetAttributeValue("mail"),
				FirstName: value.GetAttributeValue("givenName"),
				LastName:  value.GetAttributeValue("sn"),
				IsAdmin:   false,
			}

			_, userError = s.userService.CreateUser(newUser)
		}

		client.Close()
		return userError

	} else {
		fmt.Println("No Users Found")
		return userError
	}

}
