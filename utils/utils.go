package utils

import "github.com/celestebrant/processout-payment-gateway/models"

const Path = "/payments"

// ValidProcessPaymentRequest generates a valid ProcessPaymentRequest which is useful for testing.
func ValidProcessPaymentRequest() *models.ProcessPaymentRequest {
	return &models.ProcessPaymentRequest{
		CardNumber:  "1234123412341234",
		ExpiryYear:  2099,
		ExpiryMonth: 12,
		CVV:         "987",
		Amount:      10.05,
		Currency:    "GBP",
	}
}
