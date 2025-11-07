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
