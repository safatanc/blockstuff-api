package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/safatanc/blockstuff-api/internal/domain/item"
)

type Transaction struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Code              string             `gorm:"unique" json:"code"`
	MinecraftUsername string             `validate:"required,min=3,max=16" json:"minecraft_username"`
	Email             string             `validate:"required,email" json:"email"`
	Phone             *string            `validate:"omitempty,e164" json:"phone"`
	TransactionItems  []*TransactionItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" validate:"required,min=1" json:"transaction_items,omitempty"`
	Subtotal          int64              `validate:"omitempty,number,min=100" json:"subtotal"`
	// Status: WAITING_PAYMENT, EXPIRED, PAID, PAYOUT_REQUEST, PAYOUT_COMPLETE
	Status     string    `gorm:"default:WAITING_PAYMENT" validate:"omitempty,uppercase" json:"status"`
	QrisString string    `json:"qris_string"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

type TransactionItem struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	TransactionID string     `validate:"omitempty,uuid" json:"transaction_id"`
	ItemID        string     `validate:"omitempty,uuid" json:"item_id"`
	Item          *item.Item `json:"item,omitempty"`
	Quantity      int        `gorm:"default:1" validate:"omitempty,number,min=1" json:"quantity"`
	Subtotal      int64      `validate:"omitempty,number,min=100" json:"subtotal"`
}
