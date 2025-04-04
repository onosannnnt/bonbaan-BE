package transactionUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type TransactionRepository interface {
	Insert(transaction *Entities.Transaction) error
	GetAll() (*[]Entities.Transaction, error)
	GetByID(id string) (*Entities.Transaction, error)
	Update(transaction *Entities.Transaction) error
	Delete(id string) error
}
