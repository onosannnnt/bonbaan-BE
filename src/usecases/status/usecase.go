package statusUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type StatusUsecase interface {
	GetStatusByID(id *string) (*Entities.Status, error)
	GetStatusByName(name *string) (*Entities.Status, error)
	GetAll() ([]*Entities.Status, error)
	Insert(status *Entities.Status) error
	Update(status *Entities.Status) error
	Delete(id *string) error
}

type StatusService struct {
	statusRepo StatusRepository
}

func NewStatusService(repo StatusRepository) StatusUsecase {
	return &StatusService{
		statusRepo: repo,
	}
}

func (s *StatusService) GetStatusByID(id *string) (*Entities.Status, error) {
	return s.statusRepo.GetByName(id)
}

func (s *StatusService) GetStatusByName(name *string) (*Entities.Status, error) {
	return s.statusRepo.GetByName(name)
}

func (s *StatusService) GetAll() ([]*Entities.Status, error) {
	return s.statusRepo.GetAll()
}

func (s *StatusService) Insert(status *Entities.Status) error {
	return s.statusRepo.Insert(status)
}

func (s *StatusService) Update(status *Entities.Status) error {
	return s.statusRepo.Update(status)
}

func (s *StatusService) Delete(id *string) error {
	return s.statusRepo.Delete(id)
}
