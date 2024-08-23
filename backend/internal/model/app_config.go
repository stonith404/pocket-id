package model

type AppConfigVariable struct {
	Key        string `gorm:"primaryKey;not null"`
	Type       string
	IsPublic   bool
	IsInternal bool
	Value      string
}

type AppConfig struct {
	AppName             AppConfigVariable
	BackgroundImageType AppConfigVariable
	LogoImageType       AppConfigVariable
	SessionDuration     AppConfigVariable
}
