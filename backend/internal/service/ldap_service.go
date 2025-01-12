package service

import "gorm.io/gorm"

type LdapService struct {
	db *gorm.DB
}

func NewLdapService(db *gorm.DB) *LdapService {
	return &LdapService{db: db}
}
