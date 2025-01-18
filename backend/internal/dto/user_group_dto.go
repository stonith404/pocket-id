package dto

import (
	datatype "github.com/stonith404/pocket-id/backend/internal/model/types"
)

type UserGroupDtoWithUsers struct {
	ID           string            `json:"id"`
	FriendlyName string            `json:"friendlyName"`
	Name         string            `json:"name"`
	CustomClaims []CustomClaimDto  `json:"customClaims"`
	Users        []UserDto         `json:"users"`
	CreatedAt    datatype.DateTime `json:"createdAt"`
}

type UserGroupDtoWithUserCount struct {
	ID           string            `json:"id"`
	FriendlyName string            `json:"friendlyName"`
	Name         string            `json:"name"`
	CustomClaims []CustomClaimDto  `json:"customClaims"`
	UserCount    int64             `json:"userCount"`
	CreatedAt    datatype.DateTime `json:"createdAt"`
}

type UserGroupCreateDto struct {
	FriendlyName string `json:"friendlyName" binding:"required,min=2,max=50"`
	Name         string `json:"name" binding:"required,min=2,max=255"`
}

type UserGroupUpdateUsersDto struct {
	UserIDs []string `json:"userIds" binding:"required"`
}

type AssignUserToGroupDto struct {
	UserID string `json:"userId" binding:"required"`
}
