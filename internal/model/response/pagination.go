package response

import "math"

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
}

type PaginatedData struct {
	Items interface{}    `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

func NewPaginatedData(items interface{}, totalItems int64, page, pageSize int) PaginatedData {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	if totalPages < 0 {
		totalPages = 0
	}

	return PaginatedData{
		Items: items,
		Meta: PaginationMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			PageSize:    pageSize,
			TotalItems:  totalItems,
		},
	}
}
