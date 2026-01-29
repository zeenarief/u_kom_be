package service

import (
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
)

type DashboardService interface {
	GetStats() (*response.DashboardStatsResponse, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{dashboardRepo: dashboardRepo}
}

func (s *dashboardService) GetStats() (*response.DashboardStatsResponse, error) {
	// Gunakan Goroutine (WaitGroup) jika datanya sangat besar nanti.
	// Untuk sekarang, sequential saja sudah sangat cepat (< 10ms).

	students, _ := s.dashboardRepo.CountTable(&domain.Student{})
	employees, _ := s.dashboardRepo.CountTable(&domain.Employee{})
	parents, _ := s.dashboardRepo.CountTable(&domain.Parent{})
	users, _ := s.dashboardRepo.CountTable(&domain.User{})

	maleStudents, _ := s.dashboardRepo.CountStudentByGender("male")
	femaleStudents, _ := s.dashboardRepo.CountStudentByGender("female")

	return &response.DashboardStatsResponse{
		TotalStudents:  students,
		TotalEmployees: employees,
		TotalParents:   parents,
		TotalUsers:     users,
		StudentGender: struct {
			Male   int64 `json:"male"`
			Female int64 `json:"female"`
		}{
			Male:   maleStudents,
			Female: femaleStudents,
		},
	}, nil
}
