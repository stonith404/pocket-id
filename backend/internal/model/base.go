package model

import (
	"github.com/google/uuid"
	model "github.com/stonith404/pocket-id/backend/internal/model/types"
	"gorm.io/gorm"
	"time"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string `gorm:"primaryKey;not null"`
	CreatedAt model.DateTime
}

func (b *Base) BeforeCreate(_ *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	b.CreatedAt = model.DateTime(time.Now())
	return
}
