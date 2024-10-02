package dto

import "time"

type UserGroupDtoWithUsers struct {
	ID           string    `json:"id"`
	FriendlyName string    `json:"friendlyName"`
	Name         string    `json:"name"`
	Users        []UserDto `json:"users"`
	CreatedAt    time.Time `json:"createdAt"`
}

type UserGroupDtoWithUserCount struct {
	ID           string    `json:"id"`
	FriendlyName string    `json:"friendlyName"`
	Name         string    `json:"name"`
	UserCount    int64     `json:"userCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

type UserGroupCreateDto struct {
	FriendlyName string `json:"friendlyName" binding:"required,min=3,max=30"`
	Name         string `json:"name" binding:"required,min=3,max=30,userGroupName"`
}

type UserGroupUpdateUsersDto struct {
	UserIDs []string `json:"userIds" binding:"required"`
}

type AssignUserToGroupDto struct {
	UserID string `json:"userId" binding:"required"`
}
