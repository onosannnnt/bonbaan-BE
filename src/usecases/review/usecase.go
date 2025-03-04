package reviewUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type ReviewUsecase interface {
	Insert(review *Entities.Review) error
	GetAll()([]*Entities.Review,error)
	GetByID(id *string)(*Entities.Review, error)
	Update(id *string, review *Entities.Review) error
	Delete(id *string) error
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

func (s *ReviewService) GetAll()([]*Entities.Review,error){
	return s.reviewRepo.GetAll()
}
func (s *ReviewService) GetByID(id *string)(*Entities.Review, error){
	return s.reviewRepo.GetByID(id)
}
func (s *ReviewService) Update(id *string,review *Entities.Review) error{
	return s.reviewRepo.Update(id, review)
}

func (s *ReviewService) Delete(id *string) error{
	return s.reviewRepo.Delete(id)
	}
