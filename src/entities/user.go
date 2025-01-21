package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Phone     string    `json:"phone"`
	RoleID    uuid.UUID `json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID ;references:ID"`
}
