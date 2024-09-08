package dto

import (
	"github.com/stonith404/pocket-id/backend/internal/model"
	"time"
)

type AuditLogDto struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	Event     model.AuditLogEvent `json:"event"`
	IpAddress string              `json:"ipAddress"`
	Device    string              `json:"device"`
	UserID    string              `json:"userID"`
	Data      model.AuditLogData  `json:"data"`
}
