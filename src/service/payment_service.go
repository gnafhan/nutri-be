package service

import (
	"app/src/config"
	midtransutils "app/src/midtrans"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type MidtransPaymentService struct {
	SnapClient    snap.Client
	CoreAPIClient coreapi.Client
	Log           *logrus.Logger
	IsProduction  bool
}

type PaymentToken struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func NewMidtransPaymentService() *MidtransPaymentService {
	var snapClient snap.Client
	var coreAPIClient coreapi.Client

	serverKey := config.MidtransServerKey
	isProduction := config.MidtransStatus == "PRODUCTION"

	// Set environment
	env := midtrans.Sandbox
	if isProduction {
		env = midtrans.Production
	}

	// Initialize clients with server key and environment
	snapClient.New(serverKey, env)
	coreAPIClient.New(serverKey, env)

	return &MidtransPaymentService{
		SnapClient:    snapClient,
		CoreAPIClient: coreAPIClient,
		Log:           logrus.New(),
		IsProduction:  isProduction,
	}
}

func (s *MidtransPaymentService) CreateTransaction(orderID string, amount int, userDetails map[string]interface{}, paymentMethod string) (*PaymentToken, error) {
	// Create transaction request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: userDetails["first_name"].(string),
			LName: userDetails["last_name"].(string),
			Email: userDetails["email"].(string),
			Phone: userDetails["phone"].(string),
		},
	}

	// Add payment method specific configuration if specified
	if paymentMethod != "" {
		if paymentMethod == "credit_card" {
			req.CreditCard = &snap.CreditCardDetails{
				Secure: true,
			}
		} else if paymentMethod == "gopay" || paymentMethod == "shopeepay" {
			req.EnabledPayments = []snap.SnapPaymentType{snap.SnapPaymentType(paymentMethod)}
		} else if paymentMethod == "bank_transfer" {
			req.EnabledPayments = []snap.SnapPaymentType{
				snap.PaymentTypeBCAVA,
				snap.PaymentTypeBNIVA,
				snap.PaymentTypeBRIVA,
				snap.PaymentTypePermataVA,
			}
		}
	}
	// When paymentMethod is not specified, Midtrans will show all available payment methods

	// Create Snap transaction
	snapResp, err := s.SnapClient.CreateTransaction(req)
	if err != nil {
		return nil, fmt.Errorf("error creating snap transaction: %w", err)
	}

	return &PaymentToken{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

func (s *MidtransPaymentService) CheckTransactionStatus(transactionID string) (interface{}, error) {
	response, err := s.CoreAPIClient.CheckTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("error checking transaction status: %w", err)
	}

	return response, nil
}

func (s *MidtransPaymentService) HandleNotification(notificationJSON []byte) (interface{}, error) {
	// s.Log.Infof("HandleNotification called with data: %s", string(notificationJSON))

	var notificationPayload map[string]interface{}

	err := json.Unmarshal(notificationJSON, &notificationPayload)
	if err != nil {
		s.Log.Errorf("Error parsing notification JSON: %v", err)
		return nil, fmt.Errorf("error parsing notification JSON: %w", err)
	}

	// Verify signature key
	isValidSignature, err := s.verifySignatureKey(notificationPayload)
	if err != nil {
		s.Log.Errorf("Error verifying signature key: %v", err)
		return nil, fmt.Errorf("error verifying signature: %w", err)
	}

	if !isValidSignature {
		s.Log.Error("Invalid signature key, possible security threat")
		return nil, errors.New("invalid signature key")
	}

	s.Log.Info("Signature verification successful")

	orderID, exists := notificationPayload["order_id"].(string)
	if !exists {
		s.Log.Error("Notification does not contain order_id")
		return nil, errors.New("notification does not contain order_id")
	}

	s.Log.Infof("Checking transaction status for order ID: %s", orderID)

	// Get transaction status from Midtrans
	response, err := s.CoreAPIClient.CheckTransaction(orderID)
	if err != nil {
		s.Log.Errorf("Error checking transaction with Midtrans: %v", err)
		return nil, fmt.Errorf("error checking transaction: %w", err)
	}

	s.Log.Infof("Transaction status response from Midtrans: %+v", response)

	// Verify the transaction is valid
	if response.StatusCode != "200" {
		s.Log.Errorf("Invalid transaction status: %s", response.StatusMessage)
		return nil, fmt.Errorf("invalid transaction status: %s", response.StatusMessage)
	}

	s.Log.Info("Transaction successfully verified with Midtrans")
	return response, nil
}

// verifySignatureKey verifies the signature key from Midtrans notification
// The signature is generated using SHA512(order_id+status_code+gross_amount+ServerKey)
func (s *MidtransPaymentService) verifySignatureKey(notification map[string]interface{}) (bool, error) {
	// Extract required fields
	orderID, ok1 := notification["order_id"].(string)
	statusCode, ok2 := notification["status_code"].(string)
	grossAmount, ok3 := notification["gross_amount"].(string)
	signatureKey, ok4 := notification["signature_key"].(string)

	// Ensure all required fields are present
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return false, errors.New("notification missing required fields for signature verification")
	}

	// Get the server key from config
	serverKey := config.MidtransServerKey

	// Use the midtrans package for verification
	isValid, _ := midtransutils.VerifySignature(orderID, statusCode, grossAmount, signatureKey, serverKey)

	return isValid, nil
}

// Implementation of PaymentGateway interface methods
func (s *MidtransPaymentService) Charge(amount int, method string) (*PaymentResponse, error) {
	// This is a simplified implementation, in real scenarios we would create
	// a transaction and return its ID
	orderID := fmt.Sprintf("ORDER-%d-%d", amount, time.Now().Unix())

	return &PaymentResponse{
		TransactionID: orderID,
		Status:        "pending",
	}, nil
}

func (s *MidtransPaymentService) Refund(transactionID string) error {
	// Implementation would use Midtrans Core API to refund a transaction
	// For simplicity, this is not fully implemented
	return nil
}
