package main

import (
	"log"
	"net/http"

	"github.com/celestebrant/processout-payment-gateway/server"
	"github.com/celestebrant/processout-payment-gateway/utils"
	"github.com/gorilla/mux"
)

const (
	port = "8000"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc(utils.Path, server.ProcessPaymentHandler).Methods("POST")
	router.HandleFunc(utils.Path+"/{id}", server.GetPaymentHandler).Methods("GET")
	log.Printf("server listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
