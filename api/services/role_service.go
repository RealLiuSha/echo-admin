package services

import (
	"gorm.io/gorm"

	"github.com/RealLiuSha/echo-admin/api/repository"
	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
	"github.com/RealLiuSha/echo-admin/pkg/uuid"
)

// RoleService service layer
type RoleService struct {
	logger               lib.Logger
	casbinService        CasbinService
	userRepository       repository.UserRepository
	roleRepository       repository.RoleRepository
	roleMenuRepository   repository.RoleMenuRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
}

// NewRoleService creates a new roleservice
func NewRoleService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	roleMenuRepository repository.RoleMenuRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) RoleService {
	return RoleService{
		logger:               logger,
		casbinService:        casbinService,
		userRepository:       userRepository,
		roleRepository:       roleRepository,
		roleMenuRepository:   roleMenuRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a RoleService) WithTrx(trxHandle *gorm.DB) RoleService {
	a.roleRepository = a.roleRepository.WithTrx(trxHandle)
	a.userRepository = a.userRepository.WithTrx(trxHandle)
	a.roleMenuRepository = a.roleMenuRepository.WithTrx(trxHandle)

	return a
}

func (a RoleService) Query(param *models.RoleQueryParam) (roleQR *models.RoleQueryResult, err error) {
	return a.roleRepository.Query(param)
}

func (a RoleService) QueryRoleMenus(roleID string) (models.RoleMenus, error) {
	roleMenuQR, err := a.roleMenuRepository.Query(&models.RoleMenuQueryParam{
		RoleID: roleID,
	})

	if err != nil {
		return nil, err
	}

	return roleMenuQR.List, nil
}

func (a RoleService) Get(id string) (*models.Role, error) {
	role, err := a.roleRepository.Get(id)
	if err != nil {
		return nil, err
	}

	roleMenus, err := a.QueryRoleMenus(id)
	if err != nil {
		return nil, err
	}

	role.RoleMenus = roleMenus
	return role, nil
}

func (a RoleService) Check(item *models.Role) error {
	qr, err := a.roleRepository.Query(&models.RoleQueryParam{Name: item.Name})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.RoleAlreadyExists
	}

	return nil
}

func (a RoleService) CheckRoleMenu(rMenu *models.RoleMenu) error {
	if _, err := a.menuRepository.Get(rMenu.MenuID); err != nil {
		return errors.Wrap(err, "menu id")
	}

	if _, err := a.menuActionRepository.Get(rMenu.ActionID); err != nil {
		return errors.Wrap(err, "menu action id")
	}

	return nil
}

func (a RoleService) CompareRoleMenus(oRoleMenus, nRoleMenus models.RoleMenus) (aList, dList models.RoleMenus) {
	oMap := oRoleMenus.ToMap()
	nMap := nRoleMenus.ToMap()

	for k, nRoleMenu := range nMap {
		if _, ok := oMap[k]; ok {
			delete(oMap, k)
			continue
		}
		aList = append(aList, nRoleMenu)
	}

	for _, oRoleMenu := range oMap {
		dList = append(dList, oRoleMenu)
	}
	return
}

func (a RoleService) Create(role *models.Role) (id string, err error) {
	if err = a.Check(role); err != nil {
		return
	}

	role.ID = uuid.MustString()
	for _, roleMenu := range role.RoleMenus {
		roleMenu.ID = uuid.MustString()
		roleMenu.RoleID = role.ID

		if err = a.CheckRoleMenu(roleMenu); err != nil {
			return
		}

		if err = a.roleMenuRepository.Create(roleMenu); err != nil {
			return
		}
	}

	if err = a.roleRepository.Create(role); err != nil {
		return
	}

	a.casbinService.Enforcer.LoadPolicy()
	return role.ID, nil
}

func (a RoleService) Update(id string, role *models.Role) error {
	oRole, err := a.Get(id)
	if err != nil {
		return err
	} else if role.Name != oRole.Name {
		if err = a.Check(role); err != nil {
			return err
		}
	}

	role.ID = oRole.ID
	role.CreatedBy = oRole.CreatedBy
	role.CreatedAt = oRole.CreatedAt

	aRoleMenus, dRoleMenus := a.CompareRoleMenus(oRole.RoleMenus, role.RoleMenus)
	for _, aRoleMenu := range aRoleMenus {
		aRoleMenu.ID = uuid.MustString()
		aRoleMenu.RoleID = id

		if err := a.CheckRoleMenu(aRoleMenu); err != nil {
			return err
		}

		if err := a.roleMenuRepository.Create(aRoleMenu); err != nil {
			return err
		}
	}

	for _, dRoleMenu := range dRoleMenus {
		if err := a.roleMenuRepository.Delete(dRoleMenu.ID); err != nil {
			return err
		}
	}

	if err := a.roleRepository.Update(id, role); err != nil {
		return err
	}

	a.casbinService.Enforcer.LoadPolicy()
	return nil
}

func (a RoleService) Delete(id string) error {
	_, err := a.roleRepository.Get(id)
	if err != nil {
		return err
	}

	userQR, err := a.userRepository.Query(&models.UserQueryParam{
		RoleIDs: []string{id},
	})

	if err != nil {
		return err
	} else if userQR.Pagination.Total > 0 {
		return errors.RoleNotAllowDeleteWithUser
	}

	if err := a.roleMenuRepository.DeleteByRoleID(id); err != nil {
		return err
	}

	if err := a.roleRepository.Delete(id); err != nil {
		return err
	}

	a.casbinService.Enforcer.LoadPolicy()
	return nil
}

func (a RoleService) UpdateStatus(id string, status int) error {
	_, err := a.roleRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.roleRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	a.casbinService.Enforcer.LoadPolicy()
	return nil
}
