package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPayment(t *testing.T) {
	t.Parallel()

	router := mux.NewRouter()
	router.HandleFunc(path, ProcessPaymentHandler).Methods("POST")
	router.HandleFunc(path+"/{id}", GetPaymentHandler).Methods("GET")

	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("nonexistent payment returns error", func(t *testing.T) {
		// Setup
		r, a := require.New(t), assert.New(t)

		// Get the payment
		request, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s", server.URL, path, uuid.New().String()), nil)
		r.NoError(err, "failed to create retrieve request")

		response, err := http.DefaultClient.Do(request)
		r.NoError(err, "failed to retrieve payment")
		defer response.Body.Close()

		a.Equal(http.StatusNotFound, response.StatusCode)

		responseBody, err := io.ReadAll(response.Body)
		r.NoError(err, "failed to read response body")
		r.Equal("payment not found", string(bytes.TrimSpace(responseBody)))
	})

	t.Run("payment ID too long returns error", func(t *testing.T) {
		// Setup
		r, a := require.New(t), assert.New(t)

		// Get the payment
		request, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%sa", server.URL, path, uuid.New().String()), nil)
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
