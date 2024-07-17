package minecraftserver

import (
	"time"

	"github.com/google/uuid"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
)

type MinecraftServer struct {
	ID                    uuid.UUID            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	IP                    string               `gorm:"unique" validate:"required,fqdn" json:"ip,omitempty"`
	Port                  int                  `validate:"omitempty,number" json:"port,omitempty"`
	Slug                  string               `gorm:"unique" validate:"min=3,max=20" json:"slug,omitempty"`
	Name                  string               `validate:"required,min=3" json:"name,omitempty"`
	Logo                  *string              `validate:"omitempty,url" json:"logo,omitempty"`
	Description           *string              `validate:"omitempty,min=8" json:"description,omitempty"`
	Website               *string              `validate:"omitempty,url" json:"website,omitempty"`
	Discord               *string              `validate:"omitempty,url" json:"discord,omitempty"`
	AuthorID              string               `validate:"required,uuid" json:"author_id,omitempty"`
	Author                *user.User           `json:"author,omitempty"`
	MinecraftServerRconID *string              `validate:"omitempty,uuid" json:"minecraft_server_rcon_id,omitempty"`
	MinecraftServerRcon   *MinecraftServerRcon `json:"minecraft_server_rcon,omitempty"`
	CreatedAt             time.Time            `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt             time.Time            `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

type MinecraftServerRcon struct {
	ID                uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	IP                string           `gorm:"unique" validate:"required,fqdn" json:"ip,omitempty"`
	Port              int              `validate:"required,number" json:"port,omitempty"`
	Password          string           `validate:"required" json:"password,omitempty"`
	MinecraftServerID string           `gorm:"unique" validate:"required,uuid" json:"minecraft_server_id,omitempty"`
	MinecraftServer   *MinecraftServer `json:"minecraft_server,omitempty"`
}
