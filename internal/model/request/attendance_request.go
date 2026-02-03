package request

type AttendanceSubmitRequest struct {
	ScheduleID string `json:"schedule_id" binding:"required"`
	Date       string `json:"date" binding:"required,datetime=2006-01-02"` // YYYY-MM-DD
	Topic      string `json:"topic"`
	Notes      string `json:"notes"`

	// List Absensi Siswa
	Students []StudentAttendanceInput `json:"students" binding:"required,dive"`
}

type StudentAttendanceInput struct {
	StudentID string `json:"student_id" binding:"required"`
	Status    string `json:"status" binding:"required,oneof=PRESENT SICK PERMISSION ABSENT"`
	Notes     string `json:"notes"`
}
