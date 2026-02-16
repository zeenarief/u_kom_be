package repository

import (
	"smart_school_be/internal/model/domain"

	"gorm.io/gorm"
)

type DonorRepository interface {
	Create(donor *domain.Donor) error
	Update(donor *domain.Donor) error
	FindByID(id string) (*domain.Donor, error)
	FindByPhone(phone string) (*domain.Donor, error)
	FindAll(name string, limit, offset int) ([]domain.Donor, int64, error)
}

type DonationRepository interface {
	Create(donation *domain.Donation) error
	Update(donation *domain.Donation) error
	FindByID(id string) (*domain.Donation, error)
	FindAll(filter map[string]interface{}, limit, offset int) ([]domain.Donation, int64, error)
}

type donorRepository struct {
	db *gorm.DB
}

type donationRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) (DonorRepository, DonationRepository) {
	return &donorRepository{db}, &donationRepository{db}
}

// --- Donor Repository Implementation ---

func (r *donorRepository) Create(donor *domain.Donor) error {
	return r.db.Create(donor).Error
}

func (r *donorRepository) Update(donor *domain.Donor) error {
	return r.db.Save(donor).Error
}

func (r *donorRepository) FindByID(id string) (*domain.Donor, error) {
	var donor domain.Donor
	err := r.db.First(&donor, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &donor, nil
}

func (r *donorRepository) FindByPhone(phone string) (*domain.Donor, error) {
	var donor domain.Donor
	err := r.db.First(&donor, "phone = ?", phone).Error
	if err != nil {
		return nil, err
	}
	return &donor, nil
}

func (r *donorRepository) FindAll(name string, limit, offset int) ([]domain.Donor, int64, error) {
	var donors []domain.Donor
	var total int64

	query := r.db.Model(&domain.Donor{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("name ASC").Find(&donors).Error
	return donors, total, err
}

// --- Donation Repository Implementation ---

func (r *donationRepository) Create(donation *domain.Donation) error {
	// Transaction to save donation and items
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(donation).Error; err != nil {
			return err
		}

		// If there are items, they are automatically saved by GORM association usually,
		// but explicit check is good if we want custom logic.
		// Since we defined `Items []DonationItem` in domain, GORM handles it.

		return nil
	})
}

func (r *donationRepository) Update(donation *domain.Donation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update donation details
		if err := tx.Save(donation).Error; err != nil {
			return err
		}

		// Note: For items, we might need to verify if validation/replacement is needed.
		// For now, assuming domain logic handles item replacement/update logic before calling repo,
		// or GORM Association mode handles it.
		// But usually `Save` updates the main record.
		// If items are changed, we might need explicitly replace association.

		if len(donation.Items) > 0 {
			if err := tx.Model(donation).Association("Items").Replace(donation.Items); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *donationRepository) FindByID(id string) (*domain.Donation, error) {
	var donation domain.Donation
	err := r.db.Preload("Donor").Preload("Employee").Preload("Items").First(&donation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &donation, nil
}

func (r *donationRepository) FindAll(filter map[string]interface{}, limit, offset int) ([]domain.Donation, int64, error) {
	var donations []domain.Donation
	var total int64

	query := r.db.Model(&domain.Donation{}).Preload("Donor").Preload("Employee").Preload("Items")

	if val, ok := filter["date_from"]; ok {
		query = query.Where("date >= ?", val)
	}
	if val, ok := filter["date_to"]; ok {
		query = query.Where("date <= ?", val)
	}
	if val, ok := filter["type"]; ok && val != "" {
		query = query.Where("type = ?", val)
	}
	if val, ok := filter["donor_id"]; ok && val != "" {
		query = query.Where("donor_id = ?", val)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("date DESC").Limit(limit).Offset(offset).Find(&donations).Error
	return donations, total, err
}
