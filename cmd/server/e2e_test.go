package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celestebrant/processout-payment-gateway/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestEndToEndPaymentFlow(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	// Setup
	router := mux.NewRouter()
	router.HandleFunc(path, ProcessPaymentHandler).Methods("POST")
	router.HandleFunc(path+"/{id}", GetPaymentHandler).Methods("GET")

	server := httptest.NewServer(router)
	defer server.Close()

	// Process payment (pp)
	body, err := json.Marshal(validProcessPaymentRequest())
	r.NoError(err, "failed to marshal request")

	ppRequest, err := http.NewRequest("POST", server.URL+path, bytes.NewReader(body))
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
	gpRequest, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s", server.URL, path, maskedPayment.ID), nil)
	r.NoError(err, "failed to create retrieve request")

	gpResponse, err := http.DefaultClient.Do(gpRequest)
	r.NoError(err, "failed to retrieve payment")
	defer gpResponse.Body.Close()

	r.Equal(http.StatusOK, gpResponse.StatusCode)

	var retrievedMaskedPayment models.MaskedPayment
	err = json.NewDecoder(gpResponse.Body).Decode(&retrievedMaskedPayment)
	r.NoError(err, "failed to unmarshal retrieve payment response")

	r.Equal(maskedPayment, retrievedMaskedPayment, "processed payment and retrieved payment should match")
}
