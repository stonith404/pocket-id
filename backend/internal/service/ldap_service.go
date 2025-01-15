package service

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"github.com/stonith404/pocket-id/backend/internal/model"
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
	client, err := ldap.DialURL(common.EnvConfig.LDAPUrl, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: common.EnvConfig.LDAPTLSVerify}))
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

	for _, value := range result.Entries {

		var databaseGroup model.UserGroup
		groupUniqueName := value.GetAttributeValue(common.EnvConfig.LDAPGroupAttribute)
		s.db.Where("name = ?", groupUniqueName).First(&databaseGroup)

		syncGroup := dto.UserGroupCreateDto{
			Name:         value.GetAttributeValue(common.EnvConfig.LDAPGroupAttribute),
			FriendlyName: value.GetAttributeValue(common.EnvConfig.LDAPGroupAttribute),
		}

		if databaseGroup.ID == "" {
			_, err = s.groupService.Create(syncGroup)
			if err != nil {
				log.Printf("Error syncing group %s: %s", syncGroup.Name, err)
			}
		} else {
			_, err := s.groupService.Update(databaseGroup.ID, syncGroup)
			if err != nil {
				log.Printf("Error syncing group %s: %s", syncGroup.Name, err)
			}

		}

	}

	client.Close()
	return nil

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

	for _, value := range result.Entries {
		ldapId := value.GetAttributeValue(common.EnvConfig.LDAPUserIdAttribute)

		// Get the user from the database
		var databaseUser model.User
		s.db.Where("ldap_id = ?", ldapId).First(&databaseUser)

		newUser := dto.UserCreateDto{
			Username:  value.GetAttributeValue(common.EnvConfig.LDAPUsernameAttribute),
			Email:     value.GetAttributeValue("mail"),
			FirstName: value.GetAttributeValue("givenName"),
			LastName:  value.GetAttributeValue("sn"),
			IsAdmin:   false,
			LdapID:    ldapId,
		}

		if databaseUser.ID == "" {
			_, err = s.userService.CreateUser(newUser)
			if err != nil {
				log.Printf("Error syncing user %s: %s", newUser.Username, err)
			}
		} else {
			_, err = s.userService.UpdateUser(databaseUser.ID, newUser, false)
			if err != nil {
				log.Printf("Error syncing user %s: %s", newUser.Username, err)
			}

		}

	}

	client.Close()
	return nil

}
