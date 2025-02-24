package serviceUsecase

import (
	"math"

	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type ServiceUsecase interface {
	CreateService(service *Entities.Service) error
	// Updated GetAll to return ServiceOutput models instead of Entities.Service.
	GetAll(config *model.Pagination) (*[]model.ServiceOutput, *model.Pagination, error)
	GetByID(id *string) (*Entities.Service, error)
	GetPackageByServiceID(serviceID *string) (*[]Entities.Package, error)
	UpdateService(service *Entities.Service) error
	DeleteService(id *string) error
}

type ServiceAsService struct {
	ServiceRepo ServiceRepository
}

func NewServiceUsecase(repo ServiceRepository) ServiceUsecase {
	return &ServiceAsService{ServiceRepo: repo}
}

func (sc *ServiceAsService) CreateService(service *Entities.Service) error {
	return sc.ServiceRepo.Insert(service)
}

// mapServiceToOutput converts an Entities.Service to a model.ServiceOutput.
func mapServiceToOutput(s Entities.Service) model.ServiceOutput {
    // Convert categories
    categories := make([]model.CategoryOutput, len(s.Categories))
    for i, c := range s.Categories {
        categories[i] = mapCategoryToOutput(c)
    }

    // Convert packages
    packages := make([]model.PackageOutput, len(s.Packages))
    for i, p := range s.Packages {
        packages[i] = mapPackageToOutput(p)
    }

    // Convert attachments
    attachments := make([]model.AttachmentOutput, len(s.Attachments))
    for i, a := range s.Attachments {
        attachments[i] = mapAttachmentToOutput(a)
    }

    return model.ServiceOutput{
        ID:          s.ID.String(),
        Name:        s.Name,
        Description: s.Description,
        Rate:        s.Rate,
        Adress:      s.Adress,
        Categories:  categories,
        Packages:    packages,
        Attachments: attachments,
    }
}

// mapCategoryToOutput converts an Entities.Category to a model.CategoryOutput.
func mapCategoryToOutput(c Entities.Category) model.CategoryOutput {
	return model.CategoryOutput{
		Name: c.Name,
	}
}

// mapPackageToOutput converts an Entities.Package to a model.PackageOutput.
func mapPackageToOutput(p Entities.Package) model.PackageOutput {
	return model.PackageOutput{
		Name:          p.Name,
		Item:          p.Item,
		Price:         p.Price,
		Description:   p.Description,
		PackageTypeID: p.PackageTypeID.String(),
	}
}

// mapAttachmentToOutput converts an Entities.Attachment to a model.AttachmentOutput.
func mapAttachmentToOutput(a Entities.Attachment) model.AttachmentOutput {
	return model.AttachmentOutput{
		URL: a.URL,
	}
}


// Updated GetAll method.
func (sc *ServiceAsService) GetAll(config *model.Pagination) (*[]model.ServiceOutput, *model.Pagination, error) {
	// Set default pagination values.
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}

	services, totalRecords, err := sc.ServiceRepo.GetAll(config)
	if err != nil {
		return nil, nil, err
	}

	totalPages := math.Ceil(float64(totalRecords) / float64(config.PageSize))
	pagination := &model.Pagination{
		CurrentPage:  config.CurrentPage,
		PageSize:     config.PageSize,
		TotalPages:   int(totalPages),
		TotalRecords: int(totalRecords),
	}

	// Map each Entities.Service to model.ServiceOutput.
	outputs := make([]model.ServiceOutput, 0, len(*services))
	for _, s := range *services {
		outputs = append(outputs, mapServiceToOutput(s))
	}

	return &outputs, pagination, nil
}

func (sc *ServiceAsService) GetByID(id *string) (*Entities.Service, error) {
	return sc.ServiceRepo.GetByID(id)
}

func (sc *ServiceAsService) GetPackageByServiceID(serviceID *string) (*[]Entities.Package, error) {
	return sc.ServiceRepo.GetPackagebyServiceID(serviceID)
}

func (sc *ServiceAsService) UpdateService(service *Entities.Service) error {
	return sc.ServiceRepo.Update(service)
}

func (sc *ServiceAsService) DeleteService(id *string) error {
	return sc.ServiceRepo.Delete(id)
}
