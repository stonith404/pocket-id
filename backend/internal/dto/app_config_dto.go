package dto

type PublicAppConfigVariableDto struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AppConfigVariableDto struct {
	PublicAppConfigVariableDto
	IsPublic bool `json:"isPublic"`
}

type AppConfigUpdateDto struct {
	AppName         string `json:"appName" binding:"required,min=1,max=30"`
	SessionDuration string `json:"sessionDuration" binding:"required"`
	EmailEnabled    string `json:"emailEnabled" binding:"required"`
	SmtHost         string `json:"smtpHost" binding:"required"`
	SmtpPort        string `json:"smtpPort" binding:"required"`
	SmtpFrom        string `json:"smtpFrom" binding:"email"`
	SmtpUser        string `json:"smtpUser" binding:"required"`
	SmtpPassword    string `json:"smtpPassword" binding:"required"`
}
