package packageUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type PackageRepository interface {
	Insert(preset *Entities.Package) error
	GetAll() (*[]Entities.Package, error)
	GetByID(id *string) (*Entities.Package, error)
	GetByServiceID(serviceID *string) (*[]Entities.Package, error)
	Update(preset *Entities.Package) error
	Delete(id *string) error
}
