package response

type ScheduleResponse struct {
	ID            string `json:"id"`
	DayOfWeek     int    `json:"day_of_week"` // 1
	DayName       string `json:"day_name"`    // "Senin"
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	SubjectName   string `json:"subject_name"`
	TeacherName   string `json:"teacher_name"`
	ClassroomName string `json:"classroom_name"`
}
