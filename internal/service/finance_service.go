package service

import (
	"encoding/json"
	"time"

	"smart_school_be/internal/apperrors"
	"smart_school_be/internal/model/domain"
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/model/response"
	"smart_school_be/internal/repository"
	"smart_school_be/internal/utils"
)

type FinanceService interface {
	CreateDonation(req request.CreateDonationRequest, employeeID string) (*response.DonationResponse, error)
	GetDonations(filter map[string]interface{}, pagination request.PaginationRequest) (*response.PaginatedData, error)
	GetDonationByID(id string) (*response.DonationResponse, error)
	UpdateDonation(id string, req request.UpdateDonationRequest) (*response.DonationResponse, error)
	GetDonors(name string, pagination request.PaginationRequest) (*response.PaginatedData, error)
	GetDonorByID(id string) (*response.DonorResponse, error)
	UpdateDonor(id string, req request.UpdateDonorRequest) (*response.DonorResponse, error)
}

type financeService struct {
	donorRepo    repository.DonorRepository
	donationRepo repository.DonationRepository
	employeeRepo repository.EmployeeRepository
	baseURL      string
}

func NewFinanceService(
	donorRepo repository.DonorRepository,
	donationRepo repository.DonationRepository,
	employeeRepo repository.EmployeeRepository,
	baseURL string,
) FinanceService {
	return &financeService{
		donorRepo:    donorRepo,
		donationRepo: donationRepo,
		employeeRepo: employeeRepo,
		baseURL:      baseURL,
	}
}

func (s *financeService) CreateDonation(req request.CreateDonationRequest, userID string) (*response.DonationResponse, error) {
	// 0. Resolve Employee form UserID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, apperrors.NewInternalError(err.Error())
	}
	if employee == nil {
		return nil, apperrors.NewUnauthorizedError("user is not associated with an employee record")
	}
	employeeID := employee.ID

	// 1. Get or Create Donor
	var donor *domain.Donor

	// Try to find by phone if provided
	if req.DonorPhone != "" {
		existingDonor, err := s.donorRepo.FindByPhone(req.DonorPhone)
		if err == nil {
			// Found by phone, check if name matches (case-insensitive or simple string match)
			// Ideally normalize names, but for now simple check
			if existingDonor.Name == req.DonorName {
				donor = existingDonor
			} else {
				// Phone matches but Name is different.
				// Treat as different person sharing phone (e.g. spouse).
				// Do NOT set donor = existingDonor, so it falls through to create new.
			}
		}
	}

	// parsing donor name if phone lookup fails OR name didn't match
	if donor == nil {
		donors, _, err := s.donorRepo.FindAll(req.DonorName, 1, 0)
		if err == nil && len(donors) > 0 {
			// Check for exact match on name
			for _, d := range donors {
				if d.Name == req.DonorName {
					// Check if this donor has same phone?
					// If req.DonorPhone is provided, we prefer strict match.
					// But if we are here, it means phone lookup failed or name didn't match phone owner.

					// If the existing donor has a generic nil phone, we might want to claim it?
					// For now, let's just reuse if name matches EXACTLY.
					donor = &d
					break
				}
			}
		}
	}

	if donor == nil {
		// Create new donor
		newDonor := &domain.Donor{
			Name:    req.DonorName,
			Phone:   utils.StringPtr(req.DonorPhone),
			Email:   utils.StringPtr(req.DonorEmail),
			Address: utils.StringPtr(req.DonorAddress),
		}
		// Helper to handle empty strings - set to nil if empty
		if req.DonorPhone == "" {
			newDonor.Phone = nil
		}
		if req.DonorEmail == "" {
			newDonor.Email = nil
		}
		if req.DonorAddress == "" {
			newDonor.Address = nil
		}

		if err := s.donorRepo.Create(newDonor); err != nil {
			return nil, err
		}
		donor = newDonor
	} else {
		// Update donor info if provided and different?
		// For simplicity, we just use existing donor.
		// Maybe update address if provided?
	}

	// 2. Handle File Upload (Already handled by handler, path provided in req)
	var proofFileURL *string
	if req.ProofFile != "" {
		proofFileURL = utils.StringPtr(req.ProofFile)
	}

	// 3. Parse Items
	var items []domain.DonationItem
	if req.ItemsJSON != "" {
		var itemReqs []request.DonationItemRequest
		if err := json.Unmarshal([]byte(req.ItemsJSON), &itemReqs); err != nil {
			return nil, apperrors.NewBadRequestError("invalid items_json format")
		}

		for _, ir := range itemReqs {
			items = append(items, domain.DonationItem{
				ItemName:       ir.ItemName,
				Quantity:       ir.Quantity,
				Unit:           ir.Unit,
				EstimatedValue: ir.EstimatedValue,
				Notes:          utils.StringPtr(ir.Notes),
			})
		}
	}

	// 4. Create Donation
	donation := &domain.Donation{
		DonorID:       donor.ID,
		EmployeeID:    employeeID,
		Date:          time.Now(),
		Type:          req.Type,
		PaymentMethod: req.PaymentMethod,
		TotalAmount:   req.TotalAmount,
		Description:   utils.StringPtr(req.Description),
		ProofFile:     proofFileURL,
		Items:         items,
	}

	if err := s.donationRepo.Create(donation); err != nil {
		return nil, err
	}

	// Reload to get full data (e.g. created_at, relationships)
	savedDonation, err := s.donationRepo.FindByID(donation.ID)
	if err != nil {
		return nil, err
	}

	// Map response (Log potential panic area)
	// log.Println("Mapping donation response...")
	return responsePtr(response.FromDomainDonation(savedDonation, s.baseURL)), nil
}

func (s *financeService) GetDonations(filter map[string]interface{}, pagination request.PaginationRequest) (*response.PaginatedData, error) {
	limit := pagination.GetLimit()
	offset := pagination.GetOffset()

	donations, total, err := s.donationRepo.FindAll(filter, limit, offset)
	if err != nil {
		return nil, err
	}

	var res []response.DonationResponse
	for _, d := range donations {
		res = append(res, response.FromDomainDonation(&d, s.baseURL))
	}

	paginatedData := response.NewPaginatedData(res, total, pagination.GetPage(), limit)
	return &paginatedData, nil
}

func (s *financeService) GetDonors(name string, pagination request.PaginationRequest) (*response.PaginatedData, error) {
	limit := pagination.GetLimit()
	offset := pagination.GetOffset()

	donors, total, err := s.donorRepo.FindAll(name, limit, offset)
	if err != nil {
		return nil, err
	}

	var res []response.DonorResponse
	for _, d := range donors {
		res = append(res, response.FromDomainDonor(&d))
	}

	paginatedData := response.NewPaginatedData(res, total, pagination.GetPage(), limit)
	return &paginatedData, nil
}

func (s *financeService) GetDonationByID(id string) (*response.DonationResponse, error) {
	donation, err := s.donationRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("donation not found")
	}
	return responsePtr(response.FromDomainDonation(donation, s.baseURL)), nil
}

func (s *financeService) UpdateDonation(id string, req request.UpdateDonationRequest) (*response.DonationResponse, error) {
	donation, err := s.donationRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("donation not found")
	}

	// Update fields if provided
	if req.Date != "" {
		// Try parsing various formats if needed, or stick to YYYY-MM-DD
		if t, err := time.Parse("2006-01-02", req.Date); err == nil {
			donation.Date = t
		}
	}
	if req.Type != "" {
		donation.Type = req.Type
	}
	if req.PaymentMethod != "" {
		donation.PaymentMethod = req.PaymentMethod
	}
	if req.TotalAmount > 0 {
		donation.TotalAmount = req.TotalAmount
	}
	// Allow clearing description? If so check intention. For now just update if not empty.
	if req.Description != "" {
		donation.Description = utils.StringPtr(req.Description)
	}

	// Handle Proof File upload if provided (path passed in req)
	if req.ProofFile != "" {
		donation.ProofFile = utils.StringPtr(req.ProofFile)
	}

	// Handle Items update for GOODS
	if req.Type == "GOODS" && req.ItemsJSON != "" {
		var requestItems []request.DonationItemRequest
		if err := json.Unmarshal([]byte(req.ItemsJSON), &requestItems); err != nil {
			return nil, apperrors.NewBadRequestError("invalid items_json format")
		}

		// Replace items
		var newItems []domain.DonationItem
		for _, item := range requestItems {
			newItems = append(newItems, domain.DonationItem{
				DonationID:     donation.ID,
				ItemName:       item.ItemName,
				Quantity:       item.Quantity,
				Unit:           item.Unit,
				EstimatedValue: item.EstimatedValue,
				Notes:          utils.StringPtr(item.Notes),
			})
		}
		donation.Items = newItems
	}

	if err := s.donationRepo.Update(donation); err != nil {
		return nil, err
	}

	return responsePtr(response.FromDomainDonation(donation, s.baseURL)), nil
}

func (s *financeService) GetDonorByID(id string) (*response.DonorResponse, error) {
	donor, err := s.donorRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("donor not found")
	}
	res := response.FromDomainDonor(donor)

	// Fetch recent 5 donations
	filter := map[string]interface{}{"donor_id": id}
	donations, _, err := s.donationRepo.FindAll(filter, 5, 0)
	if err == nil {
		// Map to response
		var recentDonations []response.DonationResponse
		for _, d := range donations {
			// To avoid cycle/redundancy, we could nil out the Donor field in recent donations if it wasn't a value type
			// But Since DonationResponse has Donor as value, we just accept it.
			// Or we can manually set it to empty DonorResponse if we want to save bandwidth,
			// but keeping it consistent is safer for now.
			recentDonations = append(recentDonations, response.FromDomainDonation(&d, s.baseURL))
		}
		res.RecentDonations = recentDonations
	}

	return &res, nil
}

func (s *financeService) UpdateDonor(id string, req request.UpdateDonorRequest) (*response.DonorResponse, error) {
	donor, err := s.donorRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("donor not found")
	}

	if req.Name != "" {
		donor.Name = req.Name
	}
	if req.Phone != "" {
		donor.Phone = utils.StringPtr(req.Phone)
	}
	if req.Email != "" {
		donor.Email = utils.StringPtr(req.Email)
	}
	if req.Address != "" {
		donor.Address = utils.StringPtr(req.Address)
	}

	if err := s.donorRepo.Update(donor); err != nil {
		return nil, err
	}

	res := response.FromDomainDonor(donor)
	return &res, nil
}

func responsePtr(r response.DonationResponse) *response.DonationResponse {
	return &r
}
