package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celestebrant/processout-payment-gateway/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessPaymentHandler(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                 string
		modifyRequest        func(req *models.ProcessPaymentRequest)
		expectedStatusCode   int
		expectedPayment      models.Payment
		expectedErrorMessage string
	}

	testCases := []testCase{
		{
			"valid payload returns payment response",
			func(req *models.ProcessPaymentRequest) {},
			http.StatusOK,
			models.Payment{
				CardNumber:  "1234123412341234",
				ExpiryYear:  2099,
				ExpiryMonth: 12,
				CVV:         "987",
				Amount:      10.05,
				Currency:    "GBP",
			},
			"",
		}, {
			"validation error returns 400 error response",
			func(req *models.ProcessPaymentRequest) {
				req.CardNumber = ""
			},
			http.StatusBadRequest,
			models.Payment{},
			"card number should have 16 digits",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			processPaymentRequest := validProcessPaymentRequest()
			tc.modifyRequest(processPaymentRequest)
			body, err := json.Marshal(processPaymentRequest)
			r.NoError(err, "failed to marshal request")

			// Request and response implement http.ResponseWriter and http.Request
			request := httptest.NewRequest("POST", "/process-payment", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			ProcessPaymentHandler(response, request)
			defer response.Result().Body.Close()

			r.Equal(tc.expectedStatusCode, response.Result().StatusCode)

			if tc.expectedErrorMessage == "" {
				// Positive scenarios: response should contain payment data
				var payment models.Payment
				err = json.NewDecoder(response.Body).Decode(&payment)
				r.NoError(err, "failed to unmarshal response")

				expectedPayment := models.Payment{
					CardNumber:  tc.expectedPayment.CardNumber,
					ExpiryYear:  tc.expectedPayment.ExpiryYear,
					ExpiryMonth: tc.expectedPayment.ExpiryMonth,
					CVV:         tc.expectedPayment.CVV,
					Amount:      tc.expectedPayment.Amount,
					Currency:    tc.expectedPayment.Currency,
				}

				if diff := cmp.Diff(expectedPayment, payment, cmpopts.IgnoreFields(models.Payment{}, "ID", "Status")); diff != "" {
					t.Errorf("Payment mismatch (-expected +got):\n%s", diff)
				}
				r.NotEmpty(payment.ID)
				r.True(
					payment.Status == "SUCCESS" || payment.Status == "FAILED",
					`expected status to be either "SUCCESS" or "FAILED", got "%s"`, payment.Status,
				)

			} else {
				// Negative scenarios: response should contain an error
				responseBody, err := io.ReadAll(response.Body)
				r.NoError(err, "failed to read response body")
				r.Equal(tc.expectedErrorMessage, string(bytes.TrimSpace(responseBody)))
			}
		})
	}
}

func TestValidateProcessPaymentRequest(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                 string
		modifyRequest        func(req *models.ProcessPaymentRequest)
		expectedErrorMessage string
	}

	testCases := []testCase{
		{
			"valid",
			func(req *models.ProcessPaymentRequest) {},
			"",
		}, {
			"card number too short returns error",
			func(req *models.ProcessPaymentRequest) {
				req.CardNumber = "123412341234123" // 15 digits
			},
			"card number should have 16 digits",
		}, {
			"card number too long returns error",
			func(req *models.ProcessPaymentRequest) {
				req.CardNumber = "12341234123412345" // 17 digits
			},
			"card number should have 16 digits",
		}, {
			"card number alphanumeric returns error",
			func(req *models.ProcessPaymentRequest) {
				req.CardNumber = "a234123412341234"
			},
			"card number should have 16 digits",
		}, {
			"card number special characters returns error",
			func(req *models.ProcessPaymentRequest) {
				req.CardNumber = " 234123412341234"
			},
			"card number should have 16 digits",
		}, {
			"expiry year 3 digits returns error",
			func(req *models.ProcessPaymentRequest) {
				req.ExpiryYear = 999
			},
			"expiry year should have 4 digits",
		}, {
			"expiry year 4 digits lower bound",
			func(req *models.ProcessPaymentRequest) {
				req.ExpiryYear = 1000
			},
			"",
		}, {
			"expiry year 4 digits upper bound",
			func(req *models.ProcessPaymentRequest) {
				req.ExpiryYear = 9999
			},
			"",
		}, {
			"expiry year 5 digits returns error",
			func(req *models.ProcessPaymentRequest) {
				req.ExpiryYear = 10000
			},
			"expiry year should have 4 digits",
		}, {
			"expiry month 0 returns error",
			func(req *models.ProcessPaymentRequest) {
				req.ExpiryMonth = 0
			},
			"expiry month should have value of 1 to 12",
		}, {
			"expiry month 1",
			func(req *models.ProcessPaymentRequest) {
				req.ExpiryMonth = 1 // jan
			},
			"",
		}, {
			"expiry month 12",
			func(req *models.ProcessPaymentRequest) {
				req.ExpiryMonth = 12 // dec
			},
			"",
		}, {
			"expiry month 13 returns error",
			func(req *models.ProcessPaymentRequest) {
				req.ExpiryMonth = 13
			},
			"expiry month should have value of 1 to 12",
		}, {
			"CVV too short returns error",
			func(req *models.ProcessPaymentRequest) {
				req.CVV = "99"
			},
			"cvv should have 3 digits",
		}, {
			"CVV too long returns error",
			func(req *models.ProcessPaymentRequest) {
				req.CVV = "1000"
			},
			"cvv should have 3 digits",
		}, {
			"CVV alphanumeric returns error",
			func(req *models.ProcessPaymentRequest) {
				req.CVV = "a12"
			},
			"cvv should have 3 digits",
		}, {
			"CVV special characters returns error",
			func(req *models.ProcessPaymentRequest) {
				req.CVV = " 12"
			},
			"cvv should have 3 digits",
		}, {
			"amount 0 returns error",
			func(req *models.ProcessPaymentRequest) {
				req.Amount = 0
			},
			"amount must be greater than zero",
		}, {
			"amount one decimal place",
			func(req *models.ProcessPaymentRequest) {
				req.Amount = 0.1
			},
			"",
		}, {
			"amount two decimal places",
			func(req *models.ProcessPaymentRequest) {
				req.Amount = 0.01
			},
			"",
		}, {
			"amount more than two decimal places returns error",
			func(req *models.ProcessPaymentRequest) {
				req.Amount = 0.009
			},
			"amount must have up to two decimal places",
		}, {
			"amount negative returns error",
			func(req *models.ProcessPaymentRequest) {
				req.Amount = -0.01
			},
			"amount must be a positive number with up to two decimal places",
		}, {
			"currency unsupported returns error",
			func(req *models.ProcessPaymentRequest) {
				req.Currency = "USD"
			},
			"invalid currency code",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			req := validProcessPaymentRequest()
			tc.modifyRequest(req)
			err := validateProcessPaymentRequest(*req)

			if tc.expectedErrorMessage == "" {
				// Positive scenarios
				r.NoError(err)
			} else {
				// Negative scenarios
				r.ErrorContains(err, tc.expectedErrorMessage)
			}
		})
	}
}

func TestPopulatePayment(t *testing.T) {
	a := assert.New(t)

	id := "some-id"
	status := "some-status"
	request := models.ProcessPaymentRequest{
		CardNumber:  "some-card-number",
		ExpiryYear:  1,
		ExpiryMonth: 1,
		CVV:         "some-cvv",
		Amount:      1.234,
		Currency:    "some-currency",
	}

	payment := populatePayment(request, id, status)
	a.Equal(id, payment.ID)
	a.Equal(status, payment.Status)
	a.Equal(request.CardNumber, payment.CardNumber)
	a.Equal(request.ExpiryYear, payment.ExpiryYear)
	a.Equal(request.ExpiryMonth, payment.ExpiryMonth)
	a.Equal(request.CVV, payment.CVV)
	a.Equal(request.Amount, payment.Amount)
	a.Equal(request.Currency, payment.Currency)
}

func validProcessPaymentRequest() *models.ProcessPaymentRequest {
	return &models.ProcessPaymentRequest{
		CardNumber:  "1234123412341234",
		ExpiryYear:  2099,
		ExpiryMonth: 12,
		CVV:         "987",
		Amount:      10.05,
		Currency:    "GBP",
	}
}
