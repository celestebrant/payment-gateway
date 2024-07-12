package mockbank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/exp/rand"

	"github.com/google/uuid"
)

// BankClient is the client for making mocked requests to the bank.
type BankClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewBankClient instantiates a new bank client.
func NewBankClient() *BankClient {
	return &BankClient{
		BaseURL:    "mockbank.com",
		HTTPClient: &http.Client{},
	}
}

// MakePaymentRequest represents the assumed request data the bank API requires, including card details
// and data about the money to be transacted.
type MakePaymentRequest struct {
	CardNumber  string  `json:"card_number"`
	ExpiryYear  uint    `json:"expiry_year"`
	ExpiryMonth uint    `json:"expiry_month"`
	CVV         string  `json:"cvv"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
}

// MakePaymentResponse represents the assumed response the bank API returns, containing payment ID and status.
type MakePaymentResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

// MakePayment mocks a call to an external bank server and then returns the response that
// is decoded into CallBankResponse. It is assumed that the data returned are: payment_id, status.
func (b *BankClient) MakePayment(r MakePaymentRequest) (*MakePaymentResponse, error) {
	// Generate CallBankResponse with mock data
	mockData := generateMockedData()
	mockDataJSON, err := json.Marshal(mockData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mocked make payment request")
	}

	callBankResponse, err := decodeBankResponse(mockDataJSON)
	if err != nil {
		return nil, err
	}

	return callBankResponse, nil
}

// decodeBankResponse decodes the bank response into CallBankResponse.
func decodeBankResponse(responseJSON []byte) (*MakePaymentResponse, error) {
	callBankResponse := MakePaymentResponse{}
	data := bytes.NewReader(responseJSON)
	err := json.NewDecoder(data).Decode(&callBankResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payment response from bank server: %w", err)
	}

	return &callBankResponse, nil
}

// generateMockedData returns a map of mocked data with fields payment_id (36 characters) and status.
// There is an 80% chance that status is "SUCCESS" and 20% chance of "FAILED".
func generateMockedData() map[string]string {
	paymentID := uuid.New().String()

	status := "SUCCESS"
	if rand.Intn(10) < 2 {
		status = "FAILED"
	}

	// Imagine the data was received in a JSON response, then was converted to a map
	mockedData := map[string]string{
		"payment_id": paymentID,
		"status":     status,
	}

	return mockedData
}
