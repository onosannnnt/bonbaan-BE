package packageTypeUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type PackageTypeRepository interface {
	Insert(role *Entities.PackageType) error
	GetAll() (*[]Entities.PackageType, error)
	GetByID(id *string) (*Entities.PackageType, error)
	Update(id *string, role *Entities.PackageType) error
	Delete(id *string) error
}
