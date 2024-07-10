package bank

import (
	"bytes"
	"encoding/json"
	"fmt"

	"golang.org/x/exp/rand"

	"github.com/celestebrant/processout-payment-gateway/models"
	"github.com/google/uuid"
)

func MockCallBank() (*models.CallBankResponse, error) {
	mockData := generateMockData()
	mockDataJSON, _ := marshalMockData(mockData)

	callBankResponse, err := DecodeBankResponse(mockDataJSON)
	if err != nil {
		return nil, err
	}

	return callBankResponse, nil
}

func marshalMockData(data map[string]string) ([]byte, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mocked payment response from bank server")
	}

	return dataJSON, err
}

// DecodeBankResponse decodes the bank response into CallBankResponse
func DecodeBankResponse(responseJSON []byte) (*models.CallBankResponse, error) {
	callBankResponse := &models.CallBankResponse{}
	data := bytes.NewReader(responseJSON)
	err := json.NewDecoder(data).Decode(callBankResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payment response from bank server: %w", err)
	}

	return callBankResponse, nil
}

func generateMockData() map[string]string {
	paymentID := uuid.New().String()

	// 20% chance of failure
	status := "SUCCESS"
	if rand.Intn(10) < 2 {
		status = "FAILED"
	}

	// Assume JSON bank response has been converted to map[string]string
	mockedData := map[string]string{
		"payment_id": paymentID,
		"status":     status,
	}

	return mockedData
}
