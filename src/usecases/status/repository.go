package statusUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type StatusRepository interface {
	FindStatusByID(id *string) (*Entities.Status, error)
	FindStatusByName(name *string) (*Entities.Status, error)
	FindAll() ([]*Entities.Status, error)
	Insert(status *Entities.Status) error
	Update(status *Entities.Status) error
	Delete(id *string) error
}
