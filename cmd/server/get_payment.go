package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetPaymentHandler handles fetching individual payments by ID.
func GetPaymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	maskedPayment, exists := paymentStore.GetPayment(id)
	if !exists {
		http.Error(w, "payment not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(maskedPayment)
}
