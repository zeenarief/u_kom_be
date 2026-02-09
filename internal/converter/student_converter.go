package converter

import (
	"log"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/utils"
)

// StudentConverterInterface mendefinisikan kontrak untuk konverter
type StudentConverterInterface interface {
	ToStudentDetailResponse(student *domain.Student) *response.StudentDetailResponse
	ToStudentListResponse(student *domain.Student) *response.StudentListResponse
	ToStudentListResponses(students []domain.Student) []response.StudentListResponse
}

// studentConverter adalah implementasi dengan dependensi
type studentConverter struct {
	encryptionUtil  utils.EncryptionUtil
	parentConverter ParentConverterInterface
}

// NewStudentConverter membuat instance konverter baru
func NewStudentConverter(
	encryptionUtil utils.EncryptionUtil,
	parentConverter ParentConverterInterface,
) StudentConverterInterface {
	return &studentConverter{
		encryptionUtil:  encryptionUtil,
		parentConverter: parentConverter,
	}
}

// ToStudentDetailResponse mengubah domain Student (terenkripsi) ke response (plaintext)
func (c *studentConverter) ToStudentDetailResponse(student *domain.Student) *response.StudentDetailResponse {
	// Dekripsi NIK
	decryptedNIK := ""
	if student.NIK != nil {
		var err error
		decryptedNIK, err = c.encryptionUtil.Decrypt(*student.NIK)
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

	var parentResponses []response.ParentRelationshipResponse
	// Cek apakah relasi Parents di-load (tidak nil)
	if student.Parents != nil {
		for _, sp := range student.Parents {
			// sp.Parent akan terisi jika kita Preload("Parents.Parent")
			if sp.Parent.ID != "" { // Pastikan data parent ada
				parentResponses = append(parentResponses, response.ParentRelationshipResponse{
					RelationshipType: sp.RelationshipType,
					// Kita gunakan parentConverter untuk mengubah domain.Parent -> response.ParentListResponse
					Parent: *c.parentConverter.ToParentListResponse(&sp.Parent),
				})
			}
		}
	}

	return &response.StudentDetailResponse{
		ID:             student.ID,
		FullName:       student.FullName,
		NoKK:           decryptedNoKK, // <-- Data plaintext
		NIK:            decryptedNIK,  // <-- Data plaintext
		NISN:           student.NISN,
		NIM:            student.NIM,
		Gender:         student.Gender,
		PlaceOfBirth:   student.PlaceOfBirth,
		DateOfBirth:    student.DateOfBirth,
		Address:        student.Address,
		RT:             student.RT,
		RW:             student.RW,
		SubDistrict:    student.SubDistrict,
		District:       student.District,
		City:           student.City,
		Province:       student.Province,
		PostalCode:     student.PostalCode,
		Status:         student.Status,
		EntryYear:      student.EntryYear,
		GraduationYear: student.GraduationYear,
		CreatedAt:      student.CreatedAt,
		UpdatedAt:      student.UpdatedAt,
		Parents:        parentResponses,
	}
}

// ToStudentListResponse mengubah domain ke response list (ringkas)
func (c *studentConverter) ToStudentListResponse(student *domain.Student) *response.StudentListResponse {
	// Cari kelas aktif
	var className, major, level string

	// Gunakan status dari tabel students sebagai default (ACTIVE, GRADUATED, dll)
	status := student.Status

	if len(student.StudentClassrooms) > 0 {
		// Karena di repository kita filter "status = ACTIVE", maka yg muncul di sini harusnya aktif semua.
		// Kita ambil yg pertama.
		activeClass := student.StudentClassrooms[0]
		className = activeClass.Classroom.Name
		major = activeClass.Classroom.Major
		level = activeClass.Classroom.Level
		// Optional: Anda bisa gunakan activeClass.Status jika ingin override dengan status classroom
		// status = activeClass.Status
	}

	email := ""
	if student.User.ID != "" {
		email = student.User.Email
	}

	return &response.StudentListResponse{
		ID:        student.ID,
		FullName:  student.FullName,
		NISN:      student.NISN,
		NIM:       student.NIM,
		Gender:    student.Gender,
		District:  student.District,
		City:      student.City,
		ClassName: className,
		Major:     major,
		Level:     level,
		Status:    status,
		Email:     email,
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
