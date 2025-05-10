package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionDetail stores all payment-related information from Midtrans
type TransactionDetail struct {
	ID                 uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()"`
	UserSubscriptionID uuid.UUID `gorm:"not null"`
	OrderID            string    `gorm:"size:100;index"`
	TransactionID      string    `gorm:"size:100"`
	TransactionStatus  string    `gorm:"size:50"`
	TransactionTime    time.Time
	StatusCode         string `gorm:"size:10"`
	StatusMessage      string
	PaymentType        string `gorm:"size:50"`
	GrossAmount        string `gorm:"size:20"`
	Currency           string `gorm:"size:10"`
	FraudStatus        string `gorm:"size:20"`

	// Credit Card specific fields
	MaskedCard             *string `gorm:"size:50"`
	CardType               *string `gorm:"size:20"`
	Bank                   *string `gorm:"size:20"`
	ApprovalCode           *string `gorm:"size:50"`
	ECI                    *string `gorm:"size:10"`
	ChannelResponseCode    *string `gorm:"size:10"`
	ChannelResponseMessage *string

	// Bank Transfer specific fields
	VANumbers       JSON    `gorm:"type:jsonb"`
	PermataVANumber *string `gorm:"size:50"`
	BillerCode      *string `gorm:"size:20"`
	BillKey         *string `gorm:"size:50"`
	PaymentAmounts  JSON    `gorm:"type:jsonb"`

	// Store specific fields
	Store       *string `gorm:"size:50"`
	PaymentCode *string `gorm:"size:100"`

	// E-wallet specific fields
	Issuer   *string `gorm:"size:50"`
	Acquirer *string `gorm:"size:50"`

	// Settlement info
	SettlementTime *time.Time

	// Raw response for debugging
	RawResponse JSON `gorm:"type:jsonb"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// JSON custom type for handling JSON data
type JSON json.RawMessage

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON("null")
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return errors.New("invalid scan source for JSON")
	}
	*j = JSON(s)
	return nil
}

// MarshalJSON returns the JSON value of the string
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON sets the string to be the supplied JSON
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null pointer exception")
	}
	*j = JSON(data)
	return nil
}

func (t *TransactionDetail) BeforeCreate(_ *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
