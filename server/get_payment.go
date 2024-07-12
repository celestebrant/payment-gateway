package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GetPaymentHandler handles fetching individual payments by ID.
func GetPaymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) > 36 {
		http.Error(w, "payment ID should have up to 36 characters", http.StatusBadRequest)
		return
	}

	maskedPayment, exists := paymentStore.GetPayment(id)
	if !exists {
		http.Error(w, "payment not found", http.StatusNotFound)
		return
	}

	log.Println("Fetched payment:", *maskedPayment)

	json.NewEncoder(w).Encode(maskedPayment)
}
