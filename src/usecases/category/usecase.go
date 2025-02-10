package categoryUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type CategoryUsecase interface {
	CreateCategory( category *Entities.Category) error
	// AddServiceToCategory(categoryID *string, serviceID *string) error
	// RemoveServiceFromCategory(categoryID *string, serviceID *string) error
	GetByID(id *string) (*Entities.Category, error)
	GetServicesByCategoryID(categoryID *string) (*[]Entities.Service, error)
	// GetCategoriesByServiceID(serviceID *string) (*[]Entities.Category, error)
	GetAll( ) (*[]Entities.Category, error)
	Update(  category *Entities.Category) error
	Delete(  id *string) error
}

type CategoryAsService struct {
	CategoryRepo CategoryRepository
}

func NewCategoryUsecase(repo CategoryRepository) CategoryUsecase {
	return &CategoryAsService{CategoryRepo: repo}
}

func (u *CategoryAsService) CreateCategory(category *Entities.Category) error {
	return u.CategoryRepo.Insert(category)
}

// func (u *CategoryAsService) AddServiceToCategory(categoryID *string, serviceID *string) error {
// 	return u.CategoryRepo.AddServiceToCategory(categoryID, serviceID)
// }

// func (u *CategoryAsService) RemoveServiceFromCategory(categoryID *string, serviceID *string) error {
// 	return u.CategoryRepo.RemoveServiceFromCategory(categoryID, serviceID)
// }



func (u *CategoryAsService) GetAll() (*[]Entities.Category, error) {
	return u.CategoryRepo.GetAll()
}

func (u *CategoryAsService) GetByID(id *string) (*Entities.Category, error) {
	return u.CategoryRepo.GetByID(id)
}


func (u *CategoryAsService) GetServicesByCategoryID(categoryID *string) (*[]Entities.Service, error) {
	return u.CategoryRepo.GetServicesByCategoryID(categoryID)
}


// func (u *CategoryAsService) GetCategoriesByServiceID(serviceID *string) (*[]Entities.Category, error) {
// 	return u.CategoryRepo.GetCategoriesByServiceID(serviceID)
// }

func (u *CategoryAsService) Update( category *Entities.Category) error {
	return u.CategoryRepo.Update(category)
}

func (u *CategoryAsService) Delete( id *string) error {
	return u.CategoryRepo.Delete(id)
}
