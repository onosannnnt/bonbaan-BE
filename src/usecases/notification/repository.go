package NotificationUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type NotificationRepository interface {
	Insert(notification *Entities.Notification) error
	GetAll(config *model.Pagination) ([]*Entities.Notification, int64, error)
	GetByID(id *string) (*Entities.Notification, error)
	Update(id *string, notification *Entities.Notification) error
	Delete(id *string) error
	Read(id *string) error
	GetByUserID(userID *string, pagination *model.Pagination) ([]*Entities.Notification, int64, error)
}
