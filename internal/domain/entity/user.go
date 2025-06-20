package entity

import (
	"challenge-travel-api/internal/domain/enums"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	Email     string         `json:"email" gorm:"type:varchar(255);not null;unique_index"`
	Password  string         `json:"password" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;not null"`
	IsActive  bool           `json:"is_active" gorm:"type:boolean;not null;default:true"`
	Role      enums.UserType `json:"role" gorm:"type:user_type;not null"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt *time.Time     `json:"deleted_at" gorm:"type:timestamp"`
}
