package repository

import (
	"gorm.io/gorm"
	"u_kom_be/internal/model/domain"
)

type DashboardRepository interface {
	CountTable(model interface{}) (int64, error)
	CountStudentByGender(gender string) (int64, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) CountTable(model interface{}) (int64, error) {
	var count int64
	err := r.db.Model(model).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountStudentByGender(gender string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Student{}).Where("gender = ?", gender).Count(&count).Error
	return count, err
}
