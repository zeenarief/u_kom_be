package handler

import (
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service service.AttendanceService
}

func NewAttendanceHandler(service service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) Submit(c *gin.Context) {
	var req request.AttendanceSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.SubmitAttendance(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	CreatedResponse(c, "Attendance submitted successfully", res)
}

func (h *AttendanceHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	res, err := h.service.GetSessionDetail(id)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Attendance detail retrieved", res)
}

func (h *AttendanceHandler) GetHistory(c *gin.Context) {
	// Ambil user yang sedang login dari context
	// Asumsi middleware auth menyimpan object user atau claims
	_, exists := c.Get("user")
	if !exists {
		UnauthorizedError(c, "User context not found")
		return
	}

	// currentUser := user.(*domain.User)

	// Kita perlu mencari Employee ID berdasarkan User ID yang login
	// (Fitur ini membutuhkan relasi User -> Employee sudah ter-setup dengan benar)
	// Jika User adalah Admin, mungkin dia kirim ?teacher_id=...
	// Jika User adalah Guru, kita pakai ID dia sendiri.

	// SEMENTARA: Kita ambil dari Query Param teacher_id untuk fleksibilitas Admin
	teacherID := c.Query("teacher_id")

	// Jika kosong, coba cek apakah user ini punya relasi employee (Logic ini ada di UserService biasanya)
	// Untuk keamanan MVP, kita wajibkan query param dulu atau validasi logic terpisah.
	if teacherID == "" {
		BadRequestError(c, "teacher_id query parameter required", nil)
		return
	}

	res, err := h.service.GetHistoryByTeacher(teacherID)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Attendance history retrieved", res)
}

func (h *AttendanceHandler) GetHistoryByAssignment(c *gin.Context) {
	taID := c.Param("id")
	if taID == "" {
		BadRequestError(c, "teaching_assignment_id required", nil)
		return
	}

	res, err := h.service.GetHistoryByAssignment(taID)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Attendance history retrieved", res)
}

func (h *AttendanceHandler) CheckSession(c *gin.Context) {
	scheduleID := c.Query("schedule_id")
	date := c.Query("date") // YYYY-MM-DD

	if scheduleID == "" || date == "" {
		BadRequestError(c, "schedule_id and date required", nil)
		return
	}

	// Gunakan method baru yang cerdas: Return Session Existing ATAU List Siswa jika belum ada
	res, err := h.service.GetSessionOrClassList(scheduleID, date)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Session check / initiation", res)
}

func (h *AttendanceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteSession(id); err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Attendance session deleted successfully", nil)
}
