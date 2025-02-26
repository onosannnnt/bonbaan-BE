package packageTypeAdapter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	packageTypeUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/packageType"
	"gorm.io/gorm"
)

type PackageTypeDriver struct {
    db *gorm.DB
}

func NewPackageTypeDriver(db *gorm.DB) packageTypeUsecase.PackageTypeRepository {
    return &PackageTypeDriver{
        db: db,
    }
}

func (d *PackageTypeDriver) Insert(packageType *Entities.PackageType) error {
    if err := d.db.Create(packageType).Error; err != nil {
        return err
    }
    return nil
}

func (d *PackageTypeDriver) GetAll() (*[]Entities.PackageType, error) {
    // Use a slice instead of a pointer to a slice for proper GORM handling.
    var packageTypes []Entities.PackageType
    if err := d.db.Find(&packageTypes).Error; err != nil {
        return nil, err
    }
    return &packageTypes, nil
}

func (d *PackageTypeDriver) GetByID(id *string) (*Entities.PackageType, error) {
	var packageType Entities.PackageType
	if err := d.db.First(&packageType, id).Error; err != nil {
		return nil, err
	}
	return &packageType, nil
}

func (d *PackageTypeDriver) Update(id *string, packageType *Entities.PackageType) error {
	if err := d.db.Model(&Entities.PackageType{}).Where("id = ?", id).Updates(packageType).Error; err != nil {
		return err
	}
	return nil
}

func (d *PackageTypeDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.PackageType{}).Error; err != nil {
		return err
	}
	return nil
}
