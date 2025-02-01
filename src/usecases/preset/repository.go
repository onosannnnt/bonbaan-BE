package presetUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type PresetRepository interface {
	Insert(preset *Entities.Preset) error
	GetAll() (*[]Entities.Preset, error)
	GetByID(id *string) (*Entities.Preset, error)
	GetByServiceID(serviceID *string) (*[]Entities.Preset, error)
	Update(preset *Entities.Preset) error
	Delete(id *string) error
}
