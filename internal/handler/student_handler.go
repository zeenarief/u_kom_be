package handler

import (
	"fmt"
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"
	"smart_school_be/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentService service.StudentService
}

func NewStudentHandler(studentService service.StudentService) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req request.StudentCreateRequest
	// Ganti ShouldBindJSON ke ShouldBind agar support multipart/form-data
	if err := c.ShouldBind(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// List file yang akan diupload
	fileKeys := []struct {
		Key      string
		FieldPtr *string
	}{
		{"birth_certificate_file", nil},
		{"family_card_file", nil},
		{"parent_statement_file", nil},
		{"student_statement_file", nil},
		{"health_insurance_file", nil},
		{"diploma_certificate_file", nil},
		{"graduation_certificate_file", nil},
		{"financial_hardship_letter_file", nil},
	}

	uploadedPaths := make(map[string]string)

	// Clean up jika terjadi error di tengah jalan
	defer func() {
		if c.IsAborted() {
			for _, path := range uploadedPaths {
				utils.RemoveFile(path)
			}
		}
	}()

	for _, fk := range fileKeys {
		file, err := c.FormFile(fk.Key)
		if err == nil {
			// Simpan file
			path, errSave := utils.SaveUploadedFile(c, file, "students", fk.Key)
			if errSave != nil {
				// Hapus file yang sudah terlanjur terupload
				for _, p := range uploadedPaths {
					utils.RemoveFile(p)
				}
				BadRequestError(c, fmt.Sprintf("Failed to upload %s", fk.Key), errSave.Error())
				return
			}
			uploadedPaths[fk.Key] = path
		}
	}

	// Masukkan ke struct
	filesToCheck := service.StudentFiles{
		BirthCertificateFile:        uploadedPaths["birth_certificate_file"],
		FamilyCardFile:              uploadedPaths["family_card_file"],
		ParentStatementFile:         uploadedPaths["parent_statement_file"],
		StudentStatementFile:        uploadedPaths["student_statement_file"],
		HealthInsuranceFile:         uploadedPaths["health_insurance_file"],
		DiplomaCertificateFile:      uploadedPaths["diploma_certificate_file"],
		GraduationCertificateFile:   uploadedPaths["graduation_certificate_file"],
		FinancialHardshipLetterFile: uploadedPaths["financial_hardship_letter_file"],
	}

	student, err := h.studentService.CreateStudent(req, filesToCheck)
	if err != nil {
		HandleError(c, err)
		return
	}

	CreatedResponse(c, "Student created successfully", student)
}

func (h *StudentHandler) GetAllStudents(c *gin.Context) {
	searchQuery := c.Query("q")
	classroomID := c.Query("classroom_id")
	pagination := request.NewPaginationRequest(c.Query("page"), c.Query("limit"))

	students, err := h.studentService.GetAllStudents(searchQuery, classroomID, pagination)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Students retrieved successfully", students)
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id := c.Param("id")

	student, err := h.studentService.GetStudentByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Student retrieved successfully", student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")

	var req request.StudentUpdateRequest
	// Ganti ShouldBindJSON ke ShouldBind
	if err := c.ShouldBind(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// File upload loop
	fileKeys := []string{
		"birth_certificate_file",
		"family_card_file",
		"parent_statement_file",
		"student_statement_file",
		"health_insurance_file",
		"diploma_certificate_file",
		"graduation_certificate_file",
		"financial_hardship_letter_file",
	}

	uploadedPaths := make(map[string]string)

	for _, key := range fileKeys {
		file, err := c.FormFile(key)
		if err == nil {
			path, errSave := utils.SaveUploadedFile(c, file, "students", key)
			if errSave != nil {
				// Cleanup newly uploaded files
				for _, p := range uploadedPaths {
					utils.RemoveFile(p)
				}
				BadRequestError(c, fmt.Sprintf("Failed to upload %s", key), errSave.Error())
				return
			}
			uploadedPaths[key] = path
		}
	}

	filesToUpdate := service.StudentFiles{
		BirthCertificateFile:        uploadedPaths["birth_certificate_file"],
		FamilyCardFile:              uploadedPaths["family_card_file"],
		ParentStatementFile:         uploadedPaths["parent_statement_file"],
		StudentStatementFile:        uploadedPaths["student_statement_file"],
		HealthInsuranceFile:         uploadedPaths["health_insurance_file"],
		DiplomaCertificateFile:      uploadedPaths["diploma_certificate_file"],
		GraduationCertificateFile:   uploadedPaths["graduation_certificate_file"],
		FinancialHardshipLetterFile: uploadedPaths["financial_hardship_letter_file"],
	}

	student, err := h.studentService.UpdateStudent(id, req, filesToUpdate)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Student updated successfully", student)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")

	err := h.studentService.DeleteStudent(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Student deleted successfully", nil)
}

// SyncParents adalah handler untuk POST /students/:id/sync-parents
func (h *StudentHandler) SyncParents(c *gin.Context) {
	id := c.Param("id")

	var req request.StudentSyncParentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// Panggil service
	err := h.studentService.SyncParents(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Student parents synced successfully", nil)
}

// SetGuardian adalah handler untuk PUT /students/:id/set-guardian
func (h *StudentHandler) SetGuardian(c *gin.Context) {
	id := c.Param("id")

	var req request.StudentSetGuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Validasi 'oneof' gagal akan ditangkap di sini
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// Panggil service
	err := h.studentService.SetGuardian(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Student guardian set successfully", nil)
}

// RemoveGuardian adalah handler untuk DELETE /students/:id/remove-guardian
func (h *StudentHandler) RemoveGuardian(c *gin.Context) {
	id := c.Param("id")

	// Panggil service
	err := h.studentService.RemoveGuardian(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Student guardian removed successfully", nil)
}

// LinkUser menangani POST /students/:id/link-user
func (h *StudentHandler) LinkUser(c *gin.Context) {
	studentID := c.Param("id")

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload (missing 'user_id')", err.Error())
		return
	}

	err := h.studentService.LinkUser(studentID, req.UserID)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Student linked to user successfully", nil)
}

// UnlinkUser menangani DELETE /students/:id/unlink-user
func (h *StudentHandler) UnlinkUser(c *gin.Context) {
	studentID := c.Param("id")

	err := h.studentService.UnlinkUser(studentID)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Student unlinked from user successfully", nil)
}

func (h *StudentHandler) ExportExcel(c *gin.Context) {
	buffer, err := h.studentService.ExportStudentsToExcel()
	if err != nil {
		InternalServerError(c, "Failed to generate excel file")
		return
	}

	// Nama file dinamis dengan timestamp
	filename := fmt.Sprintf("data_siswa_%s.xlsx", time.Now().Format("20060102_150405"))

	// Set Headers untuk Download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Kirim binary data
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}

func (h *StudentHandler) ExportPDF(c *gin.Context) {
	buffer, err := h.studentService.ExportStudentsToPdf()
	if err != nil {
		InternalServerError(c, "Failed to generate PDF file")
		return
	}

	filename := fmt.Sprintf("data_siswa_%s.pdf", time.Now().Format("20060102_150405"))

	// Header untuk PDF
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", buffer.Len()))

	c.Data(200, "application/pdf", buffer.Bytes())
}

func (h *StudentHandler) ExportStudentBiodata(c *gin.Context) {
	id := c.Param("id") // Ambil ID dari URL

	buffer, err := h.studentService.ExportStudentBiodata(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	filename := fmt.Sprintf("biodata_%s.pdf", id)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Data(200, "application/pdf", buffer.Bytes())
}
