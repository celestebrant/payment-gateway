package bank

import (
	"bytes"
	"encoding/json"
	"fmt"

	"golang.org/x/exp/rand"

	"github.com/celestebrant/processout-payment-gateway/models"
	"github.com/google/uuid"
)

// MockCallBank mocks a call to an external bank server and then returns the response that
// is decoded into CallBankResponse. It is assumed that the data returned are: payment_id, status.
func MockCallBank() (*models.CallBankResponse, error) {
	mockData := generateMockData()
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
func DecodeBankResponse(responseJSON []byte) (*models.CallBankResponse, error) {
	callBankResponse := models.CallBankResponse{}
	data := bytes.NewReader(responseJSON)
	err := json.NewDecoder(data).Decode(&callBankResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payment response from bank server: %w", err)
	}

	return &callBankResponse, nil
}

// generateMockData generates mock data with fields payment_id (36 characters) and status, returned in a map.
// There is an 80% chance that status is "SUCCESS" and 20% chance of "FAILED".
func generateMockData() map[string]string {
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
