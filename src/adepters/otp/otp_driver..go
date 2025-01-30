package otpDriver

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
	"gorm.io/gorm"
)

type otpDriver struct {
	db *gorm.DB
}

func NewOtpDriver(db *gorm.DB) userUsecase.OtpRepository {
	return &otpDriver{
		db: db,
	}
}

func (d *otpDriver) Insert(otp *Entities.Otp) error {
	if err := d.db.Create(otp).Error; err != nil {
		return err
	}
	return nil
}

func (d *otpDriver) GetByEmail(email *string, code *string) (*Entities.Otp, error) {
	var otp Entities.Otp
	err := d.db.Where("email = ? AND otp = ?", email, code).First(&otp).Error
	return &otp, err
}

func (d *otpDriver) DeleteByEmail(email *string) error {
	if err := d.db.Where("email = ?", email).Delete(&Entities.Otp{}).Error; err != nil {
		return err
	}
	return nil
}
