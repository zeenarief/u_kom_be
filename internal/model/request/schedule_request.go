package request

type ScheduleCreateRequest struct {
	TeachingAssignmentID string `json:"teaching_assignment_id" binding:"required"`
	DayOfWeek            int    `json:"day_of_week" binding:"required,min=1,max=7"`
	StartTime            string `json:"start_time" binding:"required"` // Format HH:MM
	EndTime              string `json:"end_time" binding:"required"`   // Format HH:MM
}
