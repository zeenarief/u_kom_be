package converter

import (
	"belajar-golang/internal/model/domain"
	"belajar-golang/internal/model/response"
	"belajar-golang/internal/utils"
	"log"
)

// StudentConverterInterface mendefinisikan kontrak untuk konverter
type StudentConverterInterface interface {
	ToStudentDetailResponse(student *domain.Student) *response.StudentDetailResponse
	ToStudentListResponse(student *domain.Student) *response.StudentListResponse
	ToStudentListResponses(students []domain.Student) []response.StudentListResponse
}

// studentConverter adalah implementasi dengan dependensi
type studentConverter struct {
	encryptionUtil utils.EncryptionUtil
}

// NewStudentConverter membuat instance konverter baru
func NewStudentConverter(encryptionUtil utils.EncryptionUtil) StudentConverterInterface {
	return &studentConverter{
		encryptionUtil: encryptionUtil,
	}
}

// ToStudentDetailResponse mengubah domain Student (terenkripsi) ke response (plaintext)
func (c *studentConverter) ToStudentDetailResponse(student *domain.Student) *response.StudentDetailResponse {
	// Dekripsi NIK
	decryptedNIK := ""
	if student.NIK != "" {
		var err error
		decryptedNIK, err = c.encryptionUtil.Decrypt(student.NIK)
		if err != nil {
			log.Printf("Failed to decrypt NIK for student %s: %v", student.ID, err)
			decryptedNIK = "[DECRYPTION_ERROR]"
		}
	}

	// Dekripsi No. KK
	decryptedNoKK := ""
	if student.NoKK != "" {
		var err error
		decryptedNoKK, err = c.encryptionUtil.Decrypt(student.NoKK)
		if err != nil {
			log.Printf("Failed to decrypt NoKK for student %s: %v", student.ID, err)
			decryptedNoKK = "[DECRYPTION_ERROR]"
		}
	}

	return &response.StudentDetailResponse{
		ID:           student.ID,
		FullName:     student.FullName,
		NoKK:         decryptedNoKK, // <-- Data plaintext
		NIK:          decryptedNIK,  // <-- Data plaintext
		NISN:         student.NISN,
		NIM:          student.NIM,
		Gender:       student.Gender,
		PlaceOfBirth: student.PlaceOfBirth,
		DateOfBirth:  student.DateOfBirth,
		Address:      student.Address,
		RT:           student.RT,
		RW:           student.RW,
		SubDistrict:  student.SubDistrict,
		District:     student.District,
		City:         student.City,
		Province:     student.Province,
		PostalCode:   student.PostalCode,
		CreatedAt:    student.CreatedAt,
		UpdatedAt:    student.UpdatedAt,
	}
}

// ToStudentListResponse mengubah domain ke response list (ringkas)
func (c *studentConverter) ToStudentListResponse(student *domain.Student) *response.StudentListResponse {
	return &response.StudentListResponse{
		ID:       student.ID,
		FullName: student.FullName,
		NISN:     student.NISN,
		NIM:      student.NIM,
		Gender:   student.Gender,
		City:     student.City,
	}
}

// ToStudentListResponses adalah helper untuk list (menggantikan ToStudentDetailResponses)
func (c *studentConverter) ToStudentListResponses(students []domain.Student) []response.StudentListResponse {
	var responses []response.StudentListResponse // <-- Tipe diubah
	for _, s := range students {
		responses = append(responses, *c.ToStudentListResponse(&s)) // <-- Memanggil ToStudentListResponse
	}
	return responses
}
