package Entities

import (
	"errors"

	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/config"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Phone     string    `json:"phone"`
	RoleID    uuid.UUID `json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID ;references:ID"`
}

func InitializeUserData(db *gorm.DB) error {
	var role Role
	if err := db.Where("role = ?", constance.Admin_Role_ctx).First(&role).Error; err != nil {
		panic(err)
	}
	if role == (Role{}) {
		return errors.New("Role not found")
	}
	if config.AdminEmail == "" || config.AdminUsername == "" || config.AdminPassword == "" {
		return errors.New("admin data not found")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(config.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	admin := []User{
		{Username: config.AdminUsername, Password: string(hashPassword), Email: config.AdminEmail, Firstname: "Super", Lastname: "Admin", Phone: "0000000000", RoleID: role.ID},
	}
	for _, user := range admin {
		if err := db.FirstOrCreate(&user, User{Username: user.Username}).Error; err != nil {
			return err
		}
	}
	return nil
}
