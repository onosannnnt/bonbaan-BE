package orderUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type OrderRepository interface {
	Insert(order *Entities.Order) error
	FindAll(page *int, count *int) ([]*Entities.Order, error)
	FindOne(id *string) (*Entities.Order, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
}
