package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSubscription struct {
	ID            uuid.UUID        `gorm:"primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID        `gorm:"not null"`
	PlanID        uuid.UUID        `gorm:"not null"`
	Plan          SubscriptionPlan `gorm:"foreignKey:PlanID"`
	AIscansUsed   int              `gorm:"default:0"`
	StartDate     time.Time        `gorm:"not null"`
	EndDate       time.Time        `gorm:"not null"`
	IsActive      bool             `gorm:"default:true"`
	PaymentMethod string           `gorm:"size:50"`
	TransactionID string           `gorm:"size:100"`
	PaymentStatus string           `gorm:"size:50;default:'pending'"`
	CreatedAt     time.Time        `gorm:"autoCreateTime"`
}

type PurchaseSubscriptionRequest struct {
	PaymentMethod string `json:"payment_method" validate:"omitempty,oneof=gopay shopeepay bank_transfer credit_card"`
}

type MidtransCallbackPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

type PaymentResponse struct {
	TransactionToken string `json:"transaction_token"`
	RedirectURL      string `json:"redirect_url"`
	OrderID          string `json:"order_id"`
}

func (userSubscription *UserSubscription) BeforeCreate(_ *gorm.DB) error {
	userSubscription.ID = uuid.New()
	return nil
}
