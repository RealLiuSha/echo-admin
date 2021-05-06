package repository

import (
	"gorm.io/gorm"

	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
)

// MenuActionRepository database structure
type MenuActionResourceRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewMenuActionResourceRepository creates a new menu action resource repository
func NewMenuActionResourceRepository(db lib.Database, logger lib.Logger) MenuActionResourceRepository {
	return MenuActionResourceRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a MenuActionResourceRepository) WithTrx(trxHandle *gorm.DB) MenuActionResourceRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a MenuActionResourceRepository) Query(param *models.MenuActionResourceQueryParam) (*models.MenuActionResourceQueryResult, error) {
	db := a.db.ORM.Model(&models.MenuActionResource{})

	if v := param.MenuID; v != "" {
		subQuery := a.db.ORM.Model(&models.MenuAction{}).
			Where("menu_id=?", v).
			Select("id")

		db = db.Where("action_id IN (?)", subQuery)
	}

	if v := param.MenuIDs; len(v) > 0 {
		subQuery := a.db.ORM.Model(&models.MenuAction{}).
			Where("menu_id IN (?)", v).
			Select("id")

		db = db.Where("action_id IN (?)", subQuery)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.MenuActionResources, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.MenuActionResourceQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a MenuActionResourceRepository) Get(id string) (*models.MenuActionResource, error) {
	menuActionResource := new(models.MenuActionResource)

	if ok, err := QueryOne(a.db.ORM.Model(menuActionResource).Where("id=?", id), menuActionResource); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return menuActionResource, nil
}

func (a MenuActionResourceRepository) Create(menuActionResource *models.MenuActionResource) error {
	result := a.db.ORM.Model(menuActionResource).Create(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a MenuActionResourceRepository) Update(id string, menuActionResource *models.MenuActionResource) error {
	result := a.db.ORM.Model(menuActionResource).Where("id=?", id).Updates(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a MenuActionResourceRepository) Delete(id string) error {
	menuActionResource := new(models.MenuActionResource)

	result := a.db.ORM.Model(menuActionResource).Where("id=?", id).Delete(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a MenuActionResourceRepository) DeleteByActionID(actionID string) error {
	menuActionResource := new(models.MenuActionResource)

	result := a.db.ORM.Model(menuActionResource).Where("action_id=?", actionID).Delete(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a MenuActionResourceRepository) DeleteByMenuID(menuID string) error {
	menuAction := new(models.MenuAction)
	menuActionResource := new(models.MenuActionResource)

	subQuery := a.db.ORM.Model(menuAction).
		Where("menu_id=?", menuID).Select("id")

	result := a.db.ORM.Model(menuActionResource).
		Where("action_id IN (?)", subQuery).Delete(menuActionResource)

	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
