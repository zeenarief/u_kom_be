package domain

import "time"

// StudentParent adalah model untuk tabel pivot student_parent
// Ini BUKAN model untuk CRUD, tapi untuk relasi GORM
type StudentParent struct {
	StudentID        string    `gorm:"primaryKey;type:char(36)"`
	ParentID         string    `gorm:"primaryKey;type:char(36)"`
	RelationshipType string    `gorm:"type:varchar(50);not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`

	// Definisi relasi (agar kita bisa Preload)
	Student Student `gorm:"foreignKey:StudentID;references:ID"`
	Parent  Parent  `gorm:"foreignKey:ParentID;references:ID"`
}

// Opsional: Tentukan nama tabel secara eksplisit
func (StudentParent) TableName() string {
	return "student_parent"
}
