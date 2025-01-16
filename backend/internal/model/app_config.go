package model

type AppConfigVariable struct {
	Key          string `gorm:"primaryKey;not null"`
	Type         string
	IsPublic     bool
	IsInternal   bool
	Value        string
	DefaultValue string
}

type AppConfig struct {
	// General
	AppName             AppConfigVariable
	SessionDuration     AppConfigVariable
	EmailsVerified      AppConfigVariable
	AllowOwnAccountEdit AppConfigVariable
	// Internal
	BackgroundImageType AppConfigVariable
	LogoLightImageType  AppConfigVariable
	LogoDarkImageType   AppConfigVariable
	// Email
	EmailEnabled       AppConfigVariable
	SmtpHost           AppConfigVariable
	SmtpPort           AppConfigVariable
	SmtpFrom           AppConfigVariable
	SmtpUser           AppConfigVariable
	SmtpPassword       AppConfigVariable
	SmtpTls            AppConfigVariable
	SmtpSkipCertVerify AppConfigVariable
	// LDAP
	LdapEnabled                        AppConfigVariable
	LdapUrl                            AppConfigVariable
	LdapBindDn                         AppConfigVariable
	LdapBindPassword                   AppConfigVariable
	LdapBase                           AppConfigVariable
	LdapSkipCertVerify                 AppConfigVariable
	LdapAttributeUserUniqueIdentifier  AppConfigVariable
	LdapAttributeUserUsername          AppConfigVariable
	LdapAttributeUserEmail             AppConfigVariable
	LdapAttributeUserFirstName         AppConfigVariable
	LdapAttributeUserLastName          AppConfigVariable
	LdapAttributeGroupUniqueIdentifier AppConfigVariable
	LdapAttributeGroupName             AppConfigVariable
}
