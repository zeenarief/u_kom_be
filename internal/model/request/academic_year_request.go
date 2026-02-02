package request

type AcademicYearCreateRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" binding:"required,datetime=2006-01-02"`
}

type AcademicYearUpdateRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
}
