package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/celestebrant/processout-payment-gateway/models"
	"github.com/gorilla/mux"
)

const (
	port = "8000"
)

var GlobPayments struct {
	mu       sync.Mutex
	Payments map[string]*models.Payment
}

// Subset of ISO 4217 currency codes
var supportedCurrencies = map[string]bool{
	"EUR": true, "GBP": true,
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/process-payment", ProcessPaymentHandler).Methods("POST")
	log.Printf("server listening on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), router))
}
