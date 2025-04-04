package transactionAdepter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	transactionUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/transaction"
	"gorm.io/gorm"
)

type transactionDriver struct {
	db *gorm.DB
}

func NewTransactionDriver(db *gorm.DB) transactionUsecase.TransactionRepository {
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

func (d *transactionDriver) GetAll() (*[]Entities.Transaction, error) {
	var transactions []Entities.Transaction
	if err := d.db.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return &transactions, nil
}

func (d *transactionDriver) GetByID(id string) (*Entities.Transaction, error) {
	var transaction Entities.Transaction
	if err := d.db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (d *transactionDriver) Update(transaction *Entities.Transaction) error {
	if err := d.db.Save(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (d *transactionDriver) Delete(id string) error {
	if err := d.db.Delete(&Entities.Transaction{}, id).Error; err != nil {
		return err
	}
	return nil
}
