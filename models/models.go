package models

type MaskedPayment struct {
	ID               string  `json:"id"`
	Status           string  `json:"status"`
	MaskedCardNumber string  `json:"masked_card_number"`
	ExpiryYear       uint    `json:"expiry_year"`
	ExpiryMonth      uint    `json:"expiry_month"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
}

type BankResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ProcessPaymentRequest struct {
	CardNumber  string  `json:"card_number"`
	ExpiryYear  uint    `json:"expiry_year"`
	ExpiryMonth uint    `json:"expiry_month"`
	CVV         string  `json:"cvv"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
}

// NewProcessPaymentRequest creates a new ProcessPaymentRequest with zero values.
func NewProcessPaymentRequest() *ProcessPaymentRequest {
	return &ProcessPaymentRequest{
		CardNumber:  "",
		ExpiryYear:  0,
		ExpiryMonth: 0,
		CVV:         "",
		Amount:      0,
		Currency:    "",
	}
}

type CallBankResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

// NewCallBankResponse creates a new CallBankResponse with zero values.
func NewCallBankResponse() *CallBankResponse {
	return &CallBankResponse{
		PaymentID: "",
		Status:    "",
	}
}
