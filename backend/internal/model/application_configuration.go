package model

type ApplicationConfigurationVariable struct {
	Key        string `gorm:"primaryKey;not null" json:"key"`
	Type       string `json:"type"`
	IsPublic   bool   `json:"-"`
	IsInternal bool   `json:"-"`
	Value      string `json:"value"`
}

type ApplicationConfiguration struct {
	AppName             ApplicationConfigurationVariable
	BackgroundImageType ApplicationConfigurationVariable
	LogoImageType       ApplicationConfigurationVariable
}

type ApplicationConfigurationUpdateDto struct {
	AppName string `json:"appName" binding:"required"`
}
