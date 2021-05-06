package models

import (
	"github.com/RealLiuSha/echo-admin/models/database"
	"github.com/RealLiuSha/echo-admin/models/dto"
)

// Status - 1: Enable -1: Disable
type Role struct {
	database.Model
	ID        string    `gorm:"column:id;size:36;not null;index;" json:"id"`
	Name      string    `gorm:"column:name;not null;" json:"name" validate:"required"`
	Remark    string    `gorm:"column:remark;not null;" json:"remark" validate:"required"`
	Sequence  int       `gorm:"column:sequence;index;not null;" json:"sequence" validate:"required"`
	Status    int       `gorm:"column:status;default:0;not null;" json:"status" validate:"required,max=1,min=-1"`
	CreatedBy string    `gorm:"column:created_by;not null;" json:"created_by"`
	RoleMenus RoleMenus `gorm:"-" json:"role_menus"`
}

type Roles []*Role

type RoleQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	Name       string   `query:"name"`
	QueryValue string   `query:"query_value"`
	UserID     string   `query:"user_id"`
	Status     int      `query:"status" validate:"max=1,min=-1"`
}

type RoleQueryResult struct {
	List       Roles           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a Roles) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}

	return names
}

func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
