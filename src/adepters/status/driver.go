package statusAdapter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	"gorm.io/gorm"
)

type StatusDriver struct {
	db *gorm.DB
}

func NewStatusDriver(db *gorm.DB) statusUsecase.StatusRepository {
	return &StatusDriver{
		db: db,
	}
}

func (d *StatusDriver) Insert(status *Entities.Status) error {
	if err := d.db.Create(status).Error; err != nil {
		return err
	}
	return nil
}

func (d *StatusDriver) GetByID(id *string) (*Entities.Status, error) {
	var selectStatus Entities.Status
	if err := d.db.Where("id = ?", id).First(&selectStatus).Error; err != nil {
		return nil, err
	}
	return &selectStatus, nil
}

func (d *StatusDriver) GetByName(name *string) (*Entities.Status, error) {
	var selectStatus Entities.Status
	if err := d.db.Where("name = ?", name).First(&selectStatus).Error; err != nil {
		return nil, err
	}
	return &selectStatus, nil
}

func (d *StatusDriver) GetAll() ([]*Entities.Status, error) {
	var selectStatus []*Entities.Status
	if err := d.db.Find(&selectStatus).Error; err != nil {
		return nil, err
	}
	return selectStatus, nil
}

func (d *StatusDriver) Update(status *Entities.Status) error {
	if err := d.db.Model(status).Updates(status).Error; err != nil {
		return err
	}
	return nil
}

func (d *StatusDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", *id).Delete(&Entities.Status{}).Error; err != nil {
		return err
	}
	return nil
}
