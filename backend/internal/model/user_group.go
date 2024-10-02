package model

type UserGroup struct {
	Base
	FriendlyName string
	Name         string `gorm:"unique"`
	Users        []User `gorm:"many2many:user_groups_users;"`
}
