package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

type Student struct {
	ID                          string             `gorm:"primaryKey;type:char(36)" json:"id"`
	FullName                    string             `gorm:"type:varchar(100);not null" json:"full_name"`
	NoKK                        string             `gorm:"type:text" json:"no_kk,omitempty"`      // akan dienkripsi
	NIK                         *string            `gorm:"type:text" json:"nik,omitempty"`        // akan dienkripsi
	NIKHash                     *string            `gorm:"type:varchar(64);uniqueIndex" json:"-"` // Blind Index for Unique Check
	NISN                        *string            `gorm:"type:varchar(20);uniqueIndex" json:"nisn"`
	NIM                         *string            `gorm:"type:varchar(20);uniqueIndex" json:"nim"`
	Gender                      string             `gorm:"type:varchar(10)" json:"gender"`
	PlaceOfBirth                *string            `gorm:"type:varchar(100)" json:"place_of_birth"`
	DateOfBirth                 *utils.Date        `gorm:"type:date" json:"date_of_birth"`
	Address                     *string            `gorm:"type:text" json:"address"`
	RT                          *string            `gorm:"type:varchar(3)" json:"rt"`
	RW                          *string            `gorm:"type:varchar(3)" json:"rw"`
	SubDistrict                 *string            `gorm:"type:varchar(100)" json:"sub_district"`
	District                    *string            `gorm:"type:varchar(100)" json:"district"`
	City                        *string            `gorm:"type:varchar(100)" json:"city"`
	Province                    *string            `gorm:"type:varchar(100)" json:"province"`
	PostalCode                  *string            `gorm:"type:varchar(5)" json:"postal_code"`
	Status                      string             `gorm:"type:enum('ACTIVE','GRADUATED','DROPOUT');default:'ACTIVE'" json:"status"`
	EntryYear                   *string            `gorm:"type:varchar(4)" json:"entry_year"`
	ExitYear                    *string            `gorm:"type:varchar(4)" json:"exit_year"`
	BirthCertificateFile        *string            `gorm:"type:varchar(255)" json:"birth_certificate_file"`
	FamilyCardFile              *string            `gorm:"type:varchar(255)" json:"family_card_file"`
	ParentStatementFile         *string            `gorm:"type:varchar(255)" json:"parent_statement_file"`
	StudentStatementFile        *string            `gorm:"type:varchar(255)" json:"student_statement_file"`
	HealthInsuranceFile         *string            `gorm:"type:varchar(255)" json:"health_insurance_file"`
	DiplomaCertificateFile      *string            `gorm:"type:varchar(255)" json:"diploma_certificate_file"`
	GraduationCertificateFile   *string            `gorm:"type:varchar(255)" json:"graduation_certificate_file"`
	FinancialHardshipLetterFile *string            `gorm:"type:varchar(255)" json:"financial_hardship_letter_file"`
	CreatedAt                   time.Time          `json:"created_at"`
	UpdatedAt                   time.Time          `json:"updated_at"`
	Parents                     []StudentParent    `gorm:"foreignKey:StudentID" json:"parents,omitempty"`            // Relasi ke tabel pivot StudentParent
	StudentClassrooms           []StudentClassroom `gorm:"foreignKey:StudentID" json:"student_classrooms,omitempty"` // Relasi ke tabel pivot StudentClassroom
	GuardianID                  *string            `gorm:"type:char(36);index:idx_student_guardian" json:"guardian_id"`
	GuardianType                *string            `gorm:"type:varchar(20);index:idx_student_guardian" json:"guardian_type"`
	UserID                      *string            `gorm:"type:char(36);uniqueIndex" json:"user_id"`
	User                        User               `gorm:"foreignKey:UserID;references:ID"`
}

// Hook BeforeCreate untuk generate UUID
func (s *Student) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = utils.GenerateUUID()
	}
	return
}

func (s *Student) NISNValue() string {
	if s.NISN == nil {
		return "-"
	}
	return *s.NISN
}

func (s *Student) NIMValue() string {
	if s.NIM == nil {
		return "-"
	}
	return *s.NIM
}
