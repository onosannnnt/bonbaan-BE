package packageUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type PackageUsecase interface {
	CreatePackage(packages *Entities.Package) error
	GetAll() (*[]Entities.Package, error)
	GetByID(id *string) (*Entities.Package, error)
	GetByServiceID(serviceID *string) (*[]Entities.Package, error) // Added method to get packagess by service ID
	UpdatePackage(service *Entities.Package) error
	DeletePackage(id *string) error
}

type PackageAsService struct {
	PackageRepo PackageRepository
}

func NewPackageUsecase(repo PackageRepository) PackageUsecase {
	return &PackageAsService{PackageRepo: repo}
}

func (sc *PackageAsService) CreatePackage(packages *Entities.Package) error {
	return sc.PackageRepo.Insert(packages)
}

// Implement the GetAll method to satisfy the ServiceUsecase interface
func (sc *PackageAsService) GetAll() (*[]Entities.Package, error) {
	// Implementation of GetAll method
	return sc.PackageRepo.GetAll()
}

func (sc *PackageAsService) GetByID(id *string) (*Entities.Package, error) {
	// Implementation of GetByID method
	return sc.PackageRepo.GetByID(id)
}

func (sc *PackageAsService) GetByServiceID(serviceID *string) (*[]Entities.Package, error) {
	// Implementation of GetByServiceID method
	return sc.PackageRepo.GetByServiceID(serviceID)
}


func (sc *PackageAsService) UpdatePackage(packages *Entities.Package) error {
	// Implementation of UpdateService method
	return sc.PackageRepo.Update(packages)
}

func (sc *PackageAsService) DeletePackage(id *string) error {
	// Implementation of DeleteService method
	return sc.PackageRepo.Delete(id)
}

