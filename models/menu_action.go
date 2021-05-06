package models

import (
	"github.com/RealLiuSha/echo-admin/models/database"
	"github.com/RealLiuSha/echo-admin/models/dto"
)

type MenuAction struct {
	database.Model
	ID        string              `gorm:"column:id;size:36;not null;index;" json:"id" yaml:"-"`
	MenuID    string              `gorm:"column:menu_id;size:36;not null;index;" json:"menu_id" yaml:"-"`
	Code      string              `gorm:"column:code;not null;" json:"code" validate:"required" yaml:"code"`
	Name      string              `gorm:"column:name;not null;" json:"name" validate:"required" yaml:"name"`
	Resources MenuActionResources `gorm:"-" json:"resources" yaml:"resources"`
}

type MenuActions []*MenuAction

type MenuActionQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	MenuID string
	IDs    []string
}

type MenuActionQueryResult struct {
	List       MenuActions     `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

func (a MenuActions) FillResources(maResources map[string]MenuActionResources) {
	for i, item := range a {
		a[i].Resources = maResources[item.ID]
	}
}

func (a MenuActions) ToMenuIDMap() map[string]MenuActions {
	m := make(map[string]MenuActions)
	for _, item := range a {
		m[item.MenuID] = append(m[item.MenuID], item)
	}

	return m
}
