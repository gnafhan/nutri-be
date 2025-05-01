package service

import "github.com/google/uuid"

type MockPayment struct{}

func (m *MockPayment) Charge(amount int, method string) (*PaymentResponse, error) {
	return &PaymentResponse{
		TransactionID: "mock_" + uuid.New().String(),
		Status:        "completed",
	}, nil
}

func (m *MockPayment) Refund(transactionID string) error {
	return nil
}
