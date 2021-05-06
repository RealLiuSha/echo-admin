package models

import (
	"github.com/RealLiuSha/echo-admin/models/database"
	"github.com/RealLiuSha/echo-admin/models/dto"
)

type UserRole struct {
	database.Model
	ID     string `gorm:"column:id;size:36;not null;" json:"id"`
	UserID string `gorm:"column:user_id;size:36;index;not null;" json:"user_id"`
	RoleID string `gorm:"column:role_id;size:36;index;not null;" json:"role_id"`
}

type UserRoles []*UserRole

type UserRoleQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	UserID  string
	UserIDs []string
}

type UserRoleQueryResult struct {
	List       UserRoles       `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a UserRoles) ToMap() map[string]*UserRole {
	m := make(map[string]*UserRole)
	for _, item := range a {
		m[item.RoleID] = item
	}

	return m
}

func (a UserRoles) ToRoleIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}

	return list
}

func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, item := range a {
		m[item.UserID] = append(m[item.UserID], item)
	}

	return m
}
