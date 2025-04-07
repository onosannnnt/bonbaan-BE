package recommendationUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type RecommendationRepository interface{
	Insert(review *Entities.Recommendation) error
	SuggestNextServie(userID string, config *model.Pagination) (*[]Entities.Service, int64, error)
	InterestRating(userID string,config *model.Pagination) (*[]Entities.Service, int64, error)
	Bestseller(config *model.Pagination) (*[]Entities.Service, int64, error)	
}