package handler

import (
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/service"
	"strings"

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
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	student, err := h.studentService.CreateStudent(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			BadRequestError(c, "Student creation failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	CreatedResponse(c, "Student created successfully", student)
}

func (h *StudentHandler) GetAllStudents(c *gin.Context) {
	students, err := h.studentService.GetAllStudents()
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Students retrieved successfully", students)
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id := c.Param("id")

	student, err := h.studentService.GetStudentByID(id)
	if err != nil {
		if err.Error() == "student not found" {
			NotFoundError(c, "Student not found")
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Student retrieved successfully", student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")

	var req request.StudentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	student, err := h.studentService.UpdateStudent(id, req)
	if err != nil {
		if err.Error() == "student not found" {
			NotFoundError(c, "Student not found")
		} else if strings.Contains(err.Error(), "already exists") {
			BadRequestError(c, "Student update failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Student updated successfully", student)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")

	err := h.studentService.DeleteStudent(id)
	if err != nil {
		if err.Error() == "student not found" {
			NotFoundError(c, "Student not found")
		} else {
			InternalServerError(c, err.Error())
		}
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
		// Tangani error spesifik dari service
		if strings.Contains(err.Error(), "student not found") {
			NotFoundError(c, err.Error())
		} else if strings.Contains(err.Error(), "parent not found") {
			BadRequestError(c, "Invalid parent ID", err.Error())
		} else if strings.Contains(err.Error(), "duplicate parent_id") {
			BadRequestError(c, "Invalid request", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
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
		// Tangani error spesifik dari service
		if strings.Contains(err.Error(), "student not found") {
			NotFoundError(c, err.Error())
		} else if strings.Contains(err.Error(), "parent not found") || strings.Contains(err.Error(), "guardian not found") {
			BadRequestError(c, "Invalid guardian_id", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
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
		if strings.Contains(err.Error(), "student not found") {
			NotFoundError(c, err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Student guardian removed successfully", nil)
}
