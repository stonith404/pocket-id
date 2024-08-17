package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string    `gorm:"primaryKey;not null" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func (b *Base) BeforeCreate(_ *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}
