package example

import (
	"time"

	"github.com/google/uuid"
)

// TransactionDetailResponse is a Swagger-friendly version of model.TransactionDetail
type TransactionDetailResponse struct {
	ID                 uuid.UUID `json:"id"`
	UserSubscriptionID uuid.UUID `json:"user_subscription_id"`
	OrderID            string    `json:"order_id"`
	TransactionID      string    `json:"transaction_id"`
	TransactionStatus  string    `json:"transaction_status"`
	TransactionTime    time.Time `json:"transaction_time"`
	StatusCode         string    `json:"status_code"`
	StatusMessage      string    `json:"status_message"`
	PaymentType        string    `json:"payment_type"`
	GrossAmount        string    `json:"gross_amount"`
	Currency           string    `json:"currency"`
	FraudStatus        string    `json:"fraud_status"`

	// Credit Card specific fields
	MaskedCard             *string `json:"masked_card,omitempty"`
	CardType               *string `json:"card_type,omitempty"`
	Bank                   *string `json:"bank,omitempty"`
	ApprovalCode           *string `json:"approval_code,omitempty"`
	ECI                    *string `json:"eci,omitempty"`
	ChannelResponseCode    *string `json:"channel_response_code,omitempty"`
	ChannelResponseMessage *string `json:"channel_response_message,omitempty"`

	// Bank Transfer specific fields
	VANumbers       JSONData `json:"va_numbers,omitempty"`
	PermataVANumber *string  `json:"permata_va_number,omitempty"`
	BillerCode      *string  `json:"biller_code,omitempty"`
	BillKey         *string  `json:"bill_key,omitempty"`
	PaymentAmounts  JSONData `json:"payment_amounts,omitempty"`

	// Store specific fields
	Store       *string `json:"store,omitempty"`
	PaymentCode *string `json:"payment_code,omitempty"`

	// E-wallet specific fields
	Issuer   *string `json:"issuer,omitempty"`
	Acquirer *string `json:"acquirer,omitempty"`

	// Settlement info
	SettlementTime *time.Time `json:"settlement_time,omitempty"`

	// Raw response for debugging
	RawResponse JSONData `json:"raw_response,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// TransactionsResponse is a Swagger-friendly version of response.SuccessWithTransactions
type TransactionsResponse struct {
	Status       string                      `json:"status"`
	Message      string                      `json:"message"`
	Data         []TransactionDetailResponse `json:"data"`
	Page         int                         `json:"page,omitempty"`
	Limit        int                         `json:"limit,omitempty"`
	TotalPages   int64                       `json:"total_pages,omitempty"`
	TotalResults int64                       `json:"total_results,omitempty"`
}
