package categoryAdapter

import (
	"fmt"

	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	categoryUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/category"
	"gorm.io/gorm"
)

type Categorydriver struct {
	db *gorm.DB
}

func NewCategoryDriver(db *gorm.DB) categoryUsecase.CategoryRepository {
	return &Categorydriver{db: db,}
}

func (d *Categorydriver) Insert(category *Entities.Category) error {

	if err := d.db.Create(category).Error; err != nil {
		return err
	}
	return nil
}

func (d *Categorydriver) AddServiceToCategory(categoryID *string, serviceID *string) error {
    var Service_Category Entities.Service_Category

    // Convert *string to uuid.UUID
    catID, err := uuid.Parse(*categoryID)
    if err != nil {
        return err
    }
    servID, err := uuid.Parse(*serviceID)
    if err != nil {
        return err
    }

    Service_Category.CategoryID = catID
    Service_Category.ServiceID = servID
	fmt.Println(Service_Category) // Debugging line to print the Service_Category struct
	fmt.Println("Service_Category.CategoryID:", Service_Category.CategoryID)
	fmt.Println("Service_Category.ServiceID:", Service_Category.ServiceID)
	fmt.Println("Type of ervice_Category.CategoryID", fmt.Sprintf("%T", Service_Category.CategoryID))
	fmt.Println("Type of Service_Category.ServiceID", fmt.Sprintf("%T", Service_Category.ServiceID))


    if err := d.db.Create(&Service_Category).Error; err != nil {
        return err
    }
    return nil
}

func (d *Categorydriver) RemoveServiceFromCategory(categoryID *string, serviceID *string) error {
    // var Service_Category Entities.Service_Category

	catID, err := uuid.Parse(*categoryID)
	if err != nil {
		return err
	}
	servID, err := uuid.Parse(*serviceID)
	if err != nil {
		return err
	}
	fmt.Println(catID, servID) // Debugging line to print the Service_Category struct

	if err := d.db.Delete(&Entities.Service_Category{},"category_id = ? AND service_id = ?", catID, servID).Error; err != nil {
		return err
	}
	return nil
}

func (d *Categorydriver) GetByID(id *string) (*Entities.Category, error) {
	var category Entities.Category
	if err := d.db.First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (d *Categorydriver) GetServicesByCategoryID(categoryID *string) (*[]Entities.Service, error) {
	// Convert string to uuid.UUID
	catID, err := uuid.Parse(*categoryID)
	if err != nil {
		return nil, err
	}

	var services []Entities.Service
	if err := d.db.Joins("JOIN service_categories ON services.id = service_categories.service_id").
		Where("service_categories.category_id = ?", catID).
		Find(&services).Error; err != nil {
		return nil, err
	}
	return &services, nil
}

func (d *Categorydriver) GetCategoriesByServiceID(serviceID *string) (*[]Entities.Category, error) {
	// Convert string to uuid.UUID
	servID, err := uuid.Parse(*serviceID)
	if err != nil {
		return nil, err
	}
	var categories []Entities.Category
	if err := d.db.Joins("JOIN service_categories ON categories.id = service_categories.category_id").
		Where("service_categories.service_id = ?", servID).
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return &categories, nil
}

func (d *Categorydriver) GetAll() (*[]Entities.Category, error) {
	var categories []Entities.Category
	if err := d.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return &categories, nil
}

func (d *Categorydriver) Update(category *Entities.Category) error {
	if err := d.db.Save(category).Error; err != nil {
		return err
	}
	return nil
}

func (d *Categorydriver) Delete(id *string) error {
	if err := d.db.Delete(&Entities.Category{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
