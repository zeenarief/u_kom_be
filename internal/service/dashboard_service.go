package service

import (
	"errors"
	"time"
	"smart_school_be/internal/model/domain"
	"smart_school_be/internal/model/response"
	"smart_school_be/internal/repository"
)

type DashboardService interface {
	GetStats() (*response.DashboardStatsResponse, error)
	GetTeacherStats(userID string) (*response.TeacherDashboardStatsResponse, error)
}

type dashboardService struct {
	dashboardRepo          repository.DashboardRepository
	scheduleRepo           repository.ScheduleRepository
	attendanceRepo         repository.AttendanceRepository
	teachingAssignmentRepo repository.TeachingAssignmentRepository
	studentRepo            repository.StudentRepository
	employeeRepo           repository.EmployeeRepository
}

func NewDashboardService(
	dashboardRepo repository.DashboardRepository,
	scheduleRepo repository.ScheduleRepository,
	attendanceRepo repository.AttendanceRepository,
	taRepo repository.TeachingAssignmentRepository,
	studentRepo repository.StudentRepository,
	employeeRepo repository.EmployeeRepository,
) DashboardService {
	return &dashboardService{
		dashboardRepo:          dashboardRepo,
		scheduleRepo:           scheduleRepo,
		attendanceRepo:         attendanceRepo,
		teachingAssignmentRepo: taRepo,
		studentRepo:            studentRepo,
		employeeRepo:           employeeRepo,
	}
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

func (s *dashboardService) GetTeacherStats(userID string) (*response.TeacherDashboardStatsResponse, error) {
	// 0. Get Teacher (Employee) ID from UserID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil || employee == nil {
		return nil, errors.New("teacher not found for this user")
	}
	teacherID := employee.ID

	// 1. Get Today's Schedule count
	weekday := time.Now().Weekday()
	dayInt := int(weekday)
	if weekday == time.Sunday {
		dayInt = 7
	}

	todaySchedules, err := s.scheduleRepo.FindByTeacherIDAndDay(teacherID, dayInt)
	if err != nil {
		return nil, err
	}
	totalClassesToday := int64(len(todaySchedules))

	// 2. Count Pending Attendance
	var pendingAttendance int64
	today := time.Now()
	for _, sch := range todaySchedules {
		session, _ := s.attendanceRepo.FindSessionByScheduleDate(sch.ID, today)
		if session == nil {
			pendingAttendance++
		}
	}

	// 3. Count Total Students (unique students in classes taught by teacher)
	assignments, err := s.teachingAssignmentRepo.FindByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	uniqueStudentIDs := make(map[string]bool)
	processedClassrooms := make(map[string]bool)

	for _, a := range assignments {
		if processedClassrooms[a.ClassroomID] {
			continue
		}
		processedClassrooms[a.ClassroomID] = true

		studentsInClass, err := s.studentRepo.FindByClassroomID(a.ClassroomID)
		if err == nil {
			for _, st := range studentsInClass {
				uniqueStudentIDs[st.ID] = true
			}
		}
	}

	totalStudents := int64(len(uniqueStudentIDs))

	return &response.TeacherDashboardStatsResponse{
		TotalClassesToday: totalClassesToday,
		PendingAttendance: pendingAttendance,
		TotalStudents:     totalStudents,
	}, nil
}
