package reviewUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type ReviewRepository interface{
	Insert(review *Entities.Review) error
	GetAll() ([]*Entities.Review,error)
	GetByID(id *string)(*Entities.Review,error)
	Update(id *string, review *Entities.Review) error
	Delete(id *string) error	
}