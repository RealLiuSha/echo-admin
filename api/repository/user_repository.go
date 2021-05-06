package repository

import (
	"gorm.io/gorm"

	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
)

// UserRepository database structure
type UserRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db lib.Database, logger lib.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a UserRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

// GetAll gets all users
func (a UserRepository) Query(param *models.UserQueryParam) (*models.UserQueryResult, error) {
	db := a.db.ORM.Model(&models.User{})

	if v := param.QueryPassword; !v {
		db = db.Omit("password")
	}

	if v := param.Username; v != "" {
		db = db.Where("username = (?)", v)
	}

	if v := param.Realname; v != "" {
		db = db.Where("realname = (?)", v)
	}

	if v := param.Status; v != 0 {
		db = db.Where("status = (?)", v)
	}

	if v := param.RoleIDs; len(v) > 0 {
		subQuery := a.db.ORM.Model(&models.UserRole{}).
			Select("user_id").
			Where("role_id IN (?)", v)

		db = db.Where("id IN (?)", subQuery)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("username LIKE ? OR realname LIKE ? OR phone LIKE ? OR email LIKE ?", v, v, v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Users, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.UserQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a UserRepository) Get(id string) (*models.User, error) {
	user := new(models.User)

	if ok, err := QueryOne(a.db.ORM.Model(user).Where("id=?", id), user); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return user, nil
}

func (a UserRepository) Create(user *models.User) error {
	result := a.db.ORM.Model(user).Create(user)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a UserRepository) Update(id string, user *models.User) error {
	result := a.db.ORM.Model(user).Where("id=?", id).Updates(user)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a UserRepository) Delete(id string) error {
	user := new(models.User)

	result := a.db.ORM.Model(user).Where("id=?", id).Delete(user)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a UserRepository) UpdateStatus(id string, status int) error {
	user := new(models.User)

	result := a.db.ORM.Model(user).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a UserRepository) UpdatePassword(id, password string) error {
	user := new(models.User)

	result := a.db.ORM.Model(user).Where("id=?", id).Update("password", password)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
