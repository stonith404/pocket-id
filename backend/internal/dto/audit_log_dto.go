package dto

import (
	"github.com/stonith404/pocket-id/backend/internal/model"
	datatype "github.com/stonith404/pocket-id/backend/internal/model/types"
)

type AuditLogDto struct {
	ID        string            `json:"id"`
	CreatedAt datatype.DateTime `json:"createdAt"`

	Event     model.AuditLogEvent `json:"event"`
	IpAddress string              `json:"ipAddress"`
	Country   string              `json:"country"`
	City      string              `json:"city"`
	ISP       string              `json:"isp"`
	ASNumber  uint                `json:"asNumber"`
	Device    string              `json:"device"`
	UserID    string              `json:"userID"`
	Data      model.AuditLogData  `json:"data"`
}
