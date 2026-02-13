package response

type TeachingAssignmentResponse struct {
	ID        string                              `json:"id"`
	Classroom TeachingAssignmentClassroomResponse `json:"classroom"`
	Subject   TeachingAssignmentSubjectResponse   `json:"subject"`
	Teacher   TeachingAssignmentTeacherResponse   `json:"teacher"`
}

type TeachingAssignmentClassroomResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level string `json:"level"`
	Major string `json:"major"`
}

type TeachingAssignmentSubjectResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type TeachingAssignmentTeacherResponse struct {
	ID   string                         `json:"id"`
	NIP  string                         `json:"nip"`
	User TeachingAssignmentUserResponse `json:"user"`
}

type TeachingAssignmentUserResponse struct {
	Name string `json:"name"`
}
