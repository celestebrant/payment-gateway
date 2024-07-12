package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celestebrant/processout-payment-gateway/models"
	"github.com/celestebrant/processout-payment-gateway/server"
	"github.com/celestebrant/processout-payment-gateway/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndToEndPaymentFlow(t *testing.T) {
	t.Parallel()

	// Setup
	router := mux.NewRouter()
	router.HandleFunc(utils.Path, server.ProcessPaymentHandler).Methods("POST")
	router.HandleFunc(utils.Path+"/{id}", server.GetPaymentHandler).Methods("GET")

	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("process then fetch payment", func(t *testing.T) {
		r := require.New(t)

		// Process payment (pp)
		body, err := json.Marshal(utils.ValidProcessPaymentRequest())
		r.NoError(err, "failed to marshal request")

		ppRequest, err := http.NewRequest("POST", server.URL+utils.Path, bytes.NewReader(body))
		r.NoError(err, "failed to create request")
		ppRequest.Header.Set("Content-Type", "application/json")

		ppResponse, err := http.DefaultClient.Do(ppRequest)
		r.NoError(err, "failed to process payment request")
		defer ppResponse.Body.Close()

		r.Equal(http.StatusOK, ppResponse.StatusCode)

		var maskedPayment models.MaskedPayment
		err = json.NewDecoder(ppResponse.Body).Decode(&maskedPayment)
		r.NoError(err, "failed to unmarshal process payment response")

		// Get the payment (gp)
		gpRequest, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s", server.URL, utils.Path, maskedPayment.ID), nil)
		r.NoError(err, "failed to create retrieve request")

		gpResponse, err := http.DefaultClient.Do(gpRequest)
		r.NoError(err, "failed to retrieve payment")
		defer gpResponse.Body.Close()

		r.Equal(http.StatusOK, gpResponse.StatusCode)

		var retrievedMaskedPayment models.MaskedPayment
		err = json.NewDecoder(gpResponse.Body).Decode(&retrievedMaskedPayment)
		r.NoError(err, "failed to unmarshal retrieve payment response")

		r.Equal(maskedPayment, retrievedMaskedPayment, "processed payment and retrieved payment should match")
	})

	t.Run("process payment validation error", func(t *testing.T) {
		r := require.New(t)

		data := utils.ValidProcessPaymentRequest()
		data.CVV = "1" // invalid
		body, err := json.Marshal(data)
		r.NoError(err, "failed to marshal request")

		request, err := http.NewRequest("POST", server.URL+utils.Path, bytes.NewReader(body))
		r.NoError(err, "failed to create request")
		request.Header.Set("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(request)
		r.NoError(err, "failed to process payment request")
		defer response.Body.Close()

		r.Equal(http.StatusBadRequest, response.StatusCode)

		responseBody, err := io.ReadAll(response.Body)
		r.NoError(err, "failed to read response body")
		r.Contains(
			string(bytes.TrimSpace(responseBody)),
			"cvv",
			"response error should contain some error message related to cvv",
		)
	})

	t.Run("get nonexistent payment returns error", func(t *testing.T) {
		r, a := require.New(t), assert.New(t)

		request, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s", server.URL, utils.Path, uuid.New().String()), nil)
		r.NoError(err, "failed to create retrieve request")

		response, err := http.DefaultClient.Do(request)
		r.NoError(err, "failed to retrieve payment")
		defer response.Body.Close()

		a.Equal(http.StatusNotFound, response.StatusCode)

		responseBody, err := io.ReadAll(response.Body)
		r.NoError(err, "failed to read response body")
		r.Equal("payment not found", string(bytes.TrimSpace(responseBody)))
	})

	t.Run("get payment ID too long returns error", func(t *testing.T) {
		r, a := require.New(t), assert.New(t)

		request, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%sa", server.URL, utils.Path, uuid.New().String()), nil)
		r.NoError(err, "failed to create retrieve request")

		response, err := http.DefaultClient.Do(request)
		r.NoError(err, "failed to retrieve payment")
		defer response.Body.Close()

		a.Equal(http.StatusBadRequest, response.StatusCode)

		responseBody, err := io.ReadAll(response.Body)
		r.NoError(err, "failed to read response body")
		r.Equal("payment ID should have up to 36 characters", string(bytes.TrimSpace(responseBody)))
	})

}
