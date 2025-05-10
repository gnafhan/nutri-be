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

func (m *MockPayment) CreateTransaction(orderID string, amount int, userDetails map[string]interface{}, paymentMethod string) (*PaymentToken, error) {
	return &PaymentToken{
		Token:       "mock_token_" + uuid.New().String(),
		RedirectURL: "https://example.com/mock_payment",
	}, nil
}

func (m *MockPayment) CheckTransactionStatus(transactionID string) (interface{}, error) {
	return map[string]string{
		"transaction_id": transactionID,
		"status":         "settlement",
	}, nil
}

func (m *MockPayment) HandleNotification(notificationJSON []byte) (interface{}, error) {
	return map[string]string{
		"status": "success",
	}, nil
}
