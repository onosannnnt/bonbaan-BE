package presetUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type PresetUsecase interface {
	CreatePreset(preset *Entities.Preset) error
	GetAll() (*[]Entities.Preset, error)
	GetByID(id *string) (*Entities.Preset, error)
	GetByServiceID(serviceID *string) (*[]Entities.Preset, error) // Added method to get presets by service ID
	UpdatePreset(service *Entities.Preset) error
	DeletePreset(id *string) error
}

type PresetAsService struct {
	PresetRepo PresetRepository
}

func NewPresetUsecase(repo PresetRepository) PresetUsecase {
	return &PresetAsService{PresetRepo: repo}
}

func (sc *PresetAsService) CreatePreset(preset *Entities.Preset) error {
	return sc.PresetRepo.Insert(preset)
}

// Implement the GetAll method to satisfy the ServiceUsecase interface
func (sc *PresetAsService) GetAll() (*[]Entities.Preset, error) {
	// Implementation of GetAll method
	return sc.PresetRepo.GetAll()
}

func (sc *PresetAsService) GetByID(id *string) (*Entities.Preset, error) {
	// Implementation of GetByID method
	return sc.PresetRepo.GetByID(id)
}
func (sc *PresetAsService) GetByServiceID(serviceID *string) (*[]Entities.Preset, error) {
	// Implementation of GetByServiceID method
	return sc.PresetRepo.GetByServiceID(serviceID)
}


func (sc *PresetAsService) UpdatePreset(preset *Entities.Preset) error {
	// Implementation of UpdateService method
	return sc.PresetRepo.Update(preset)
}

func (sc *PresetAsService) DeletePreset(id *string) error {
	// Implementation of DeleteService method
	return sc.PresetRepo.Delete(id)
}

