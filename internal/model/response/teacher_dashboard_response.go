package response

type TeacherDashboardStatsResponse struct {
	TotalClassesToday int64 `json:"total_classes_today"`
	TotalStudents     int64 `json:"total_students"`
	PendingAttendance int64 `json:"pending_attendance"`
}
