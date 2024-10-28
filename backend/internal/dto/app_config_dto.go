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
	AppName             string `json:"appName" binding:"required,min=1,max=30"`
	SessionDuration     string `json:"sessionDuration" binding:"required"`
	EmailsVerified      string `json:"emailsVerified" binding:"required"`
	AllowOwnAccountEdit string `json:"allowOwnAccountEdit" binding:"required"`
	EmailEnabled        string `json:"emailEnabled" binding:"required"`
	SmtHost             string `json:"smtpHost"`
	SmtpPort            string `json:"smtpPort"`
	SmtpFrom            string `json:"smtpFrom" binding:"omitempty,email"`
	SmtpUser            string `json:"smtpUser"`
	SmtpPassword        string `json:"smtpPassword"`
}
