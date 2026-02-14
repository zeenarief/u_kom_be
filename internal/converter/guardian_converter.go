package converter

import (
	"log"
	"smart_school_be/internal/model/domain"
	"smart_school_be/internal/model/response"
	"smart_school_be/internal/utils"
)

// GuardianConverterInterface mendefinisikan kontrak untuk konverter Guardian
type GuardianConverterInterface interface {
	ToGuardianDetailResponse(guardian *domain.Guardian) *response.GuardianDetailResponse
	ToGuardianListResponse(guardian *domain.Guardian) *response.GuardianListResponse
	ToGuardianListResponses(guardians []domain.Guardian) []response.GuardianListResponse
}

// guardianConverter adalah implementasi dengan dependensi
type guardianConverter struct {
	encryptionUtil utils.EncryptionUtil
}

// NewGuardianConverter membuat instance konverter baru
func NewGuardianConverter(encryptionUtil utils.EncryptionUtil) GuardianConverterInterface {
	return &guardianConverter{
		encryptionUtil: encryptionUtil,
	}
}

// ToGuardianDetailResponse mengubah domain Guardian (terenkripsi) ke response detail (plaintext)
func (c *guardianConverter) ToGuardianDetailResponse(guardian *domain.Guardian) *response.GuardianDetailResponse {
	// Dekripsi NIK
	var decryptedNIK *string
	if guardian.NIK != nil {
		decrypted, err := c.encryptionUtil.Decrypt(*guardian.NIK)
		if err != nil {
			log.Printf("Failed to decrypt NIK for guardian %s: %v", guardian.ID, err)
			errStr := "[DECRYPTION_ERROR]"
			decryptedNIK = &errStr
		} else {
			decryptedNIK = &decrypted
		}
	}

	return &response.GuardianDetailResponse{
		ID:                    guardian.ID,
		FullName:              guardian.FullName,
		NIK:                   decryptedNIK, // <-- Data plaintext *string
		Gender:                guardian.Gender,
		PhoneNumber:           guardian.PhoneNumber,
		Email:                 guardian.Email,
		Address:               guardian.Address,
		RT:                    guardian.RT,
		RW:                    guardian.RW,
		SubDistrict:           guardian.SubDistrict,
		District:              guardian.District,
		City:                  guardian.City,
		Province:              guardian.Province,
		PostalCode:            guardian.PostalCode,
		RelationshipToStudent: guardian.RelationshipToStudent,
		CreatedAt:             guardian.CreatedAt,
		UpdatedAt:             guardian.UpdatedAt,
	}
}

// ToGuardianListResponse mengubah domain ke response list (ringkas)
func (c *guardianConverter) ToGuardianListResponse(guardian *domain.Guardian) *response.GuardianListResponse {
	return &response.GuardianListResponse{
		ID:                    guardian.ID,
		FullName:              guardian.FullName,
		PhoneNumber:           guardian.PhoneNumber,
		Email:                 guardian.Email,
		RelationshipToStudent: guardian.RelationshipToStudent,
	}
}

// ToGuardianListResponses adalah helper untuk list
func (c *guardianConverter) ToGuardianListResponses(guardians []domain.Guardian) []response.GuardianListResponse {
	var responses []response.GuardianListResponse
	for _, g := range guardians {
		responses = append(responses, *c.ToGuardianListResponse(&g))
	}
	return responses
}
