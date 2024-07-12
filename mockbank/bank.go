package mockbank

import (
	"bytes"
	"encoding/json"
	"fmt"

	"golang.org/x/exp/rand"

	"github.com/google/uuid"
)

// CallBankRequest represents the assumed request data the bank API requires, including card details
// and data about the money to be transacted.
type CallBankRequest struct {
	CardNumber  string  `json:"card_number"`
	ExpiryYear  uint    `json:"expiry_year"`
	ExpiryMonth uint    `json:"expiry_month"`
	CVV         string  `json:"cvv"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
}

// CallBankResponse represents the assumed response the bank API returns, containing payment ID and status.
type CallBankResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

// CallBank mocks a call to an external bank server and then returns the response that
// is decoded into CallBankResponse. It is assumed that the data returned are: payment_id, status.
func CallBank(r CallBankRequest) (*CallBankResponse, error) {
	// Pretend a call is made to the bank API using the CallBankRequest data
	_ = r

	// Generate CallBankResponse with mock data
	mockData := generateMockedData()
	mockDataJSON, _ := marshalMockData(mockData)

	callBankResponse, err := DecodeBankResponse(mockDataJSON)
	if err != nil {
		return nil, err
	}

	return callBankResponse, nil
}

// marshalMockData converts input map data to []byte, which can be used for decoding.
func marshalMockData(data map[string]string) ([]byte, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mocked payment response from bank server")
	}

	return dataJSON, err
}

// DecodeBankResponse decodes the bank response into CallBankResponse.
func DecodeBankResponse(responseJSON []byte) (*CallBankResponse, error) {
	callBankResponse := CallBankResponse{}
	data := bytes.NewReader(responseJSON)
	err := json.NewDecoder(data).Decode(&callBankResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payment response from bank server: %w", err)
	}

	return &callBankResponse, nil
}

// generateMockedData generates mock data with fields payment_id (36 characters) and status, returned in a map.
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
