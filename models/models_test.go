package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProcessPaymentRequest(t *testing.T) {
	r := require.New(t)
	request := NewProcessPaymentRequest()
	expectedRequest := ProcessPaymentRequest{
		CardNumber:  "",
		ExpiryYear:  0,
		ExpiryMonth: 0,
		CVV:         "",
		Amount:      0,
		Currency:    "",
	}
	r.Equal(expectedRequest, *request)
}

func TestNewCallBankResponse(t *testing.T) {
	r := require.New(t)
	request := NewCallBankResponse()
	expectedRequest := CallBankResponse{
		PaymentID: "",
		Status:    "",
	}
	r.Equal(expectedRequest, *request)
}
