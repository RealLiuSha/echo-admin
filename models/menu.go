package models

import (
	"strings"

	"github.com/RealLiuSha/echo-admin/models/database"
	"github.com/RealLiuSha/echo-admin/models/dto"
)

// ShowStatus - 1: show; -1: hide;
// Status - 1: Enable -1: Disable
type Menu struct {
	database.Model
	ID         string      `gorm:"column:id;size:36;not null;index;" json:"id"`
	Name       string      `gorm:"column:name;not null;index;" json:"name" validate:"required"`
	Sequence   int         `gorm:"column:sequence;not null;index;" json:"sequence" validate:"required"`
	Icon       string      `gorm:"column:icon;" json:"icon" validate:"required"`
	Router     string      `gorm:"column:router;" json:"router"`
	Component  string      `gorm:"column:component;" json:"component"`
	ParentID   string      `gorm:"column:parent_id;size:36;index;" json:"parent_id"`
	ParentPath string      `gorm:"column:parent_path;" json:"parent_path"`
	Hidden     int         `gorm:"column:hidden;not null;" json:"hidden" validate:"required,max=1,min=-1"`
	Status     int         `gorm:"column:status;not null;" json:"status" validate:"required,max=1,min=-1"`
	Remark     string      `gorm:"column:remark;" json:"remark" validate:"required"`
	CreatedBy  string      `gorm:"column:created_by;not null;" json:"created_by"`
	Actions    MenuActions `gorm:"-" json:"actions,omitempty"`
}

type MenuTree struct {
	ID         string      `yaml:"-" json:"id"`
	Name       string      `yaml:"name" json:"name"`
	Icon       string      `yaml:"icon" json:"icon"`
	Router     string      `yaml:"router,omitempty" json:"router"`
	Component  string      `yaml:"component,omitempty" json:"component"`
	ParentID   string      `yaml:"-" json:"parent_id"`
	ParentPath string      `yaml:"-" json:"parent_path"`
	Sequence   int         `yaml:"sequence" json:"sequence"`
	Hidden     int         `yaml:"-" json:"hidden"`
	Status     int         `yaml:"-" json:"status"`
	Actions    MenuActions `yaml:"actions,omitempty" json:"actions"`
	Children   MenuTrees   `yaml:"children,omitempty" json:"children,omitempty"`
}

type Menus []*Menu
type MenuTrees []*MenuTree

type MenuQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs              []string `query:"ids"`
	Name             string   `query:"name"`
	PrefixParentPath string   `query:"prefix_parent_path"`
	QueryValue       string   `query:"query_value"`
	ParentID         string   `query:"parent_id"`
	Hidden           int      `query:"hidden" validate:"max=1,min=-1"`
	Status           int      `query:"status" validate:"max=1,min=-1"`
	Tree             bool     `query:"tree"`
	IncludeActions   bool     `query:"include_actions"`
}

type MenuQueryResult struct {
	List       Menus           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	return a[i].Sequence > a[j].Sequence
}

func (a Menus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Menus) ToMap() map[string]*Menu {
	m := make(map[string]*Menu)

	for _, menu := range a {
		m[menu.ID] = menu
	}

	return m
}

func (a Menus) SplitParentIDs() []string {
	idList := make([]string, 0, len(a))
	mIDList := make(map[string]struct{})

	for _, item := range a {
		if _, ok := mIDList[item.ID]; ok || item.ParentPath == "" {
			continue
		}

		for _, pp := range strings.Split(item.ParentPath, "/") {
			if _, ok := mIDList[pp]; ok {
				continue
			}

			idList = append(idList, pp)
			mIDList[pp] = struct{}{}
		}
	}

	return idList
}

func (a Menus) ToIDs() []string {
	ids := make([]string, len(a))
	for i, item := range a {
		ids[i] = item.ID
	}
	return ids
}

func (a Menus) ToMenuTrees() MenuTrees {
	menuTrees := make(MenuTrees, len(a))
	for i, menu := range a {
		menuTrees[i] = &MenuTree{
			ID:         menu.ID,
			Name:       menu.Name,
			Icon:       menu.Icon,
			Router:     menu.Router,
			Component:  menu.Component,
			ParentID:   menu.ParentID,
			ParentPath: menu.ParentPath,
			Sequence:   menu.Sequence,
			Hidden:     menu.Hidden,
			Status:     menu.Status,
			Actions:    menu.Actions,
		}
	}

	return menuTrees.ToTree()
}

func (a MenuTrees) ToTree() MenuTrees {
	// tree map
	menuTreeMap := make(map[string]*MenuTree)
	for _, menuTree := range a {
		menuTreeMap[menuTree.ID] = menuTree
	}

	menuTrees := make(MenuTrees, 0)
	for _, menuTree := range a {
		if menuTree.ParentID == "" {
			menuTrees = append(menuTrees, menuTree)
			continue
		}

		if parentMenuTree, ok := menuTreeMap[menuTree.ParentID]; ok {
			if parentMenuTree.Children == nil {
				children := MenuTrees{menuTree}
				parentMenuTree.Children = children
				continue
			}

			parentMenuTree.Children = append(parentMenuTree.Children, menuTree)
		}
	}

	return menuTrees
}

func (a Menus) FillMenuAction(mActions map[string]MenuActions, mResources map[string]MenuActionResources) Menus {
	for _, item := range a {
		if v, ok := mActions[item.ID]; ok {
			item.Actions = v
			item.Actions.FillResources(mResources)
		}
	}
	return a
}
