package repository

import (
	"smart_school_be/internal/model/domain"
	"time"

	"gorm.io/gorm"
)

type AttendanceRepository interface {
	CreateSession(session *domain.AttendanceSession) error
	FindSessionByScheduleDate(scheduleID string, date time.Time) (*domain.AttendanceSession, error)
	FindSessionByID(id string) (*domain.AttendanceSession, error)
	GetHistoryByTeacher(teacherID string, limit, offset int) ([]domain.AttendanceSession, int64, error)
	GetHistoryByTeachingAssignmentID(taID string, limit, offset int) ([]domain.AttendanceSession, int64, error)
	// Update logic jika guru ingin mengedit absen
	UpdateSession(session *domain.AttendanceSession, newDetails []domain.AttendanceDetail) error
	DeleteSession(id string) error
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

// UPDATE: Pastikan Preload Details ada di sini
func (r *attendanceRepository) FindSessionByScheduleDate(scheduleID string, date time.Time) (*domain.AttendanceSession, error) {
	var session domain.AttendanceSession

	// PERBAIKAN: Gunakan DATE() dan format string tanggal yyyy-mm-dd
	// Ini memastikan kita membandingkan tanggalnya saja, tanpa peduli jam 00:00 atau 07:00
	dateString := date.Format("2006-01-02")

	err := r.db.Preload("Details").Preload("Details.Student").
		Where("schedule_id = ? AND DATE(date) = ?", scheduleID, dateString).
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

func (r *attendanceRepository) GetHistoryByTeacher(teacherID string, limit, offset int) ([]domain.AttendanceSession, int64, error) {
	var sessions []domain.AttendanceSession
	var total int64
	// Join kompleks untuk mendapatkan sesi berdasarkan Guru
	query := r.db.Model(&domain.AttendanceSession{}).
		Joins("JOIN schedules s ON s.id = attendance_sessions.schedule_id").
		Joins("JOIN teaching_assignments ta ON ta.id = s.teaching_assignment_id").
		Where("ta.teacher_id = ?", teacherID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Schedule.TeachingAssignment.Subject").
		Preload("Schedule.TeachingAssignment.Classroom").
		Preload("Details"). // Added Preload Details
		Order("attendance_sessions.date DESC").
		Limit(limit).Offset(offset).
		Find(&sessions).Error
	return sessions, total, err
}

func (r *attendanceRepository) GetHistoryByTeachingAssignmentID(taID string, limit, offset int) ([]domain.AttendanceSession, int64, error) {
	var sessions []domain.AttendanceSession
	var total int64
	query := r.db.Model(&domain.AttendanceSession{}).
		Select("attendance_sessions.*").
		Joins("JOIN schedules s ON s.id = attendance_sessions.schedule_id").
		Joins("JOIN teaching_assignments ta ON ta.id = s.teaching_assignment_id").
		Where("ta.id = ?", taID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Schedule.TeachingAssignment.Subject").
		Preload("Schedule.TeachingAssignment.Classroom").
		Preload("Details"). // Added Preload Details
		Order("attendance_sessions.date DESC").
		Limit(limit).Offset(offset).
		Find(&sessions).Error
	return sessions, total, err
}

// NEW: Update Session dengan Transaksi
func (r *attendanceRepository) UpdateSession(session *domain.AttendanceSession, newDetails []domain.AttendanceDetail) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update Header (Topic, Notes)
		if err := tx.Model(session).Updates(domain.AttendanceSession{
			Topic: session.Topic,
			Notes: session.Notes,
		}).Error; err != nil {
			return err
		}

		// 2. Hapus Detail Lama (Hard Delete berdasarkan Session ID)
		if err := tx.Where("attendance_session_id = ?", session.ID).Delete(&domain.AttendanceDetail{}).Error; err != nil {
			return err
		}

		// 3. Masukkan Detail Baru
		// Kita harus set ID session ke detail baru sebelum insert
		for i := range newDetails {
			newDetails[i].AttendanceSessionID = session.ID
		}

		if err := tx.Create(&newDetails).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *attendanceRepository) DeleteSession(id string) error {
	return r.db.Delete(&domain.AttendanceSession{}, "id = ?", id).Error
}
