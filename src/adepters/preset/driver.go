package presetAdapter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	presetUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/preset"
	"gorm.io/gorm"
)

type  PresetDriver struct {
	db *gorm.DB
}

func NewPresetDriver(db *gorm.DB) presetUsecase.PresetRepository {
	return &PresetDriver{
		db: db,
	}
}

// Implement the Insert method to satisfy the ServiceRepository interface
func (d *PresetDriver) Insert(preset *Entities.Preset) error {
	if err := d.db.Create(preset).Error; err != nil {
		return err
	}
	return nil
}

func (d *PresetDriver) GetAll() (*[]Entities.Preset, error) {
	var preset []Entities.Preset
	if err := d.db.Find(&preset).Error; err != nil {
		return nil, err
	}
	return &preset, nil
}
func (d *PresetDriver) GetByID(id *string) (*Entities.Preset, error) {
	var preset Entities.Preset
	if err := d.db.Where("id = ?", id).First(&preset).Error; err != nil {
		return nil, err
	}
	return &preset, nil
}

func (d *PresetDriver) GetByServiceID(serviceID *string) (*[]Entities.Preset, error) {
	var preset []Entities.Preset
	if err := d.db.Where("service_id = ?", serviceID).Find(&preset).Error; err != nil {
		return nil, err
	}
	return &preset, nil
}




func (d *PresetDriver) Update(preset *Entities.Preset) error {
	if err := d.db.Model(&Entities.Preset{}).Where("id = ?", preset.ID).Updates(preset).Error; err != nil {
		return err
	}
	return nil
}

func (d *PresetDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.Preset{}).Error; err != nil {
		return err
	}
	return nil
}
