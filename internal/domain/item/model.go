package item

import (
	"time"

	"github.com/google/uuid"
	"github.com/safatanc/blockstuff-api/internal/domain/minecraftserver"
)

type Item struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Name        string    `validate:"omitempty,min=3,max=20" json:"name"`
	Slug        string    `validate:"omitempty,min=3,max=20" json:"slug"`
	Description *string   `validate:"omitempty,min=8" json:"description"`
	Price       int64     `validate:"omitempty,number,min=100" json:"price"`
	// Category: RANK, COIN, PERKS
	Category          string                           `validate:"omitempty,uppercase" json:"category"`
	MinecraftServerID *string                          `validate:"omitempty,uuid" json:"minecraft_server_id"`
	MinecraftServer   *minecraftserver.MinecraftServer `json:"minecraft_server,omitempty"`
	ItemImages        []*ItemImage                     `json:"item_images,omitempty"`
	ItemActions       []*ItemAction                    `json:"item_actions,omitempty"`
	Visible           bool                             `gorm:"default:true" json:"visible"`
	CreatedAt         time.Time                        `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt         time.Time                        `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

type ItemImage struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	ObjectName string    `json:"object_name,omitempty"`
	URL        string    `validate:"omitempty,url" json:"url"`
	Position   int       `validate:"omitempty,number,min=1" json:"position"`
	ItemID     string    `validate:"omitempty,uuid" json:"item_id"`
	Item       *Item     `json:"item,omitempty"`
}

type ItemAction struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	// Type: COMMAND, VOUCHER
	Type        string `validate:"omitempty,uppercase" json:"type"`
	Action      string `json:"action"`
	Position    int    `validate:"omitempty,number,min=1" json:"position"`
	Description string `validate:"omitempty,min=8" json:"description"`
	ItemID      string `validate:"omitempty,uuid" json:"item_id"`
	Item        *Item  `json:"item,omitempty"`
}
