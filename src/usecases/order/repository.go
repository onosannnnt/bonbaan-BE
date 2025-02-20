package orderUsecase

import (
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type OrderRepository interface {
	Insert(order *Entities.Order) error
	GetAll(config *model.Pagination) ([]*Entities.Order, int64, error)
	GetByID(id *string) (*Entities.Order, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
	GetByStatusName(status *string, config *model.Pagination) ([]*Entities.Order, int64, error)
	GetDefaultStatus() (*Entities.Status, error)
	GetStatusIDByName(name string) (uuid.UUID, error)
	GetAndUpdateByChargeID(chargeID string) error
}
