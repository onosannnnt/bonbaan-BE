package userAdepter

import (
	"errors"

	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
	"gorm.io/gorm"
)

// ส่วนที่ต่อกับ driver handler
type UserDriver struct {
	db *gorm.DB
}

// สร้าง instance ของ UserDriver
func NewUserDriver(db *gorm.DB) userUsecase.UserRepository {
	return &UserDriver{
		db: db,
	}
}

// ส่วนของการทำงานของ UserDriver
func (d *UserDriver) Insert(user *Entities.User) error {
	if err := d.db.Create(user).Error; err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"uni_users_username\" (SQLSTATE 23505)" {
			return errors.New("this account already exists")
		}
		return err
	}
	return nil
}

func (d *UserDriver) GetByEmailOrUsername(user *Entities.User) (*Entities.User, error) {
	var selectUser Entities.User
	if err := d.db.Preload("Role").Where("email = ? OR username = ?", user.Email, user.Username).First(&selectUser).Error; err != nil {
		return nil, err
	}
	return &selectUser, nil
}

func (d *UserDriver) GetByID(id *string) (*Entities.User, error) {
	var selectUser Entities.User
	if err := d.db.Omit("password").Preload("Role").Where("id = ?", id).First(&selectUser).Error; err != nil {
		return nil, err
	}
	return &selectUser, nil
}

func (d *UserDriver) Update(user *Entities.User) (*Entities.User, error) {
	if err := d.db.Model(user).Updates(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (d *UserDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", *id).Delete(&Entities.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserDriver) GetAll() (*[]Entities.User, error) {
	var users []Entities.User
	if err := d.db.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (d *UserDriver) GetInterestByUserID(id *string) (*Entities.User, error) {
	var user Entities.User
	if err := d.db.Preload("Category").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDriver) DeleteInterest(userID *string, categoryID *string) error {
	categoryUUID, err := uuid.Parse(*categoryID)
	if err != nil {
		return err
	}

	if err := d.db.Model(&Entities.User{ID: uuid.MustParse(*userID)}).Association("Category").Delete(&Entities.Category{ID: categoryUUID}); err != nil {
		return err
	}
	return nil
}
