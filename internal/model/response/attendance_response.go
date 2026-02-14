package response

import (
	"u_kom_be/internal/utils"
)

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
	ID          string     `json:"id"`
	Date        utils.Date `json:"date"`
	ScheduleID  string     `json:"schedule_id"` // Added for Edit Feature
	SubjectName string     `json:"subject_name"`
	ClassName   string     `json:"class_name"`
	Topic       string     `json:"topic"`
	CountAbsent int        `json:"count_absent"` // Jumlah yg tidak hadir
}
