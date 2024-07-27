package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Username  string    `gorm:"unique" validate:"required,min=3,max=20" json:"username,omitempty"`
	Email     string    `gorm:"unique" validate:"required,email" json:"email,omitempty"`
	Phone     string    `gorm:"unique" validate:"omitempty,e164" json:"phone,omitempty"`
	FullName  string    `validate:"required,min=3,max=50" json:"full_name,omitempty"`
	Password  string    `validate:"required,min=8" json:"password,omitempty"`
	Role      string    `gorm:"default:SELLER" validate:"omitempty,uppercase" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

func (u *User) ToResponse() *User {
	u.Password = ""
	return u
}
