package models

import (
	"github.com/RealLiuSha/echo-admin/models/database"
	"github.com/RealLiuSha/echo-admin/models/dto"
)

type RoleMenu struct {
	database.Model
	ID       string `gorm:"column:id;size:36;not null;" json:"id"`
	RoleID   string `gorm:"column:role_id;size:36;not null;index;" json:"role_id" validate:"required"`
	MenuID   string `gorm:"column:menu_id;size:36;not null;index;" json:"menu_id" validate:"required"`
	ActionID string `gorm:"column:action_id;size:36;not null;index;" json:"action_id" validate:"required"`
}

type RoleMenus []*RoleMenu

type RoleMenuQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	RoleID  string
	RoleIDs []string
}

type RoleMenuQueryResult struct {
	List       RoleMenus       `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a RoleMenus) ToMap() map[string]*RoleMenu {
	m := make(map[string]*RoleMenu)
	for _, item := range a {
		m[item.MenuID+"-"+item.ActionID] = item
	}

	return m
}

func (a RoleMenus) ToRoleIDMap() map[string]RoleMenus {
	m := make(map[string]RoleMenus)
	for _, item := range a {
		m[item.RoleID] = append(m[item.RoleID], item)
	}

	return m
}

func (a RoleMenus) ToMenuIDs() []string {
	var idList []string
	m := make(map[string]struct{})

	for _, item := range a {
		if _, ok := m[item.MenuID]; ok {
			continue
		}
		idList = append(idList, item.MenuID)
		m[item.MenuID] = struct{}{}
	}

	return idList
}

func (a RoleMenus) ToActionIDs() []string {
	idList := make([]string, len(a))

	m := make(map[string]struct{})
	for i, item := range a {
		if _, ok := m[item.ActionID]; ok {
			continue
		}
		idList[i] = item.ActionID
		m[item.ActionID] = struct{}{}
	}

	return idList
}
