package handler

import (
	"net/http"

	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"
	"smart_school_be/internal/utils"

	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	financeService service.FinanceService
}

func NewFinanceHandler(financeService service.FinanceService) *FinanceHandler {
	return &FinanceHandler{financeService: financeService}
}

func (h *FinanceHandler) CreateDonation(c *gin.Context) {
	var req request.CreateDonationRequest

	// Bind multipart form
	if err := c.ShouldBind(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// Handle file upload
	file, err := c.FormFile("proof_file")
	if err == nil {
		path, errSave := utils.SaveUploadedFile(c, file, "donations", "donation_proof")
		if errSave != nil {
			BadRequestError(c, "Failed to save proof file", errSave.Error())
			return
		}
		// Set the path in request struct
		req.ProofFile = path
	} else if err != http.ErrMissingFile {
		BadRequestError(c, "File upload error", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		UnauthorizedError(c, "User not authenticated")
		return
	}

	res, err := h.financeService.CreateDonation(req, userID.(string))
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Donation created successfully", res)
}

func (h *FinanceHandler) GetDonations(c *gin.Context) {
	pagination := request.NewPaginationRequest(c.Query("page"), c.Query("limit"))

	filter := make(map[string]interface{})
	if v := c.Query("type"); v != "" {
		filter["type"] = v
	}
	if v := c.Query("donor_id"); v != "" {
		filter["donor_id"] = v
	}
	if v := c.Query("date_from"); v != "" {
		filter["date_from"] = v
	}
	if v := c.Query("date_to"); v != "" {
		filter["date_to"] = v
	}

	res, err := h.financeService.GetDonations(filter, pagination)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Donations retrieved successfully", res)
}

func (h *FinanceHandler) GetDonors(c *gin.Context) {
	pagination := request.NewPaginationRequest(c.Query("page"), c.Query("limit"))

	name := c.Query("name")

	res, err := h.financeService.GetDonors(name, pagination)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Donors retrieved successfully", res)
}

func (h *FinanceHandler) GetDonationByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.financeService.GetDonationByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Donation retrieved successfully", res)
}

func (h *FinanceHandler) UpdateDonation(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateDonationRequest

	// Bind multipart form
	if err := c.ShouldBind(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// Handle file upload if present
	file, err := c.FormFile("proof_file")
	if err == nil {
		path, errSave := utils.SaveUploadedFile(c, file, "donations", "donation_proof")
		if errSave != nil {
			BadRequestError(c, "Failed to save proof file", errSave.Error())
			return
		}
		req.ProofFile = path
	}

	res, err := h.financeService.UpdateDonation(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Donation updated successfully", res)
}

func (h *FinanceHandler) GetDonorByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.financeService.GetDonorByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Donor retrieved successfully", res)
}

func (h *FinanceHandler) UpdateDonor(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateDonorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.financeService.UpdateDonor(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Donor updated successfully", res)
}
