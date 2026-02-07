package converter

import (
	"log"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/utils"
)

// ParentConverterInterface mendefinisikan kontrak untuk konverter Parent
type ParentConverterInterface interface {
	ToParentDetailResponse(parent *domain.Parent) *response.ParentDetailResponse
	ToParentListResponse(parent *domain.Parent) *response.ParentListResponse
	ToParentListResponses(parents []domain.Parent) []response.ParentListResponse
}

// parentConverter adalah implementasi dengan dependensi
type parentConverter struct {
	encryptionUtil utils.EncryptionUtil
}

// NewParentConverter membuat instance konverter baru
func NewParentConverter(encryptionUtil utils.EncryptionUtil) ParentConverterInterface {
	return &parentConverter{
		encryptionUtil: encryptionUtil,
	}
}

// ToParentDetailResponse mengubah domain Parent (terenkripsi) ke response detail (plaintext)
func (c *parentConverter) ToParentDetailResponse(parent *domain.Parent) *response.ParentDetailResponse {
	// Dekripsi NIK
	var decryptedNIK *string
	if parent.NIK != nil {
		decrypted, err := c.encryptionUtil.Decrypt(*parent.NIK)
		if err != nil {
			log.Printf("Failed to decrypt NIK for parent %s: %v", parent.ID, err)
			errStr := "[DECRYPTION_ERROR]"
			decryptedNIK = &errStr
		} else {
			decryptedNIK = &decrypted
		}
	}

	return &response.ParentDetailResponse{
		ID:             parent.ID,
		FullName:       parent.FullName,
		NIK:            decryptedNIK, // <-- Data plaintext *string
		Gender:         parent.Gender,
		PlaceOfBirth:   parent.PlaceOfBirth,
		DateOfBirth:    parent.DateOfBirth,
		LifeStatus:     parent.LifeStatus,
		MaritalStatus:  parent.MaritalStatus,
		PhoneNumber:    parent.PhoneNumber,
		Email:          parent.Email,
		EducationLevel: parent.EducationLevel,
		Occupation:     parent.Occupation,
		IncomeRange:    parent.IncomeRange,
		Address:        parent.Address,
		RT:             parent.RT,
		RW:             parent.RW,
		SubDistrict:    parent.SubDistrict,
		District:       parent.District,
		City:           parent.City,
		Province:       parent.Province,
		PostalCode:     parent.PostalCode,
		CreatedAt:      parent.CreatedAt,
		UpdatedAt:      parent.UpdatedAt,
	}
}

// ToParentListResponse mengubah domain ke response list (ringkas)
func (c *parentConverter) ToParentListResponse(parent *domain.Parent) *response.ParentListResponse {
	return &response.ParentListResponse{
		ID:          parent.ID,
		FullName:    parent.FullName,
		Gender:      parent.Gender,
		LifeStatus:  parent.LifeStatus,
		PhoneNumber: parent.PhoneNumber,
		Email:       parent.Email,
		Occupation:  parent.Occupation,
	}
}

// ToParentListResponses adalah helper untuk list
func (c *parentConverter) ToParentListResponses(parents []domain.Parent) []response.ParentListResponse {
	var responses []response.ParentListResponse
	for _, p := range parents {
		responses = append(responses, *c.ToParentListResponse(&p))
	}
	return responses
}
