package request

type CreateDonationRequest struct {
	// Donor Info
	DonorName    string `form:"donor_name" binding:"required"`
	DonorPhone   string `form:"donor_phone"` // Optional, used for lookup
	DonorEmail   string `form:"donor_email"`
	DonorAddress string `form:"donor_address"`

	// Donation Info
	Type          string  `form:"type" binding:"required,oneof=MONEY GOODS MIXED"`
	PaymentMethod string  `form:"payment_method" binding:"required,oneof=CASH TRANSFER QRIS GOODS"`
	TotalAmount   float64 `form:"total_amount"`
	Description   string  `form:"description"`
	ProofFile     string  `form:"-"` // Manually handled by handler

	// Items (JSON string or multiple fields? Form-data with complex arrays is tricky)
	// Strategy: We will accept a JSON string for items and parse it in the handler
	ItemsJSON string `form:"items_json"`
}

type DonationItemRequest struct {
	ItemName       string  `json:"item_name"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	EstimatedValue float64 `json:"estimated_value"`
	Notes          string  `json:"notes"`
}

type UpdateDonationRequest struct {
	Date          string  `form:"date"`
	Type          string  `form:"type" binding:"omitempty,oneof=MONEY GOODS MIXED"`
	PaymentMethod string  `form:"payment_method" binding:"omitempty,oneof=CASH TRANSFER QRIS GOODS"`
	TotalAmount   float64 `form:"total_amount"`
	Description   string  `form:"description"`
	ItemsJSON     string  `form:"items_json"`
	ProofFile     string  `form:"-"`
}

type UpdateDonorRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}
