package response

type DashboardStatsResponse struct {
	TotalStudents  int64 `json:"total_students"`
	TotalEmployees int64 `json:"total_employees"`
	TotalParents   int64 `json:"total_parents"`
	TotalUsers     int64 `json:"total_users"`

	// Bonus: Statistik Gender Siswa untuk Grafik
	StudentGender struct {
		Male   int64 `json:"male"`
		Female int64 `json:"female"`
	} `json:"student_gender"`
}
