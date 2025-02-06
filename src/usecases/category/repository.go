package categoryUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type CategoryRepository interface {
	Insert(category *Entities.Category) error
	// AddServiceToCategory(categoryID *string, serviceID *string) error
	
	GetByID(id *string) (*Entities.Category, error)
	// GetServicesByCategoryID(categoryID *string) (*[]Entities.Service, error)
	// GetCategoriesByServiceID(serviceID *string) (*[]Entities.Category, error)
	GetAll() (*[]Entities.Category, error)
	
	Update(category *Entities.Category) error
	
	Delete(id *string) error
	// RemoveServiceFromCategory(categoryID *string, serviceID *string) error


}
