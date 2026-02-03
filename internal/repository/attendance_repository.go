package repository

import (
	"time"
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type AttendanceRepository interface {
	CreateSession(session *domain.AttendanceSession) error
	FindSessionByScheduleDate(scheduleID string, date time.Time) (*domain.AttendanceSession, error)
	FindSessionByID(id string) (*domain.AttendanceSession, error)
	GetHistoryByTeacher(teacherID string) ([]domain.AttendanceSession, error)
	// Update logic jika guru ingin mengedit absen
	UpdateSession(session *domain.AttendanceSession) error
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) CreateSession(session *domain.AttendanceSession) error {
	// Menggunakan Transaksi GORM (Create session + Details otomatis jika struct terisi)
	return r.db.Create(session).Error
}

func (r *attendanceRepository) FindSessionByScheduleDate(scheduleID string, date time.Time) (*domain.AttendanceSession, error) {
	var session domain.AttendanceSession
	err := r.db.Preload("Details").Preload("Details.Student").
		Where("schedule_id = ? AND date = ?", scheduleID, date).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *attendanceRepository) FindSessionByID(id string) (*domain.AttendanceSession, error) {
	var session domain.AttendanceSession
	err := r.db.Preload("Schedule").
		Preload("Schedule.TeachingAssignment.Subject").
		Preload("Schedule.TeachingAssignment.Classroom").
		Preload("Details").
		Preload("Details.Student").
		First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *attendanceRepository) GetHistoryByTeacher(teacherID string) ([]domain.AttendanceSession, error) {
	var sessions []domain.AttendanceSession
	// Join kompleks untuk mendapatkan sesi berdasarkan Guru
	err := r.db.Joins("JOIN schedules s ON s.id = attendance_sessions.schedule_id").
		Joins("JOIN teaching_assignments ta ON ta.id = s.teaching_assignment_id").
		Preload("Schedule.TeachingAssignment.Subject").
		Preload("Schedule.TeachingAssignment.Classroom").
		Where("ta.teacher_id = ?", teacherID).
		Order("attendance_sessions.date DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *attendanceRepository) UpdateSession(session *domain.AttendanceSession) error {
	// Logic update (biasanya hapus details lama, insert baru, atau update one-by-one)
	// Untuk simplifikasi awal, kita gunakan Save session (header) dulu
	return r.db.Save(session).Error
}
