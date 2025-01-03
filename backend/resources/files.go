package resources

import "embed"

// Embedded file systems for the project

//go:embed email-templates images migrations
var FS embed.FS
