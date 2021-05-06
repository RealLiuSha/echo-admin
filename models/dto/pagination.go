package dto

type Pagination struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"page_size"`
}

type PaginationParam struct {
	Current  int `query:"current"`
	PageSize int `query:"page_size" validate:"max=128"`
}

func (a *PaginationParam) GetCurrent() int {
	return a.Current
}

func (a *PaginationParam) GetPageSize() int {
	pageSize := a.PageSize
	if a.PageSize == 0 {
		pageSize = 15
	}

	return pageSize
}
