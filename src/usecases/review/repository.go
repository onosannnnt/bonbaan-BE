package reviewUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type ReviewRepository interface{
	Insert(review *Entities.Review) error
}