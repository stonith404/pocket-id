package model

type UserGroup struct {
	Base
	FriendlyName string `sortable:"true"`
	Name         string `gorm:"unique" sortable:"true"`
	Users        []User `gorm:"many2many:user_groups_users;"`
	CustomClaims []CustomClaim
}
