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

type ProcessPaymentRequest struct {
	CardNumber  string  `json:"card_number"`
	ExpiryYear  uint    `json:"expiry_year"`
	ExpiryMonth uint    `json:"expiry_month"`
	CVV         string  `json:"cvv"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
}
