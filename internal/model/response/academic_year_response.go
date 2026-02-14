package response

import (
	"time"
	"smart_school_be/internal/utils"
)

type AcademicYearResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Status    string     `json:"status"`
	StartDate utils.Date `json:"start_date"`
	EndDate   utils.Date `json:"end_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
