package orderUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type OrderRepository interface {
	Insert(order *Entities.Order) error
	GetAll() ([]*Entities.Order, error)
	GetByID(id *string) (*Entities.Order, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
	GetDefaultStatus() (*Entities.Status, error)
}

type TransactionRepository interface {
	Insert(transaction *Entities.Transaction) error
}
