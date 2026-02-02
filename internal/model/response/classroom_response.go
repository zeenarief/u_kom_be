package response

import "time"

type ClassroomResponse struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Level        string               `json:"level"`
	Major        string               `json:"major"`
	Description  string               `json:"description"`
	AcademicYear AcademicYearResponse `json:"academic_year"` // Nested struct ringkas

	HomeroomTeacherID   *string   `json:"homeroom_teacher_id"`
	HomeroomTeacherName string    `json:"homeroom_teacher_name"` // Nama guru saja biar ringan
	TotalStudents       int64     `json:"total_students"`
	CreatedAt           time.Time `json:"created_at"`
}

type ClassroomDetailResponse struct {
	ClassroomResponse
	Students []StudentInClassResponse `json:"students"`
}

type StudentInClassResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	NISN     string `json:"nisn"`
	Gender   string `json:"gender"`
	Status   string `json:"status_in_class"` // Active, etc
}
