package model

type AppConfigVariable struct {
	Key        string `gorm:"primaryKey;not null" json:"key"`
	Type       string `json:"type"`
	IsPublic   bool   `json:"-"`
	IsInternal bool   `json:"-"`
	Value      string `json:"value"`
}

type AppConfig struct {
	AppName             AppConfigVariable
	BackgroundImageType AppConfigVariable
	LogoImageType       AppConfigVariable
	SessionDuration     AppConfigVariable
}

type AppConfigUpdateDto struct {
	AppName string `json:"appName" binding:"required"`
}
