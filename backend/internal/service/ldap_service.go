package service

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"gorm.io/gorm"
)

type LdapService struct {
	db               *gorm.DB
	appConfigService *AppConfigService
	userService      *UserService
	groupService     *UserGroupService
}

func NewLdapService(db *gorm.DB, appConfigService *AppConfigService, userService *UserService, groupService *UserGroupService) *LdapService {
	return &LdapService{db: db, appConfigService: appConfigService, userService: userService, groupService: groupService}
}

func (s *LdapService) createClient() (*ldap.Conn, error) {
	if s.appConfigService.DbConfig.LdapEnabled.Value != "true" {
		return nil, fmt.Errorf("LDAP is not enabled")
	}
	// Setup AD Connection
	ldapURL := s.appConfigService.DbConfig.LdapUrl.Value
	skipTLSVerify := s.appConfigService.DbConfig.LdapSkipCertVerify.Value == "true"
	client, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: skipTLSVerify}))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP: %w", err)
	}

	// Bind as Service Account
	bindDn := s.appConfigService.DbConfig.LdapBindDn.Value
	bindPassword := s.appConfigService.DbConfig.LdapBindPassword.Value
	err = client.Bind(bindDn, bindPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to LDAP: %w", err)
	}
	return client, nil
}

func (s *LdapService) SyncAll() error {
	err := s.SyncUsers()
	if err != nil {
		return fmt.Errorf("failed to sync users: %w", err)
	}

	err = s.SyncGroups()
	if err != nil {
		return fmt.Errorf("failed to sync groups: %w", err)
	}

	return nil
}

func (s *LdapService) SyncGroups() error {
	// Setup LDAP Connection
	client, err := s.createClient()
	if err != nil {
		return fmt.Errorf("failed to create LDAP client: %w", err)
	}
	defer client.Close()

	baseDN := s.appConfigService.DbConfig.LdapBase.Value
	nameAttribute := s.appConfigService.DbConfig.LdapAttributeGroupName.Value
	uniqueIdentifierAttribute := s.appConfigService.DbConfig.LdapAttributeGroupUniqueIdentifier.Value
	filter := "(objectClass=groupOfUniqueNames)"

	searchAttrs := []string{
		nameAttribute,
		uniqueIdentifierAttribute,
		"member",
	}

	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, searchAttrs, []ldap.Control{})
	result, err := client.Search(searchReq)
	if err != nil {
		return fmt.Errorf("failed to query LDAP: %w", err)
	}

	for _, value := range result.Entries {
		var usersToAddDto dto.UserGroupUpdateUsersDto
		var userIDStrings []string

		// Try to find the group in the database
		ldapId := value.GetAttributeValue(uniqueIdentifierAttribute)
		var databaseGroup model.UserGroup
		s.db.Where("ldap_id = ?", ldapId).First(&databaseGroup)

		//Get Group Members and add to the correct Group
		groupMembers := value.GetAttributeValues("member")
		for _, member := range groupMembers {
			// Normal output of this would be CN=username,ou=people,dc=example,dc=com
			// Splitting at the "=" and "," then just grabbing the username for that string
			singleMember := strings.Split(strings.Split(member, "=")[1], ",")[0]

			var databaseUser model.User
			s.db.Where("username = ?", singleMember).First(&databaseUser)
			userIDStrings = append(userIDStrings, databaseUser.ID)
		}

		syncGroup := dto.UserGroupCreateDto{
			Name:         value.GetAttributeValue(nameAttribute),
			FriendlyName: value.GetAttributeValue(nameAttribute),
			LdapID:       value.GetAttributeValue(uniqueIdentifierAttribute),
		}

		usersToAddDto = dto.UserGroupUpdateUsersDto{
			UserIDs: userIDStrings,
		}

		if databaseGroup.ID == "" {
			newGroup, err := s.groupService.Create(syncGroup)
			if err != nil {
				log.Printf("Error syncing group %s: %s", syncGroup.Name, err)
			} else {
				if _, err = s.groupService.UpdateUsers(newGroup.ID, usersToAddDto); err != nil {
					log.Printf("Error syncing group %s: %s", syncGroup.Name, err)
				}
			}
		} else {
			_, err = s.groupService.Update(databaseGroup.ID, syncGroup, true)
			_, err = s.groupService.UpdateUsers(databaseGroup.ID, usersToAddDto)
			if err != nil {
				log.Printf("Error syncing group %s: %s", syncGroup.Name, err)
				return err
			}

		}

	}

	return nil

}

func (s *LdapService) SyncUsers() error {
	// Setup LDAP Connection
	client, err := s.createClient()
	if err != nil {
		return fmt.Errorf("failed to create LDAP client: %w", err)
	}
	defer client.Close()

	baseDN := s.appConfigService.DbConfig.LdapBase.Value
	uniqueIdentifierAttribute := s.appConfigService.DbConfig.LdapAttributeUserUniqueIdentifier.Value
	usernameAttribute := s.appConfigService.DbConfig.LdapAttributeUserUsername.Value
	emailAttribute := s.appConfigService.DbConfig.LdapAttributeUserEmail.Value
	firstNameAttribute := s.appConfigService.DbConfig.LdapAttributeUserFirstName.Value
	lastNameAttribute := s.appConfigService.DbConfig.LdapAttributeUserLastName.Value

	filter := "(objectClass=person)"

	searchAttrs := []string{
		"memberOf",
		"sn",
		"cn",
		uniqueIdentifierAttribute,
		usernameAttribute,
		emailAttribute,
		firstNameAttribute,
		lastNameAttribute,
	}

	// Filters must start and finish with ()!
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, searchAttrs, []ldap.Control{})

	result, err := client.Search(searchReq)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to query LDAP: %w", err))
	}

	//Get all Current Database Users
	var databaseUsers []model.User
	if err := s.db.Find(&databaseUsers).Error; err != nil {
		fmt.Println(fmt.Errorf("Failed to Fetch Users from Database: %v", err))
	}

	//Create Mapping for Users that exsist
	ldapUsers := make(map[string]bool)
	missingUsers := []model.User{}

	for _, value := range result.Entries {
		ldapId := value.GetAttributeValue(uniqueIdentifierAttribute)

		//This Maps the Users to this array if they exsist
		ldapUsers[ldapId] = true

		// Get the user from the database
		var databaseUser model.User
		s.db.Where("ldap_id = ?", ldapId).First(&databaseUser)

		newUser := dto.UserCreateDto{
			Username:  value.GetAttributeValue(usernameAttribute),
			Email:     value.GetAttributeValue(emailAttribute),
			FirstName: value.GetAttributeValue(firstNameAttribute),
			LastName:  value.GetAttributeValue(lastNameAttribute),
			IsAdmin:   false,
			LdapID:    ldapId,
		}

		if databaseUser.ID == "" {
			_, err = s.userService.CreateUser(newUser)
			if err != nil {
				log.Printf("Error syncing user %s: %s", newUser.Username, err)
			}
		} else {
			_, err = s.userService.UpdateUser(databaseUser.ID, newUser, false, true)
			if err != nil {
				log.Printf("Error syncing user %s: %s", newUser.Username, err)
			}

		}

	}

	dbUserCount := len(databaseUsers) - 1 //Accounting for the built in Admin User
	//Compare Database Users with LDAP Users
	if dbUserCount > len(ldapUsers) {
		for _, dbUser := range databaseUsers {
			if dbUser.LdapID == nil {
				continue
			}
			if _, exists := ldapUsers[*dbUser.LdapID]; !exists {
				fmt.Printf("Ldap id: %s, username: %s\n", *dbUser.LdapID, dbUser.Username)
				missingUsers = append(missingUsers, dbUser)
			}
		}

		//Remove Users from Database if they no longer exsist in LDAP
		for _, missingUser := range missingUsers {
			if err := s.db.Delete(&missingUser).Error; err != nil {
				log.Printf("Failed to delete user %s: %v", missingUser.Username, err)
			} else {
				fmt.Printf("Removed missing user: %s\n", missingUser.Username)
			}
		}
	}

	return nil
}
