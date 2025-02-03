package resetpasswordDriver

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
	"gorm.io/gorm"
)

type resetPasswordDriver struct {
	db *gorm.DB
}

func NewOtpDriver(db *gorm.DB) userUsecase.ResetPasswordRepository {
	return &resetPasswordDriver{
		db: db,
	}
}

func (d *resetPasswordDriver) Insert(resetPassword *Entities.ResetPassword) error {
	if err := d.db.Create(resetPassword).Error; err != nil {
		return err
	}
	return nil
}

func (d *resetPasswordDriver) GetByToken(token *string) (*Entities.ResetPassword, error) {
	var selectResetPassword Entities.ResetPassword
	if err := d.db.Where("reset_password = ?", token).First(&selectResetPassword).Error; err != nil {
		return nil, err
	}
	return &selectResetPassword, nil
}

func (d *resetPasswordDriver) DeleteByEmail(email *string) error {
	if err := d.db.Where("email = ?", *email).Delete(&Entities.ResetPassword{}).Error; err != nil {
		return err
	}
	return nil
}
