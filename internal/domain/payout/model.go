package payout

import (
	"time"

	"github.com/google/uuid"
	"github.com/safatanc/blockstuff-api/internal/domain/transaction"
)

type Payout struct {
	ID                 uuid.UUID            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	PayoutTransactions []*PayoutTransaction `validate:"required,min=1" json:"payout_transactions"`
	// Status: WAITING_APPROVAL, APPROVED
	Status              string    `validate:"uppercase" json:"status"`
	PayoutProofImageUrl *string   `validate:"url" json:"payout_proof_image_url,omitempty"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

type PayoutTransaction struct {
	ID            uuid.UUID                `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	TransactionID string                   `gorm:"unique" validate:"uuid" json:"transaction_id"`
	Transaction   *transaction.Transaction `json:"transaction,omitempty"`
}
