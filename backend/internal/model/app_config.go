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
	BackgroundImageType AppConfigVariable
	LogoLightImageType  AppConfigVariable
	LogoDarkImageType   AppConfigVariable
	SessionDuration     AppConfigVariable
	EmailsVerified      AppConfigVariable

	EmailEnabled AppConfigVariable
	SmtpHost     AppConfigVariable
	SmtpPort     AppConfigVariable
	SmtpFrom     AppConfigVariable
	SmtpUser     AppConfigVariable
	SmtpPassword AppConfigVariable
}
