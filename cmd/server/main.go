package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	port = "8000"
	path = "/payments"
)

var paymentStore = NewPaymentStore()

// Subset of ISO 4217 currency codes
var supportedCurrencies = map[string]bool{
	"EUR": true, "GBP": true,
}

func main() {
	log.Println("new payment store instantiated: ", paymentStore)
	router := mux.NewRouter()
	router.HandleFunc(path, ProcessPaymentHandler).Methods("POST")
	router.HandleFunc(path+"/{id}", GetPaymentHandler).Methods("GET")
	log.Printf("server listening on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), router))
}
