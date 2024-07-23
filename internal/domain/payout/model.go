package payout

import (
	"time"

	"github.com/google/uuid"
	"github.com/safatanc/blockstuff-api/internal/domain/transaction"
)

type Payout struct {
	ID                 uuid.UUID            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	PayoutTransactions []*PayoutTransaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" validate:"required,min=1" json:"payout_transactions"`
	// Status: WAITING_APPROVAL, APPROVED
	TransactionSubtotal int64     `validate:"omitempty,number,min=10000" json:"transaction_subtotal"`
	Fee                 int64     `validate:"omitempty,number" json:"fee"`
	Subtotal            int64     `validate:"omitempty,number,min=10000" json:"subtotal"`
	Status              string    `gorm:"default:WAITING_APPROVAL" validate:"omitempty,uppercase" json:"status"`
	PayoutProofImageUrl *string   `validate:"omitempty,url" json:"payout_proof_image_url,omitempty"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

type PayoutTransaction struct {
	ID            uuid.UUID                `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	PayoutID      string                   `json:"payout_id"`
	TransactionID string                   `gorm:"unique" validate:"uuid" json:"transaction_id"`
	Transaction   *transaction.Transaction `json:"transaction,omitempty"`
}
