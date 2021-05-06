package repository

import (
	"github.com/RealLiuSha/echo-admin/models/dto"
	"gorm.io/gorm"
)

func QueryPagination(db *gorm.DB, pp dto.PaginationParam, out interface{}) (*dto.Pagination, error) {
	pagination := new(dto.Pagination)

	total, err := QeuryPage(db, pp, out)
	if err != nil {
		return pagination, err
	}

	pagination.Current = pp.GetCurrent()
	pagination.PageSize = pp.GetPageSize()
	pagination.Total = total

	return pagination, nil
}

func QeuryPage(db *gorm.DB, pp dto.PaginationParam, out interface{}) (n int64, err error) {
	n, err = QueryCount(db)
	if err != nil {
		return
	} else if n == 0 {
		return
	}

	current, pageSize := pp.GetCurrent(), pp.GetPageSize()
	if current > 0 && pageSize > 0 {
		db = db.Offset((current - 1) * pageSize).Limit(pageSize)
	} else if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	err = db.Find(out).Error
	return
}

func QueryOne(db *gorm.DB, out interface{}) (bool, error) {
	result := db.First(out)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func QueryCount(db *gorm.DB) (n int64, err error) {
	result := db.Count(&n)
	if err = result.Error; err != nil {
		return
	}

	return
}
