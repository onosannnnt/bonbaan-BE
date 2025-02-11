package transactionDriver

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	"gorm.io/gorm"
)

type transactionDriver struct {
	db *gorm.DB
}

func NewTransactionDriver(db *gorm.DB) orderUsecase.TransactionRepository {
	return &transactionDriver{
		db: db,
	}
}

func (d *transactionDriver) Insert(transaction *Entities.Transaction) error {
	if err := d.db.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}
