package services

import (
	"gorm.io/gorm"

	"github.com/RealLiuSha/echo-admin/api/repository"
	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
	"github.com/RealLiuSha/echo-admin/models/dto"
	"github.com/RealLiuSha/echo-admin/pkg/uuid"
)

// MenuService service layer
type MenuService struct {
	logger                       lib.Logger
	menuRepository               repository.MenuRepository
	menuActionRepository         repository.MenuActionRepository
	menuActionResourceRepository repository.MenuActionResourceRepository
}

// NewMenuService creates a new menu service
func NewMenuService(
	logger lib.Logger,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
	menuActionResourceRepository repository.MenuActionResourceRepository,
) MenuService {
	return MenuService{
		logger:                       logger,
		menuRepository:               menuRepository,
		menuActionRepository:         menuActionRepository,
		menuActionResourceRepository: menuActionResourceRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a MenuService) WithTrx(trxHandle *gorm.DB) MenuService {
	a.menuRepository = a.menuRepository.WithTrx(trxHandle)
	a.menuActionRepository = a.menuActionRepository.WithTrx(trxHandle)
	a.menuActionResourceRepository = a.menuActionResourceRepository.WithTrx(trxHandle)

	return a
}

func (a MenuService) Check(item *models.Menu) error {
	result, err := a.menuRepository.Query(&models.MenuQueryParam{
		Name:     item.Name,
		ParentID: item.ParentID,
	})

	if err != nil {
		return err
	} else if len(result.List) > 0 {
		return errors.MenuAlreadyExists
	}

	return nil
}

func (a MenuService) Query(param *models.MenuQueryParam) (*models.MenuQueryResult, error) {
	menuQR, err := a.menuRepository.Query(param)
	if err != nil {
		return nil, err
	}

	if !param.IncludeActions {
		return menuQR, nil
	}

	menuActionQR, err := a.menuActionRepository.Query(&models.MenuActionQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return nil, err
	}

	menuResourceQR, err := a.menuActionResourceRepository.Query(&models.MenuActionResourceQueryParam{
		MenuIDs: menuQR.List.ToIDs(), PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return nil, err
	}

	menuQR.List.FillMenuAction(menuActionQR.List.ToMenuIDMap(), menuResourceQR.List.ToActionIDMap())
	return menuQR, nil
}

func (a MenuService) GetMenuActions(id string) (models.MenuActions, error) {
	paginationParam := dto.PaginationParam{PageSize: 999, Current: 1}

	menuActionQR, err := a.menuActionRepository.Query(&models.MenuActionQueryParam{
		MenuID: id, PaginationParam: paginationParam,
	})

	if err != nil {
		return nil, err
	} else if len(menuActionQR.List) == 0 {
		return nil, nil
	}

	menuResourceQR, err := a.menuActionResourceRepository.Query(&models.MenuActionResourceQueryParam{
		MenuID: id, PaginationParam: paginationParam,
	})

	if err != nil {
		return nil, err
	}

	menuActionQR.List.FillResources(menuResourceQR.List.ToActionIDMap())
	return menuActionQR.List, nil
}

func (a MenuService) Get(id string) (*models.Menu, error) {
	menu, err := a.menuRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (a MenuService) Create(menu *models.Menu) (id string, err error) {
	if err = a.Check(menu); err != nil {
		return
	}

	if menu.ParentPath, err = a.GetParentPath(menu.ParentID); err != nil {
		return
	}

	menu.ID = uuid.MustString()
	if err = a.menuRepository.Create(menu); err != nil {
		return
	}

	return menu.ID, nil
}

func (a MenuService) CreateMenus(parentID string, mTrees models.MenuTrees) error {
	for _, mTree := range mTrees {
		menu := &models.Menu{
			Name:      mTree.Name,
			Sequence:  mTree.Sequence,
			Icon:      mTree.Icon,
			Router:    mTree.Router,
			Component: mTree.Component,
			ParentID:  parentID,
			Status:    1,
			Hidden:    -1,
		}

		if v := mTree.Hidden; v != 0 {
			menu.Hidden = v
		}

		menuID, err := a.Create(menu)
		if err != nil {
			return err
		}

		if err := a.CreateActions(menu.ID, mTree.Actions); err != nil {
			return err
		}

		if mTree.Children != nil && len(mTree.Children) > 0 {
			err := a.CreateMenus(menuID, mTree.Children)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a MenuService) CreateActions(menuID string, menuActions models.MenuActions) error {
	for _, menuAction := range menuActions {
		menuAction.ID = uuid.MustString()
		menuAction.MenuID = menuID

		if err := a.menuActionRepository.Create(menuAction); err != nil {
			return err
		}

		for _, resource := range menuAction.Resources {
			resource.ID = uuid.MustString()
			resource.ActionID = menuAction.ID

			if err := a.menuActionResourceRepository.Create(resource); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a MenuService) Update(id string, menu *models.Menu) error {
	if id == menu.ParentID {
		return errors.MenuInvalidParent
	}

	// get old menu
	oMenu, err := a.Get(id)
	if err != nil {
		return err
	} else if oMenu.Name != menu.Name {
		if err = a.Check(menu); err != nil {
			return err
		}
	}

	menu.ID = oMenu.ID
	menu.CreatedBy = oMenu.CreatedBy
	menu.CreatedAt = oMenu.CreatedAt

	if menu.ParentID != oMenu.ParentID {
		parentPath, err := a.GetParentPath(menu.ParentID)
		if err != nil {
			return err
		}

		menu.ParentPath = parentPath
	} else {
		menu.ParentPath = oMenu.ParentPath
	}

	if err = a.UpdateChildParentPath(oMenu, menu); err != nil {
		return err
	}

	if err = a.menuRepository.Update(id, menu); err != nil {
		return err
	}

	return nil
}

func (a MenuService) UpdateActions(menuID string, actions models.MenuActions) error {
	oActions, err := a.GetMenuActions(menuID)
	if err != nil {
		return err
	}

	aActions, dActions, uActions := a.CompareActions(oActions, actions)

	err = a.CreateActions(menuID, aActions)
	if err != nil {
		return err
	}

	for _, dAction := range dActions {
		if err = a.menuActionRepository.Delete(dAction.ID); err != nil {
			return err
		}

		if err = a.menuActionResourceRepository.DeleteByActionID(dAction.ID); err != nil {
			return err
		}
	}

	oMap := oActions.ToMap()
	for _, uAction := range uActions {
		// old menu action
		oAction := oMap[uAction.Code]

		// update action name
		if uAction.Name != oAction.Name {
			oAction.Name = uAction.Name
			if err = a.menuActionRepository.Update(uAction.ID, oAction); err != nil {
				return err
			}
		}

		// compare resources to update
		aResources, dResources := a.CompareResources(oAction.Resources, uAction.Resources)
		for _, aResource := range aResources {
			aResource.ID = uuid.MustString()
			aResource.ActionID = oAction.ID

			err := a.menuActionResourceRepository.Create(aResource)
			if err != nil {
				return err
			}
		}

		for _, dResource := range dResources {
			err := a.menuActionResourceRepository.Delete(dResource.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a MenuService) Delete(id string) error {
	_, err := a.menuRepository.Get(id)
	if err != nil {
		return err
	}

	menuQR, err := a.menuRepository.Query(&models.MenuQueryParam{
		ParentID: id,
	})

	if err != nil {
		return err
	} else if menuQR.Pagination.Total > 0 {
		return errors.MenuNotAllowDeleteWithChild
	}

	if err = a.menuActionResourceRepository.DeleteByMenuID(id); err != nil {
		return err
	}

	if err = a.menuActionRepository.DeleteByMenuID(id); err != nil {
		return err
	}

	if err = a.menuRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a MenuService) UpdateStatus(id string, status int) error {
	_, err := a.menuRepository.Get(id)
	if err != nil {
		return err
	}

	return a.menuRepository.UpdateStatus(id, status)
}

func (a MenuService) GetParentPath(parentID string) (string, error) {
	if parentID == "" {
		return "", nil
	}

	parentMenu, err := a.menuRepository.Get(parentID)
	if err != nil {
		return "", err
	}

	return a.JoinParentPath(parentMenu.ParentPath, parentMenu.ID), nil
}

func (a MenuService) JoinParentPath(parent, id string) string {
	if parent != "" {
		return parent + "/" + id
	}

	return id
}

func (a MenuService) CompareActions(oActions, nActions models.MenuActions) (aList, dList, uList models.MenuActions) {
	oMap := oActions.ToMap()
	nMap := nActions.ToMap()

	for k, item := range nMap {
		if _, ok := oMap[k]; ok {
			uList = append(uList, item)
			delete(oMap, k)

			continue
		}

		aList = append(aList, item)
	}

	for _, item := range oMap {
		dList = append(dList, item)
	}

	return
}

func (a MenuService) CompareResources(oResources, nResources models.MenuActionResources) (aList, dList models.MenuActionResources) {
	oMap := oResources.ToMap()
	nMap := nResources.ToMap()

	for k, item := range nMap {
		if _, ok := oMap[k]; ok {
			delete(oMap, k)
			continue
		}

		aList = append(aList, item)
	}

	for _, item := range oMap {
		dList = append(dList, item)
	}

	return
}

func (a MenuService) UpdateChildParentPath(oMenu, nMenu *models.Menu) error {
	if oMenu.ParentID == nMenu.ParentID {
		return nil
	}

	oPath := a.JoinParentPath(oMenu.ParentPath, oMenu.ID)
	menuQR, err := a.menuRepository.Query(&models.MenuQueryParam{
		PrefixParentPath: oPath,
	})

	if err != nil {
		return err
	}

	nPath := a.JoinParentPath(nMenu.ParentPath, nMenu.ID)
	for _, menu := range menuQR.List {
		err = a.menuRepository.UpdateParentPath(menu.ID, nPath+menu.ParentPath[len(oPath):])
		if err != nil {
			return err
		}
	}

	return nil
}
