package domain

import (
	"time"

	"smart_school_be/internal/utils"

	"gorm.io/gorm"
)

// Donor represents a person or entity making a donation
type Donor struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;index" json:"name"`
	Phone     *string   `gorm:"type:varchar(20);index" json:"phone"`
	Email     *string   `gorm:"type:varchar(255)" json:"email"`
	Address   *string   `gorm:"type:text" json:"address"`
	Notes     *string   `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Donor) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = utils.GenerateUUID()
	}
	return
}

func (d *Donor) TableName() string {
	return "finance_donors"
}

// Donation represents the transaction record
type Donation struct {
	ID            string    `gorm:"primaryKey;type:char(36)" json:"id"`
	DonorID       string    `gorm:"type:char(36);not null;index" json:"donor_id"`
	EmployeeID    string    `gorm:"type:char(36);not null;index" json:"employee_id"` // Receiver
	Date          time.Time `gorm:"type:datetime(3);not null;index" json:"date"`
	Type          string    `gorm:"type:enum('MONEY', 'GOODS', 'MIXED');not null" json:"type"`
	PaymentMethod string    `gorm:"type:enum('CASH', 'TRANSFER', 'QRIS', 'GOODS');not null" json:"payment_method"`
	TotalAmount   float64   `gorm:"type:decimal(15,2);default:0" json:"total_amount"`
	ProofFile     *string   `gorm:"type:varchar(255)" json:"proof_file"`
	Description   *string   `gorm:"type:text" json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Donor    Donor          `gorm:"foreignKey:DonorID" json:"donor"`
	Employee Employee       `gorm:"foreignKey:EmployeeID" json:"employee"`
	Items    []DonationItem `gorm:"foreignKey:DonationID" json:"items,omitempty"`
}

func (d *Donation) TableName() string {
	return "finance_donations"
}

func (d *Donation) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = utils.GenerateUUID()
	}
	return
}

// DonationItem represents goods details
type DonationItem struct {
	ID             string    `gorm:"primaryKey;type:char(36)" json:"id"`
	DonationID     string    `gorm:"type:char(36);not null;index" json:"donation_id"`
	ItemName       string    `gorm:"type:varchar(255);not null" json:"item_name"`
	Quantity       float64   `gorm:"type:decimal(10,2);not null;default:1" json:"quantity"`
	Unit           string    `gorm:"type:varchar(50);not null" json:"unit"`
	EstimatedValue float64   `gorm:"type:decimal(15,2);default:0" json:"estimated_value"`
	Notes          *string   `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (di *DonationItem) TableName() string {
	return "finance_donation_items"
}

func (di *DonationItem) BeforeCreate(tx *gorm.DB) (err error) {
	if di.ID == "" {
		di.ID = utils.GenerateUUID()
	}
	return
}
