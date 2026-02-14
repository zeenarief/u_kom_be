package service

import (
	"time"
	"u_kom_be/internal/apperrors"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
)

type ScheduleService interface {
	Create(req request.ScheduleCreateRequest) (*response.ScheduleResponse, error)
	GetByClassroom(classroomID string) ([]response.ScheduleResponse, error)
	GetByTeacher(teacherID string) ([]response.ScheduleResponse, error)
	GetByTeachingAssignment(taID string) ([]response.ScheduleResponse, error)
	Delete(id string) error
	GetTodaySchedule(userID string) ([]response.ScheduleResponse, error)
}

type scheduleService struct {
	repo                   repository.ScheduleRepository
	teachingAssignmentRepo repository.TeachingAssignmentRepository // Butuh ini untuk cek TeacherID & ClassID
	employeeRepo           repository.EmployeeRepository
}

func NewScheduleService(
	repo repository.ScheduleRepository,
	taRepo repository.TeachingAssignmentRepository,
	employeeRepo repository.EmployeeRepository,
) ScheduleService {
	return &scheduleService{
		repo:                   repo,
		teachingAssignmentRepo: taRepo,
		employeeRepo:           employeeRepo,
	}
}

// Helper: Convert int day to string indonesian
func getDayName(day int) string {
	days := []string{"", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Ahad"}
	if day >= 1 && day <= 7 {
		return days[day]
	}
	return "Unknown"
}

func (s *scheduleService) toResponse(d *domain.Schedule) response.ScheduleResponse {
	return response.ScheduleResponse{
		ID:            d.ID,
		DayOfWeek:     d.DayOfWeek,
		DayName:       getDayName(d.DayOfWeek),
		StartTime:     d.StartTime,
		EndTime:       d.EndTime,
		SubjectName:   d.TeachingAssignment.Subject.Name,
		TeacherName:   d.TeachingAssignment.Teacher.FullName,
		ClassroomName: d.TeachingAssignment.Classroom.Name,
	}
}

func (s *scheduleService) Create(req request.ScheduleCreateRequest) (*response.ScheduleResponse, error) {
	// 1. Ambil detail Assignment (untuk tahu Guru & Kelas siapa)
	// Kita reuse FindOne dari repo assignment (tapi repo aslinya butuh classroom_id & subject_id)
	// Jadi lebih aman kita fetch by ID assignment-nya langsung.
	// *Catatan: Kita perlu tambahkan FindByID di TeachingAssignmentRepository dulu (lihat bawah)*
	assignment, err := s.teachingAssignmentRepo.FindByID(req.TeachingAssignmentID)
	if err != nil || assignment == nil {
		return nil, apperrors.NewNotFoundError("Teaching assignment not found")
	}

	// 2. Validasi Logic: EndTime harus > StartTime
	if req.EndTime <= req.StartTime {
		return nil, apperrors.NewBadRequestError("End time must be greater than start time")
	}

	// 3. Cek Bentrok KELAS (Apakah kelas ini sedang belajar mapel lain?)
	conflictClass, err := s.repo.CheckClassroomConflict(assignment.ClassroomID, req.DayOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	if conflictClass {
		return nil, apperrors.NewConflictError("Classroom is occupied at this time")
	}

	// 4. Cek Bentrok GURU (Apakah guru ini sedang mengajar di kelas lain?)
	conflictTeacher, err := s.repo.CheckTeacherConflict(assignment.TeacherID, req.DayOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	if conflictTeacher {
		return nil, apperrors.NewConflictError("Teacher is teaching in another class at this time")
	}

	// 5. Simpan
	schedule := &domain.Schedule{
		TeachingAssignmentID: req.TeachingAssignmentID,
		DayOfWeek:            req.DayOfWeek,
		StartTime:            req.StartTime,
		EndTime:              req.EndTime,
	}

	if err := s.repo.Create(schedule); err != nil {
		return nil, err
	}

	// Attach relasi manual untuk response
	schedule.TeachingAssignment = *assignment

	res := s.toResponse(schedule)
	return &res, nil
}

func (s *scheduleService) GetByClassroom(classroomID string) ([]response.ScheduleResponse, error) {
	data, err := s.repo.FindByClassroomID(classroomID)
	if err != nil {
		return nil, err
	}
	var res []response.ScheduleResponse
	for _, d := range data {
		res = append(res, s.toResponse(&d))
	}
	return res, nil
}

func (s *scheduleService) GetByTeacher(teacherID string) ([]response.ScheduleResponse, error) {
	data, err := s.repo.FindByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}
	var res []response.ScheduleResponse
	for _, d := range data {
		res = append(res, s.toResponse(&d))
	}
	return res, nil
}

func (s *scheduleService) GetByTeachingAssignment(taID string) ([]response.ScheduleResponse, error) {
	data, err := s.repo.FindByTeachingAssignmentID(taID)
	if err != nil {
		return nil, err
	}
	var res []response.ScheduleResponse
	for _, d := range data {
		res = append(res, s.toResponse(&d))
	}
	return res, nil
}

func (s *scheduleService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *scheduleService) GetTodaySchedule(userID string) ([]response.ScheduleResponse, error) {
	// Find TeacherID by UserID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil || employee == nil {
		return nil, apperrors.NewNotFoundError("Teacher/Employee profile not found")
	}
	teacherID := employee.ID

	// Map time.Weekday (Sun=0, Mon=1...) to DB (Mon=1... Sun=7)
	weekday := time.Now().Weekday()
	dayInt := int(weekday)
	if weekday == time.Sunday {
		dayInt = 7
	}

	data, err := s.repo.FindByTeacherIDAndDay(teacherID, dayInt)
	if err != nil {
		return nil, err
	}
	var res []response.ScheduleResponse
	for _, d := range data {
		res = append(res, s.toResponse(&d))
	}
	return res, nil
}
