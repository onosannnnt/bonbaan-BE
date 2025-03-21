package packageAdapter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	packageUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/package"
	"gorm.io/gorm"
)

type PackageDriver struct {
	db *gorm.DB
}

func NewPackageDriver(db *gorm.DB) packageUsecase.PackageRepository {
	return &PackageDriver{
		db: db,
	}
}

// Implement the Insert method to satisfy the ServiceRepository interface
func (d *PackageDriver) Insert(packages *Entities.Package) error {

	// fmt.Println("in driver ", Entities.Package)

	if err := d.db.Create(packages).Error; err != nil {
		return err
	}
	return nil
}

func (d *PackageDriver) GetAll() (*[]Entities.Package, error) {
	var packages []Entities.Package
	if err := d.db.Preload("OrderType").Find(&packages).Error; err != nil {
		return nil, err
	}
	return &packages, nil
}
func (d *PackageDriver) GetByID(id *string) (*Entities.Package, error) {
	var packages Entities.Package
	if err := d.db.Preload("OrderType").Where("id = ?", id).First(&packages).Error; err != nil {
		return nil, err
	}
	return &packages, nil
}

func (d *PackageDriver) GetByServiceID(serviceID *string) (*[]Entities.Package, error) {
	var packages []Entities.Package
	if err := d.db.Where("service_id = ?", serviceID).Find(&packages).Error; err != nil {
		return nil, err
	}
	return &packages, nil
}

func (d *PackageDriver) Update(packages *Entities.Package) error {
	if err := d.db.Model(&Entities.Package{}).Where("id = ?", packages.ID).Updates(packages).Error; err != nil {
		return err
	}
	return nil
}

func (d *PackageDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.Package{}).Error; err != nil {
		return err
	}
	return nil
}
