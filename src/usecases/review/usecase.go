package reviewUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type ReviewUsecase interface {
	Insert(review *Entities.Review) error
}

type ReviewService struct {
	reviewRepo ReviewRepository
}

func NewReviewService(repo ReviewRepository) ReviewUsecase{
	return &ReviewService{
		reviewRepo: repo,
	}
}

func (s *ReviewService) Insert(review *Entities.Review) error{
	return s.reviewRepo.Insert(review)
	}
