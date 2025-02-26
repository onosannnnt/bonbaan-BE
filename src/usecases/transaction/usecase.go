package transactionUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type TransectionUsecase interface {
	Insert(transaction *Entities.Transaction) error
	GetAll() (*[]Entities.Transaction, error)
	GetByID(id string) (*Entities.Transaction, error)
	Update(transaction *Entities.Transaction) error
	Delete(id string) error
}

type TransectionService struct {
	transectionRepo TransactionRepository
}

func NewTransectionService(repo TransactionRepository) TransectionUsecase {
	return &TransectionService{
		transectionRepo: repo,
	}
}

func (s *TransectionService) Insert(transaction *Entities.Transaction) error {
	return s.transectionRepo.Insert(transaction)
}
func (s *TransectionService) GetAll() (*[]Entities.Transaction, error) {
	return s.transectionRepo.GetAll()
}

func (s *TransectionService) GetByID(id string) (*Entities.Transaction, error) {
	return s.transectionRepo.GetByID(id)
}

func (s *TransectionService) Update(transaction *Entities.Transaction) error {
	return s.transectionRepo.Update(transaction)
}

func (s *TransectionService) Delete(id string) error {
	return s.transectionRepo.Delete(id)
}
