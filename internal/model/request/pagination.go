package request

import "strconv"

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type PaginationRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return DefaultPage
	}
	return p.Page
}

func (p *PaginationRequest) GetLimit() int {
	if p.Limit <= 0 {
		return DefaultLimit
	}
	if p.Limit > MaxLimit {
		return MaxLimit
	}
	return p.Limit
}

func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

// Helper to create from string values (e.g. from query params if not using binding)
func NewPaginationRequest(pageStr, limitStr string) PaginationRequest {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	return PaginationRequest{
		Page:  page,
		Limit: limit,
	}
}
