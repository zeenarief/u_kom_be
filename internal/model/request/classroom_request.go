package request

type ClassroomCreateRequest struct {
	AcademicYearID    string  `json:"academic_year_id" binding:"required"`
	HomeroomTeacherID *string `json:"homeroom_teacher_id"` // Optional
	Name              string  `json:"name" binding:"required"`
	Level             string  `json:"level" binding:"required"` // 10, 11, 12
	Major             string  `json:"major"`                    // IPA, IPS
	Description       string  `json:"description"`
}

type ClassroomUpdateRequest struct {
	HomeroomTeacherID *string `json:"homeroom_teacher_id"`
	Name              string  `json:"name"`
	Level             string  `json:"level"`
	Major             string  `json:"major"`
	Description       string  `json:"description"`
}

// Request untuk menambahkan siswa ke kelas (Bulk)
type AddStudentsToClassRequest struct {
	StudentIDs []string `json:"student_ids" binding:"required"`
}
