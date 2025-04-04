package roleAdapter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	roleUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/role"
	"gorm.io/gorm"
)

type RoleDriver struct {
	db *gorm.DB
}

func NewRoleDriver(db *gorm.DB) roleUsecase.RoleRepository {
	return &RoleDriver{
		db: db,
	}
}

func (d *RoleDriver) Insert(role *Entities.Role) error {
	if err := d.db.Create(role).Error; err != nil {
		return err
	}
	return nil
}

func (d *RoleDriver) GetAll() (*[]Entities.Role, error) {
	var roles *[]Entities.Role
	if err := d.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (d *RoleDriver) GetByName(name *string) (*Entities.Role, error) {
	var role Entities.Role
	if err := d.db.Where("role = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (d *RoleDriver) GetByID(id *string) (*Entities.Role, error) {
	var role Entities.Role
	if err := d.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (d *RoleDriver) Update(role *Entities.Role) error {
	if err := d.db.Save(role).Error; err != nil {
		return err
	}
	return nil
}

func (d *RoleDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.Role{}).Error; err != nil {
		return err
	}
	return nil
}
