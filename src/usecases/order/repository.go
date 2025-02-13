package orderUsecase

import (
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type OrderRepository interface {
	Insert(order *Entities.Order) error
	GetAll() ([]*Entities.Order, error)
	GetByID(id *string) (*Entities.Order, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
	GetDefaultStatus() (*Entities.Status, error)
	GetStatusIDByName(name string) (uuid.UUID, error)
	GetAndUpdateByChargeID(chargeID string) error
}
