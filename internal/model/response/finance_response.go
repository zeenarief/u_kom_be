package response

import (
	"fmt"
	"smart_school_be/internal/model/domain"
	"time"
)

type DonorResponse struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Phone           string             `json:"phone"`           // Always return
	Email           string             `json:"email,omitempty"` // Keep optional
	Address         string             `json:"address"`         // Always return
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	RecentDonations []DonationResponse `json:"recent_donations,omitempty"`
}

type DonationItemResponse struct {
	ID             string  `json:"id"`
	ItemName       string  `json:"item_name"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	EstimatedValue float64 `json:"estimated_value"`
	Notes          string  `json:"notes,omitempty"`
}

type DonationResponse struct {
	ID            string                 `json:"id"`
	Date          time.Time              `json:"date"`
	Type          string                 `json:"type"`
	PaymentMethod string                 `json:"payment_method"`
	TotalAmount   float64                `json:"total_amount"`
	ProofFileURL  string                 `json:"proof_file_url,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Donor         DonorResponse          `json:"donor"`
	Employee      SimpleEmployeeResponse `json:"employee"` // Simplified employee info
	Items         []DonationItemResponse `json:"items,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
}

type SimpleEmployeeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func FromDomainDonor(d *domain.Donor) DonorResponse {
	res := DonorResponse{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
	if d.Phone != nil {
		res.Phone = *d.Phone
	}
	if d.Email != nil {
		res.Email = *d.Email
	}
	if d.Address != nil {
		res.Address = *d.Address
	}
	return res
}

func FromDomainDonation(d *domain.Donation, baseURL string) DonationResponse {
	res := DonationResponse{
		ID:            d.ID,
		Date:          d.Date,
		Type:          d.Type,
		PaymentMethod: d.PaymentMethod,
		TotalAmount:   d.TotalAmount,
		Description:   utilsPtrToString(d.Description),
		ProofFileURL:  GenerateFileURL(d.ProofFile, baseURL),
		Donor:         FromDomainDonor(&d.Donor),
		Items:         []DonationItemResponse{},
		CreatedAt:     d.CreatedAt,
		Employee: SimpleEmployeeResponse{
			ID:   d.EmployeeID,
			Name: d.Employee.FullName,
		},
	}

	for _, item := range d.Items {
		res.Items = append(res.Items, DonationItemResponse{
			ID:             item.ID,
			ItemName:       item.ItemName,
			Quantity:       item.Quantity,
			Unit:           item.Unit,
			EstimatedValue: item.EstimatedValue,
			Notes:          utilsPtrToString(item.Notes),
		})
	}

	return res
}

func utilsPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func GenerateFileURL(path *string, baseURL string) string {
	if path != nil && *path != "" {
		// If path already contains http, return as is (just in case)
		// But usually it's relative path "donations/..."
		return fmt.Sprintf("%s/api/v1/files/%s", baseURL, *path)
	}
	return ""
}
