package notificationAdepter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	notificantionUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/notification"
	"gorm.io/gorm"
)

type NotificationDriver struct {
	db *gorm.DB
}

func NewNotificationDriver(db *gorm.DB) notificantionUsecase.NotificationRepository {
	return &NotificationDriver{
		db: db,
	}
}

func (d *NotificationDriver) Insert(notification *Entities.Notification) error {
	if err := d.db.Create(notification).Error; err != nil {
		return err
	}
	return nil
}

func (d *NotificationDriver) GetAll(pagination *model.Pagination) ([]*Entities.Notification, int64, error) {
	var selectNotification []*Entities.Notification
	var count int64
	if err := d.db.Find(&selectNotification).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return selectNotification, count, nil
}

func (d *NotificationDriver) GetByID(id *string) (*Entities.Notification, error) {
	var selectNotification Entities.Notification
	if err := d.db.Where("id = ?", id).First(&selectNotification).Error; err != nil {
		return nil, err
	}
	return &selectNotification, nil
}

func (d *NotificationDriver) Update(id *string, notification *Entities.Notification) error {
	if err := d.db.Where("id = ?", id).Updates(notification).Error; err != nil {
		return err
	}
	return nil
}

func (d *NotificationDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.Notification{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *NotificationDriver) Read(id *string) error {
	if err := d.db.Model(&Entities.Notification{}).Where("id = ?", id).Update("is_read", true).Error; err != nil {
		return err
	}
	return nil
}

func (d *NotificationDriver) GetByUserID(userID *string, pagination *model.Pagination) ([]*Entities.Notification, int64, error) {
	var selectNotification []*Entities.Notification
	var count int64
	if err := d.db.Where("user_id = ?", &userID).Find(&selectNotification).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return selectNotification, count, nil
}
