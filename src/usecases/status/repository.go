package statusUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type StatusRepository interface {
	GetStatusByID(id *string) (*Entities.Status, error)
	GetStatusByName(name *string) (*Entities.Status, error)
	GetAll() ([]*Entities.Status, error)
	Insert(status *Entities.Status) error
	Update(status *Entities.Status) error
	Delete(id *string) error
}
