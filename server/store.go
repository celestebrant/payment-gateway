package server

import (
	"sync"

	"github.com/celestebrant/processout-payment-gateway/models"
)

type PaymentStore struct {
	mu       sync.Mutex
	payments map[string]*models.MaskedPayment
}

func NewPaymentStore() *PaymentStore {
	return &PaymentStore{
		payments: make(map[string]*models.MaskedPayment),
	}
}

func (s *PaymentStore) AddPayment(payment *models.MaskedPayment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payments[payment.ID] = payment
}

func (s *PaymentStore) GetPayment(id string) (*models.MaskedPayment, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	payment, exists := s.payments[id]
	return payment, exists
}
