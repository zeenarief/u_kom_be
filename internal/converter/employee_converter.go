package converter

import (
	"belajar-golang/internal/model/domain"
	"belajar-golang/internal/model/response"
	"belajar-golang/internal/utils"
	"log"
)

// EmployeeConverterInterface mendefinisikan kontrak untuk konverter Employee
type EmployeeConverterInterface interface {
	ToEmployeeDetailResponse(employee *domain.Employee) *response.EmployeeDetailResponse
	ToEmployeeListResponse(employee *domain.Employee) *response.EmployeeListResponse
	ToEmployeeListResponses(employees []domain.Employee) []response.EmployeeListResponse
}

// employeeConverter adalah implementasi dengan dependensi
type employeeConverter struct {
	encryptionUtil utils.EncryptionUtil
}

// NewEmployeeConverter membuat instance konverter baru
func NewEmployeeConverter(encryptionUtil utils.EncryptionUtil) EmployeeConverterInterface {
	return &employeeConverter{
		encryptionUtil: encryptionUtil,
	}
}

// ToEmployeeDetailResponse mengubah domain Employee (terenkripsi) ke response detail (plaintext)
func (c *employeeConverter) ToEmployeeDetailResponse(employee *domain.Employee) *response.EmployeeDetailResponse {
	// Dekripsi NIK
	decryptedNIK := ""
	if employee.NIK != "" {
		var err error
		decryptedNIK, err = c.encryptionUtil.Decrypt(employee.NIK)
		if err != nil {
			log.Printf("Failed to decrypt NIK for employee %s: %v", employee.ID, err)
			decryptedNIK = "[DECRYPTION_ERROR]"
		}
	}

	return &response.EmployeeDetailResponse{
		ID:               employee.ID,
		FullName:         employee.FullName,
		NIP:              employee.NIP,
		JobTitle:         employee.JobTitle,
		NIK:              decryptedNIK, // <-- Data plaintext
		Gender:           employee.Gender,
		PhoneNumber:      employee.PhoneNumber,
		Address:          employee.Address,
		DateOfBirth:      employee.DateOfBirth,
		JoinDate:         employee.JoinDate,
		EmploymentStatus: employee.EmploymentStatus,
		CreatedAt:        employee.CreatedAt,
		UpdatedAt:        employee.UpdatedAt,
	}
}

// ToEmployeeListResponse mengubah domain ke response list (ringkas)
func (c *employeeConverter) ToEmployeeListResponse(employee *domain.Employee) *response.EmployeeListResponse {
	return &response.EmployeeListResponse{
		ID:               employee.ID,
		FullName:         employee.FullName,
		NIP:              employee.NIP,
		JobTitle:         employee.JobTitle,
		PhoneNumber:      employee.PhoneNumber,
		EmploymentStatus: employee.EmploymentStatus,
	}
}

// ToEmployeeListResponses adalah helper untuk list
func (c *employeeConverter) ToEmployeeListResponses(employees []domain.Employee) []response.EmployeeListResponse {
	var responses []response.EmployeeListResponse
	for _, e := range employees {
		responses = append(responses, *c.ToEmployeeListResponse(&e))
	}
	return responses
}
