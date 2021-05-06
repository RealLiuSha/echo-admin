package models

import (
	"github.com/RealLiuSha/echo-admin/models/database"
	"github.com/RealLiuSha/echo-admin/models/dto"
)

type MenuActionResource struct {
	database.Model
	ID       string `gorm:"column:id;size:36;index;not null;" json:"-" yaml:"-"`
	ActionID string `gorm:"column:action_id;size:36;index;not null;" json:"-" yaml:"-"`
	Method   string `gorm:"column:method;not null;" json:"method" validate:"required" yaml:"method"`
	Path     string `gorm:"column:path;not null;" json:"path" validate:"required" yaml:"path"`
}

type MenuActionResources []*MenuActionResource

type MenuActionResourceQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	MenuID  string
	MenuIDs []string
}

type MenuActionResourceQueryResult struct {
	List       MenuActionResources `json:"list"`
	Pagination *dto.Pagination     `json:"pagination"`
}

func (a MenuActionResources) ToMap() map[string]*MenuActionResource {
	m := make(map[string]*MenuActionResource)
	for _, item := range a {
		m[item.Method+item.Path] = item
	}
	return m
}

func (a MenuActionResources) ToActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)
	for _, item := range a {
		m[item.ActionID] = append(m[item.ActionID], item)
	}

	return m
}
