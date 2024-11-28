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
	AppName             AppConfigVariable
	SessionDuration     AppConfigVariable
	EmailsVerified      AppConfigVariable
	AllowOwnAccountEdit AppConfigVariable

	BackgroundImageType AppConfigVariable
	LogoLightImageType  AppConfigVariable
	LogoDarkImageType   AppConfigVariable

	EmailEnabled       AppConfigVariable
	SmtpHost           AppConfigVariable
	SmtpPort           AppConfigVariable
	SmtpFrom           AppConfigVariable
	SmtpUser           AppConfigVariable
	SmtpPassword       AppConfigVariable
	SmtpTls            AppConfigVariable
	SmtpSkipCertVerify AppConfigVariable
}
