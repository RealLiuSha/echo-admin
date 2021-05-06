package services

import (
	"sort"

	"gorm.io/gorm"

	"github.com/RealLiuSha/echo-admin/api/repository"
	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
	"github.com/RealLiuSha/echo-admin/models/dto"
	"github.com/RealLiuSha/echo-admin/pkg/hash"
	"github.com/RealLiuSha/echo-admin/pkg/uuid"
)

// UserService service layer
type UserService struct {
	logger               lib.Logger
	config               lib.Config
	casbinService        CasbinService
	userRepository       repository.UserRepository
	userRoleRepository   repository.UserRoleRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
	roleRepository       repository.RoleRepository
	roleMenuRepository   repository.RoleMenuRepository
}

// NewUserService creates a new userservice
func NewUserService(
	logger lib.Logger,
	userRepository repository.UserRepository,
	userRoleRepository repository.UserRoleRepository,
	roleRepository repository.RoleRepository,
	roleMenuRepository repository.RoleMenuRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
	casbinService CasbinService,
	config lib.Config,
) UserService {
	return UserService{
		logger:               logger,
		config:               config,
		userRepository:       userRepository,
		userRoleRepository:   userRoleRepository,
		roleRepository:       roleRepository,
		roleMenuRepository:   roleMenuRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
		casbinService:        casbinService,
	}
}

func (a UserService) GetSuperAdmin() *models.User {
	admin := a.config.SuperAdmin
	return &models.User{
		ID:       admin.Username,
		Username: admin.Username,
		Realname: admin.Realname,
		Password: admin.Password,
	}
}

// WithTrx delegates transaction to repository database
func (a UserService) WithTrx(trxHandle *gorm.DB) UserService {
	a.userRepository = a.userRepository.WithTrx(trxHandle)
	a.userRoleRepository = a.userRoleRepository.WithTrx(trxHandle)

	return a
}

func (a UserService) Query(param *models.UserQueryParam) (userQR *models.UserQueryResult, err error) {
	if userQR, err = a.userRepository.Query(param); err != nil {
		return
	}

	uRoleQR, err := a.userRoleRepository.Query(
		&models.UserRoleQueryParam{UserIDs: userQR.List.ToIDs()},
	)

	if err != nil {
		return
	}

	m := uRoleQR.List.ToUserIDMap()
	for _, user := range userQR.List {
		if uRoles, ok := m[user.ID]; ok {
			user.UserRoles = uRoles
		}
	}

	return
}

func (a UserService) Verify(username, password string) (*models.User, error) {
	// super admin user
	admin := a.GetSuperAdmin()
	if admin.Username == username && admin.Password == password {
		return admin, nil
	}

	user, err := a.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if user.Password != hash.SHA256(password) {
		return nil, errors.UserInvalidPassword
	} else if user.Status != 1 {
		return nil, errors.UserIsDisable
	}

	return user, nil
}

func (a UserService) Check(user *models.User) error {
	if user.Username == a.GetSuperAdmin().Username {
		return errors.UserInvalidUsername
	}

	if qr, err := a.Query(&models.UserQueryParam{Username: user.Username}); err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.UserAlreadyExists
	}

	return nil
}

func (a UserService) GetUserInfo(ID string) (*models.UserInfo, error) {
	if a.GetSuperAdmin().ID == ID {
		user := a.GetSuperAdmin()
		return &models.UserInfo{
			ID:       user.Username,
			Username: user.Username,
			Realname: user.Realname,
		}, nil
	}

	user, err := a.Get(ID)
	if err != nil {
		return nil, err
	}

	userinfo := &models.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Realname: user.Realname,
	}

	userRoleQR, err := a.userRoleRepository.Query(&models.UserRoleQueryParam{
		UserID: ID,
	})

	if err != nil {
		return nil, err
	}

	if roleIDs := userRoleQR.List.ToRoleIDs(); len(roleIDs) > 0 {
		roleQR, err := a.roleRepository.Query(&models.RoleQueryParam{
			IDs:    roleIDs,
			Status: 1,
		})

		if err != nil {
			return nil, err
		}

		userinfo.Roles = roleQR.List
	}

	return userinfo, nil
}

func (a UserService) GetUserMenuTrees(ID string) (models.MenuTrees, error) {
	if a.GetSuperAdmin().ID == ID {
		menuQR, err := a.menuRepository.Query(&models.MenuQueryParam{
			Status:     1,
			OrderParam: dto.OrderParam{Key: "sequence", Direction: dto.OrderByASC},
		})

		if err != nil {
			return nil, err
		}

		return menuQR.List.ToMenuTrees(), nil
	}

	var (
		userRoleQR *models.UserRoleQueryResult
		roleMenuQR *models.RoleMenuQueryResult
		menuQR     *models.MenuQueryResult
		err        error
	)

	if userRoleQR, err = a.userRoleRepository.Query(&models.UserRoleQueryParam{
		UserID: ID,
	}); err != nil {
		return nil, err
	} else if len(userRoleQR.List) == 0 {
		return nil, errors.UserNoPermission
	}

	if roleMenuQR, err = a.roleMenuRepository.Query(&models.RoleMenuQueryParam{
		RoleIDs: userRoleQR.List.ToRoleIDs(),
	}); err != nil {
		return nil, err
	} else if len(roleMenuQR.List) == 0 {
		return nil, errors.UserNoPermission
	}

	if menuQR, err = a.menuRepository.Query(&models.MenuQueryParam{
		IDs:        roleMenuQR.List.ToMenuIDs(),
		Status:     1,
		OrderParam: dto.OrderParam{Key: "sequence", Direction: dto.OrderByASC},
	}); err != nil {
		return nil, err
	} else if len(menuQR.List) == 0 {
		return nil, errors.UserNoPermission
	}

	menuMap := menuQR.List.ToMap()
	// 获取授权菜单的父级菜单，判断哪些父级菜单不在之前的授权菜单中，存放于 parentIDs 切片
	var parentIDs []string
	for _, parentID := range menuQR.List.SplitParentIDs() {
		if _, ok := menuMap[parentID]; !ok {
			parentIDs = append(parentIDs, parentID)
		}
	}

	// 获取这些差异的父级菜单的信息，补充到menuResult.Data中
	if len(parentIDs) > 0 {
		parentMenuQR, err := a.menuRepository.Query(&models.MenuQueryParam{
			IDs: parentIDs,
		})

		if err != nil {
			return nil, err
		}

		menuQR.List = append(menuQR.List, parentMenuQR.List...)
	}

	sort.Sort(menuQR.List)
	return menuQR.List.ToMenuTrees(), nil
}

func (a UserService) GetByUsername(username string) (*models.User, error) {
	userQR, err := a.Query(
		&models.UserQueryParam{Username: username, QueryPassword: true},
	)

	if err != nil {
		return nil, err
	} else if len(userQR.List) == 0 {
		return nil, errors.UserRecordNotFound
	}

	// set schema
	user := userQR.List[0]

	// get user roles
	userRoleQR, err := a.userRoleRepository.Query(
		&models.UserRoleQueryParam{UserID: user.ID},
	)

	if err != nil {
		return nil, err
	}

	user.UserRoles = userRoleQR.List
	return user, nil
}

func (a UserService) Get(id string) (*models.User, error) {
	user, err := a.userRepository.Get(id)
	if err != nil {
		return nil, err
	}

	userRoleQR, err := a.userRoleRepository.Query(
		&models.UserRoleQueryParam{UserID: id},
	)

	if err != nil {
		return nil, err
	}

	user.UserRoles = userRoleQR.List
	return user, nil
}

func (a UserService) Create(user *models.User) (id string, err error) {
	if err = a.Check(user); err != nil {
		return
	}

	user.Password = hash.SHA256(user.Password)
	user.ID = uuid.MustString()

	for _, userRole := range user.UserRoles {
		userRole.ID = uuid.MustString()
		userRole.UserID = user.ID

		if err = a.userRoleRepository.Create(userRole); err != nil {
			return
		}
	}

	if err = a.userRepository.Create(user); err != nil {
		return
	}

	a.casbinService.Enforcer.LoadPolicy()
	return user.ID, nil
}

func (a UserService) Update(id string, user *models.User) error {
	oUser, err := a.Get(id)
	if err != nil {
		return err
	} else if user.Username != oUser.Username {
		if err := a.Check(user); err != nil {
			return err
		}
	}

	if user.Password != "" {
		user.Password = hash.SHA256(user.Password)
	} else {
		user.Password = oUser.Password
	}

	user.ID = oUser.ID
	user.CreatedBy = oUser.CreatedBy
	user.CreatedAt = oUser.CreatedAt

	aUserRoles, dUserRoles := a.CompareUserRoles(oUser.UserRoles, user.UserRoles)
	for _, aUserRole := range aUserRoles {
		aUserRole.ID = uuid.MustString()
		aUserRole.UserID = id
		if err := a.userRoleRepository.Create(aUserRole); err != nil {
			return err
		}
	}

	for _, dUserRole := range dUserRoles {
		if err := a.userRoleRepository.Delete(dUserRole.ID); err != nil {
			return err
		}
	}

	if err := a.userRepository.Update(id, user); err != nil {
		return err
	}

	a.casbinService.Enforcer.LoadPolicy()
	return nil
}

func (a UserService) Delete(id string) error {
	_, err := a.userRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.userRoleRepository.DeleteByUserID(id); err != nil {
		return err
	}

	a.casbinService.Enforcer.LoadPolicy()
	return a.userRepository.Delete(id)
}

func (a UserService) UpdateStatus(id string, status int) error {
	_, err := a.userRepository.Get(id)
	if err != nil {
		return err
	}

	if err = a.userRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	a.casbinService.Enforcer.LoadPolicy()
	return nil
}

func (a UserService) CompareUserRoles(oUserRoles, nUserRoles models.UserRoles) (aList, dList models.UserRoles) {
	oMap := oUserRoles.ToMap()
	nMap := nUserRoles.ToMap()

	for k, nUserRole := range nMap {
		if _, ok := oMap[k]; ok {
			delete(oMap, k)
			continue
		}

		aList = append(aList, nUserRole)
	}

	for _, oUserRole := range oMap {
		dList = append(dList, oUserRole)
	}

	return
}
