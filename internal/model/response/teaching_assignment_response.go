package response

type TeachingAssignmentResponse struct {
	ID            string `json:"id"`
	ClassroomName string `json:"classroom_name"`
	SubjectName   string `json:"subject_name"`
	TeacherName   string `json:"teacher_name"`
	TeacherNIP    string `json:"teacher_nip"`
}
