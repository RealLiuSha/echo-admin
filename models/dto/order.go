package dto

import (
	"fmt"
)

type OrderDirection string

const (
	OrderByASC      OrderDirection = "ASC"
	OrderByDESC     OrderDirection = "DESC"
	OrderDefaultKey                = "record_id"
)

type OrderParam struct {
	Key       string         `query:"order_key"`
	Direction OrderDirection `query:"order_direction"`
}

func (a OrderParam) ParseOrder() string {
	if a.Key == "" {
		a.Key = OrderDefaultKey
	}

	key := a.Key
	direction := "DESC"
	if a.Direction == OrderByASC {
		direction = "ASC"
	}

	return fmt.Sprintf("%s %s", key, direction)
}
