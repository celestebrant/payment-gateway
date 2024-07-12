package mockbank

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockData(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	// Verify mock data generation
	mockDataMap := generateMockedData()
	a.NotEmpty(mockDataMap["payment_id"])
	r.True(
		mockDataMap["status"] == "SUCCESS" || mockDataMap["status"] == "FAILED",
		`expected status to be either "SUCCESS" or "FAILED", got "%s"`, mockDataMap["status"],
	)

	// Verify decoding process propagates values correctly
	mockDataJSON, err := json.Marshal(mockDataMap)
	r.NoError(err)

	callBankResponse, err := decodeBankResponse(mockDataJSON)
	r.NoError(err)
	a.Equal(mockDataMap["payment_id"], callBankResponse.PaymentID, "payment ID value should be propagated")
	r.Equal(mockDataMap["status"], callBankResponse.Status, "payment ID value should be propagated")
}

func TestCallBank(t *testing.T) {
	r, a := require.New(t), assert.New(t)
	bankClient := NewBankClient()
	callBankRequest := MakePaymentRequest{}
	callBankResponse, err := bankClient.MakePayment(callBankRequest)
	r.NoError(err)
	a.NotEmpty(callBankResponse.PaymentID)
	a.NotEmpty(callBankResponse.Status)
}
