package response

import "time"

// Response detail sesi beserta list siswanya
type AttendanceSessionDetailResponse struct {
	ID           string                     `json:"id"`
	Date         string                     `json:"date"`
	Topic        string                     `json:"topic"`
	ScheduleInfo ScheduleResponse           `json:"schedule_info"` // Reuse struct ScheduleResponse
	Details      []AttendanceDetailResponse `json:"details"`
	Summary      map[string]int             `json:"summary"` // Hadir: 30, Sakit: 1, dll
}

type AttendanceDetailResponse struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	NISN        string `json:"nisn"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}

// Response ringkas untuk history
type AttendanceHistoryResponse struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	SubjectName string    `json:"subject_name"`
	ClassName   string    `json:"class_name"`
	Topic       string    `json:"topic"`
	CountAbsent int       `json:"count_absent"` // Jumlah yg tidak hadir
}
