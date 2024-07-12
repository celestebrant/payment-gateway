package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/celestebrant/processout-payment-gateway/bank"
	"github.com/celestebrant/processout-payment-gateway/models"
)

// Subset of ISO 4217 currency codes
var supportedCurrencies = map[string]bool{
	"EUR": true, "GBP": true,
}

var paymentStore *PaymentStore
var once sync.Once

func init() {
	once.Do(func() {
		paymentStore = NewPaymentStore()
	})
}

// ProcessPaymentHandler handles process payment requests.
func ProcessPaymentHandler(w http.ResponseWriter, r *http.Request) {
	request := models.ProcessPaymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "failed to unmarshal the request", http.StatusBadRequest)
		return
	}

	if err := validateProcessPaymentRequest(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bankResponse, err := bank.MockCallBank()
	if err != nil {
		http.Error(w, "unexpected error with bank", http.StatusInternalServerError)
		return
	}

	// TODO: Add payment to store
	maskedPayment := populateMaskedPayment(request, bankResponse.PaymentID, bankResponse.Status)
	paymentStore.AddPayment(maskedPayment)
	log.Println("Processed payment:", *maskedPayment)

	json.NewEncoder(w).Encode(maskedPayment)
}

/*
validateProcessPaymentRequest validates the data in request with the following rules:
  - Card number must be exactly 16 digits long with numerical characters only
  - Expiry year must be exactly 4 digits long, however no validation is performed
    relative to current time
  - Expiry month must have an integer value of 1 to 12, inclusive
  - CVV must be exactly 3 digits long with numerical characters only
  - Amount must be a positive number with up to 2 decimal places
  - Currency must be either GBP or EUR
*/
func validateProcessPaymentRequest(request models.ProcessPaymentRequest) error {
	cardNumberPattern := regexp.MustCompile(`^\d{16}$`)
	if !cardNumberPattern.MatchString(request.CardNumber) {
		return fmt.Errorf("card number should have 16 digits")
	}

	if request.ExpiryYear < 1000 || request.ExpiryYear > 9999 {
		return fmt.Errorf("expiry year should have 4 digits")
	}

	if request.ExpiryMonth == 0 || request.ExpiryMonth > 12 {
		return fmt.Errorf("expiry month should have value of 1 to 12")
	}

	cvvPattern := regexp.MustCompile(`^\d{3}$`)
	if !cvvPattern.MatchString(request.CVV) {
		return fmt.Errorf("cvv should have 3 digits")
	}

	amountStr := fmt.Sprintf("%.2f", request.Amount)
	amountParsedBack, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return fmt.Errorf("an unexpected error occurred when parsing amount")
	}
	if amountParsedBack != request.Amount {
		return fmt.Errorf("amount must have up to two decimal places")
	}

	amountPattern := regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	if !amountPattern.MatchString(fmt.Sprintf("%.2f", request.Amount)) {
		return fmt.Errorf("amount must be a positive number with up to two decimal places")
	}
	if request.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	if !supportedCurrencies[request.Currency] {
		return fmt.Errorf("invalid currency code")
	}

	return nil
}

// populateMaskedPayment returns a MaskedPayment with values from the provided request, id and status.
func populateMaskedPayment(request models.ProcessPaymentRequest, id, status string) *models.MaskedPayment {
	return &models.MaskedPayment{
		ID:               id,
		Status:           status,
		MaskedCardNumber: maskCardNumber(request.CardNumber),
		ExpiryYear:       request.ExpiryYear,
		ExpiryMonth:      request.ExpiryMonth,
		Amount:           request.Amount,
		Currency:         request.Currency,
	}
}

// maskCardNumber returns the card number with * for all digits but the final 4, like ************XXXX.
func maskCardNumber(cardNumber string) string {
	return "************" + cardNumber[len(cardNumber)-4:]
}
