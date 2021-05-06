package repository

import (
	"gorm.io/gorm"

	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
)

// UserRoleRepository database structure
type UserRoleRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewUserRepository creates a new user role repository
func NewUserRoleRepository(db lib.Database, logger lib.Logger) UserRoleRepository {
	return UserRoleRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a UserRoleRepository) WithTrx(trxHandle *gorm.DB) UserRoleRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a UserRoleRepository) Query(param *models.UserRoleQueryParam) (*models.UserRoleQueryResult, error) {
	db := a.db.ORM.Model(models.UserRole{})

	if v := param.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}
	if v := param.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.UserRoles, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.UserRoleQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a UserRoleRepository) Get(id string) (*models.UserRole, error) {
	userRole := new(models.UserRole)

	if ok, err := QueryOne(a.db.ORM.Model(userRole).Where("id=?", id), userRole); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return userRole, nil
}

func (a UserRoleRepository) Create(userRole *models.UserRole) error {
	result := a.db.ORM.Model(userRole).Create(userRole)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a UserRoleRepository) Update(id string, userRole *models.UserRole) error {
	result := a.db.ORM.Model(userRole).Where("id=?", id).Updates(userRole)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a UserRoleRepository) Delete(id string) error {
	userRole := new(models.UserRole)

	result := a.db.ORM.Model(userRole).Where("id=?", id).Delete(userRole)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a UserRoleRepository) DeleteByUserID(userID string) error {
	userRole := new(models.UserRole)

	result := a.db.ORM.Model(userRole).Where("user_id=?", userID).Delete(userRole)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
