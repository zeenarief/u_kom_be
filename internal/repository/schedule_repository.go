package repository

import (
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type ScheduleRepository interface {
	Create(schedule *domain.Schedule) error
	FindByClassroomID(classroomID string) ([]domain.Schedule, error)
	FindByTeacherID(teacherID string) ([]domain.Schedule, error)
	FindByTeachingAssignmentID(taID string) ([]domain.Schedule, error)
	FindByID(id string) (*domain.Schedule, error)
	Delete(id string) error

	// Validasi Bentrok
	CheckClassroomConflict(classroomID string, day int, start, end string) (bool, error)
	CheckTeacherConflict(teacherID string, day int, start, end string) (bool, error)
}

type scheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(schedule *domain.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *scheduleRepository) FindByClassroomID(classroomID string) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	// Join dengan teaching_assignment untuk filter by classroom_id
	err := r.db.Preload("TeachingAssignment").
		Preload("TeachingAssignment.Subject").
		Preload("TeachingAssignment.Teacher").
		Preload("TeachingAssignment.Classroom").
		Joins("JOIN teaching_assignments ta ON ta.id = schedules.teaching_assignment_id").
		Where("ta.classroom_id = ?", classroomID).
		Order("day_of_week ASC, start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) FindByTeacherID(teacherID string) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := r.db.Preload("TeachingAssignment").
		Preload("TeachingAssignment.Subject").
		Preload("TeachingAssignment.Teacher").
		Preload("TeachingAssignment.Classroom").
		Joins("JOIN teaching_assignments ta ON ta.id = schedules.teaching_assignment_id").
		Where("ta.teacher_id = ?", teacherID).
		Order("day_of_week ASC, start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) FindByTeachingAssignmentID(taID string) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := r.db.Preload("TeachingAssignment").
		Preload("TeachingAssignment.Subject").
		Preload("TeachingAssignment.Teacher").
		Preload("TeachingAssignment.Classroom").
		Where("teaching_assignment_id = ?", taID).
		Order("day_of_week ASC, start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) FindByID(id string) (*domain.Schedule, error) {
	var schedule domain.Schedule
	err := r.db.Preload("TeachingAssignment").
		Preload("TeachingAssignment.Subject").
		Preload("TeachingAssignment.Teacher").
		Preload("TeachingAssignment.Classroom").
		First(&schedule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) Delete(id string) error {
	return r.db.Delete(&domain.Schedule{}, "id = ?", id).Error
}

// CheckClassroomConflict: Cek apakah KELAS ini sudah ada jadwal di jam tersebut
func (r *scheduleRepository) CheckClassroomConflict(classroomID string, day int, start, end string) (bool, error) {
	var count int64
	// Logika Overlap: (StartA < EndB) AND (EndA > StartB)
	err := r.db.Model(&domain.Schedule{}).
		Joins("JOIN teaching_assignments ta ON ta.id = schedules.teaching_assignment_id").
		Where("ta.classroom_id = ?", classroomID).
		Where("schedules.day_of_week = ?", day).
		Where("? < schedules.end_time AND ? > schedules.start_time", start, end).
		Count(&count).Error

	return count > 0, err
}

// CheckTeacherConflict: Cek apakah GURU ini sudah mengajar di kelas lain di jam tersebut
func (r *scheduleRepository) CheckTeacherConflict(teacherID string, day int, start, end string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Schedule{}).
		Joins("JOIN teaching_assignments ta ON ta.id = schedules.teaching_assignment_id").
		Where("ta.teacher_id = ?", teacherID).
		Where("schedules.day_of_week = ?", day).
		Where("? < schedules.end_time AND ? > schedules.start_time", start, end).
		Count(&count).Error

	return count > 0, err
}
