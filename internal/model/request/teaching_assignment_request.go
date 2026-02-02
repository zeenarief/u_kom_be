package request

type TeachingAssignmentCreateRequest struct {
	ClassroomID string `json:"classroom_id" binding:"required"`
	SubjectID   string `json:"subject_id" binding:"required"`
	TeacherID   string `json:"teacher_id" binding:"required"`
}

// Tidak butuh update, biasanya flow-nya adalah Delete lalu Create ulang,
// atau timpa data (Upsert). Kita pakai Create saja (yang akan handle duplikat).
