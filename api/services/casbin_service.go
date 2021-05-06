package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/RealLiuSha/echo-admin/models/dto"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"go.uber.org/zap"

	"github.com/RealLiuSha/echo-admin/api/repository"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
)

type CasbinAdapter struct {
	logger                       lib.Logger
	userRepository               repository.UserRepository
	userRoleRepository           repository.UserRoleRepository
	roleRepository               repository.RoleRepository
	roleMenuRepository           repository.RoleMenuRepository
	menuActionResourceRepository repository.MenuActionResourceRepository
}

type CasbinLogger struct {
	zap     *zap.Logger
	enabled bool
}

// CasbinService service layer
type CasbinService struct {
	Enforcer *casbin.SyncedEnforcer
}

// NewCasbinService creates a new userservice
func NewCasbinService(
	logger lib.Logger,
	config lib.Config,

	userRepository repository.UserRepository,
	userRoleRepository repository.UserRoleRepository,
	roleRepository repository.RoleRepository,
	roleMenuRepository repository.RoleMenuRepository,
	menuActionResourceRepository repository.MenuActionResourceRepository,
) CasbinService {
	adapter := &CasbinAdapter{
		logger:                       logger,
		userRepository:               userRepository,
		userRoleRepository:           userRoleRepository,
		roleRepository:               roleRepository,
		roleMenuRepository:           roleMenuRepository,
		menuActionResourceRepository: menuActionResourceRepository,
	}

	enforcer, err := casbin.NewSyncedEnforcer(
		config.Casbin.Model,
	)

	if err != nil {
		logger.Zap.Fatalf("error to new casbin enforcer: %v", err)
	}

	enforcer.EnableEnforce(true)
	enforcer.EnableLog(config.Casbin.Debug)
	enforcer.SetLogger(&CasbinLogger{
		zap:     logger.DesugarZap.With(zap.String("module", "casbin")),
		enabled: true,
	})

	service := CasbinService{
		Enforcer: enforcer,
	}

	err = enforcer.InitWithModelAndAdapter(enforcer.GetModel(), adapter)
	if err != nil {
		logger.Zap.Fatalf("error to init model and adapter: %v", err)
	}

	if config.Casbin.AutoLoad {
		enforcer.StartAutoLoadPolicy(time.Duration(config.Casbin.AutoLoadInternal) * time.Second)
	}

	return service
}

// LoadPolicy loads all policy rules from the storage.
func (a CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	err := a.loadRolePolicy(model)
	if err != nil {
		a.logger.Zap.Errorf("Load casbin role policy error: %s", err.Error())
		return err
	}

	err = a.loadUserPolicy(model)
	if err != nil {
		a.logger.Zap.Errorf("Load casbin user policy error: %s", err.Error())
		return err
	}

	return nil
}

// load role policy (p,role_id,path,method)
func (a CasbinAdapter) loadRolePolicy(m casbinModel.Model) error {
	paginationParam := dto.PaginationParam{PageSize: 9999, Current: 1}
	roleQR, err := a.roleRepository.Query(&models.RoleQueryParam{
		Status: 1, PaginationParam: paginationParam,
	})

	if err != nil {
		return err
	} else if len(roleQR.List) == 0 {
		return nil
	}

	roleMenuQR, err := a.roleMenuRepository.Query(&models.RoleMenuQueryParam{
		PaginationParam: paginationParam,
	})

	if err != nil {
		return err
	}

	mRoleMenus := roleMenuQR.List.ToRoleIDMap()

	menuResourceQR, err := a.menuActionResourceRepository.Query(
		&models.MenuActionResourceQueryParam{PaginationParam: paginationParam},
	)

	if err != nil {
		return err
	}

	mMenuResources := menuResourceQR.List.ToActionIDMap()

	for _, role := range roleQR.List {
		mcache := make(map[string]struct{})
		roleMenus, ok := mRoleMenus[role.ID]
		if !ok {
			continue
		}

		for _, actionID := range roleMenus.ToActionIDs() {
			mrs, ok := mMenuResources[actionID]
			if !ok {
				continue
			}

			for _, mr := range mrs {
				if mr.Path == "" || mr.Method == "" {
					continue
				} else if _, ok := mcache[mr.Path+mr.Method]; ok {
					continue
				}

				mcache[mr.Path+mr.Method] = struct{}{}
				line := fmt.Sprintf("p,%s,%s,%s", role.ID, mr.Path, mr.Method)
				persist.LoadPolicyLine(line, m)
			}
		}
	}

	return nil
}

// load user policy (g,user_id,role_id)
func (a CasbinAdapter) loadUserPolicy(m casbinModel.Model) error {
	paginationParam := dto.PaginationParam{PageSize: 9999, Current: 1}

	userQR, err := a.userRepository.Query(&models.UserQueryParam{
		Status: 1, PaginationParam: paginationParam,
	})

	if err != nil {
		return err
	}

	if len(userQR.List) > 0 {
		userRoleQR, err := a.userRoleRepository.Query(&models.UserRoleQueryParam{
			PaginationParam: paginationParam,
		})

		if err != nil {
			return err
		}

		mUserRoles := userRoleQR.List.ToUserIDMap()
		for _, uitem := range userQR.List {
			urs, ok := mUserRoles[uitem.ID]
			if !ok {
				continue
			}

			for _, ur := range urs {
				line := fmt.Sprintf("g,%s,%s", ur.UserID, ur.RoleID)
				persist.LoadPolicyLine(line, m)
			}
		}
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a CasbinAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}

func (a *CasbinLogger) EnableLog(enable bool) {
	a.enabled = enable
}

func (a *CasbinLogger) IsEnabled() bool {
	return a.enabled
}

func (a *CasbinLogger) LogModel(model [][]string) {
	if !a.enabled {
		return
	}
	var str strings.Builder
	str.WriteString("Model: ")
	for _, v := range model {
		str.WriteString(fmt.Sprintf("%v", v))
	}

	a.zap.Info(str.String())
}

func (a *CasbinLogger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	if !a.enabled {
		return
	}

	var reqStr strings.Builder
	reqStr.WriteString("Request: ")
	for i, rval := range request {
		if i != len(request)-1 {
			reqStr.WriteString(fmt.Sprintf("%v, ", rval))
		} else {
			reqStr.WriteString(fmt.Sprintf("%v", rval))
		}
	}

	reqStr.WriteString(fmt.Sprintf(" ---> %t, ", result))
	reqStr.WriteString("Hit Policy: ")

	for i, pval := range explains {
		if i != len(explains)-1 {
			reqStr.WriteString(fmt.Sprintf("%v, ", pval))
		} else {
			reqStr.WriteString(fmt.Sprintf("%v ", pval))
		}
	}

	a.zap.Info(reqStr.String())
}

func (a *CasbinLogger) LogPolicy(policy map[string][][]string) {
	if !a.enabled {
		return
	}

	var str strings.Builder
	str.WriteString("Policy: ")
	for k, v := range policy {
		str.WriteString(fmt.Sprintf("%s : %v", k, v))
	}

	a.zap.Info(str.String())
}

func (a *CasbinLogger) LogRole(roles []string) {
	if !a.enabled {
		return
	}

	str := fmt.Sprintf("Roles: %s", roles)
	a.zap.Info(str)
}
