package service

import (
	"fmt"
	"github.com/stonith404/pocket-id/backend/internal/utils/email"
	"time"
)

/**
How to add new template:
- pick unique and descriptive template ${name} (for example "login-with-new-device")
- in backend/email-templates/ create "${name}_html.tmpl" and "${name}_text.tmpl"
- create xxxxTemplate and xxxxTemplateData (for example NewLoginTemplate and NewLoginTemplateData)
  - Path *must* be ${name}
- add xxxTemplate.Path to "emailTemplatePaths" at the end

Notes:
- backend app must be restarted to reread all the template files
- root "." object in templates is `email.TemplateData`
- xxxxTemplateData structure is visible under .Data in templates
*/

var NewLoginTemplate = email.Template[NewLoginTemplateData]{
	Path: "login-with-new-device",
	Title: func(data *email.TemplateData[NewLoginTemplateData]) string {
		return fmt.Sprintf("New device login with %s", data.AppName)
	},
}

type NewLoginTemplateData struct {
	IPAddress string
	Device    string
	DateTime  time.Time
}

// this is list of all template paths used for preloading templates
var emailTemplatesPaths = []string{NewLoginTemplate.Path}
