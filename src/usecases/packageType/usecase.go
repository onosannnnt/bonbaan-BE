package packageTypeUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type PackageTypeUsecase interface {
    Insert(packageType *Entities.PackageType) error
    GetAll() (*[]Entities.PackageType, error)
    GetByID(id *string) (*Entities.PackageType, error)
    Update(id *string, packageType *Entities.PackageType) error
    Delete(id *string) error
}

type PackageTypeService struct {
    packageTypeRepo PackageTypeRepository
}

func NewPackageTypeService(repo PackageTypeRepository) PackageTypeUsecase {
    return &PackageTypeService{
        packageTypeRepo: repo,
    }
}

func (s *PackageTypeService) Insert(packageType *Entities.PackageType) error {
    return s.packageTypeRepo.Insert(packageType)
}

func (s *PackageTypeService) GetAll() (*[]Entities.PackageType, error) {
    return s.packageTypeRepo.GetAll()
}

func (s *PackageTypeService) GetByID(id *string) (*Entities.PackageType, error) {
    return s.packageTypeRepo.GetByID(id)
}

func (s *PackageTypeService) Update(id *string, packageType *Entities.PackageType) error {
    return s.packageTypeRepo.Update(id, packageType)
}

func (s *PackageTypeService) Delete(id *string) error {
    return s.packageTypeRepo.Delete(id)
}