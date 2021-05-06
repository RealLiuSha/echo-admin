package repository

import (
	"gorm.io/gorm"

	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
)

// MenuActionRepository database structure
type MenuActionRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewMenuActionRepository creates a new menu action repository
func NewMenuActionRepository(db lib.Database, logger lib.Logger) MenuActionRepository {
	return MenuActionRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a MenuActionRepository) WithTrx(trxHandle *gorm.DB) MenuActionRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a *MenuActionRepository) Query(param *models.MenuActionQueryParam) (*models.MenuActionQueryResult, error) {
	db := a.db.ORM.Model(&models.MenuAction{})

	if v := param.MenuID; v != "" {
		db = db.Where("menu_id=?", v)
	}

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.MenuActions, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.MenuActionQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a MenuActionRepository) Get(id string) (*models.MenuAction, error) {
	menuAction := new(models.MenuAction)

	if ok, err := QueryOne(a.db.ORM.Model(menuAction).Where("id=?", id), menuAction); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return menuAction, nil
}

func (a MenuActionRepository) Create(menuAction *models.MenuAction) error {
	result := a.db.ORM.Model(menuAction).Create(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a MenuActionRepository) Update(id string, menuAction *models.MenuAction) error {
	result := a.db.ORM.Model(menuAction).Where("id=?", id).Updates(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a MenuActionRepository) Delete(id string) error {
	menuAction := new(models.MenuAction)

	result := a.db.ORM.Model(menuAction).Where("id=?", id).Delete(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a MenuActionRepository) DeleteByMenuID(menuID string) error {
	menuAction := new(models.MenuAction)

	result := a.db.ORM.Model(menuAction).Where("menu_id=?", menuID).Delete(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
