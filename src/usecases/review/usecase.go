package reviewUsecase

import (
	"fmt"

	"github.com/onosannnnt/bonbaan-BE/src/constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
)

type ReviewUsecase interface {
	Insert(review *Entities.Review) error
	GetAll() ([]*Entities.Review, error)
	GetByID(id *string) (*Entities.Review, error)
	Update(id *string, review *Entities.Review) error
	Delete(id *string) error
}

type ReviewService struct {
	reviewRepo ReviewRepository
	orderRepo  orderUsecase.OrderRepository
	statusRepo statusUsecase.StatusRepository
}

func NewReviewService(repo ReviewRepository, orderRepo orderUsecase.OrderRepository, statusRepo statusUsecase.StatusRepository) ReviewUsecase {
	return &ReviewService{
		reviewRepo: repo,
		orderRepo:  orderRepo,
		statusRepo: statusRepo,
	}
}

func (s *ReviewService) Insert(review *Entities.Review) error {
	status, err := s.statusRepo.GetByName(&constance.Status_Completed)
	if err != nil {
		return err
	}
	orderIDStr := review.OrderID.String()
	order, err := s.orderRepo.GetByID(&orderIDStr)
	if err != nil {
		return err
	}
	order.Status = *status
	order.StatusID = status.ID
	if err := s.orderRepo.Update(&orderIDStr, order); err != nil {
		return err
	}
	fmt.Println(order.Status)
	fmt.Println(status)
	return s.reviewRepo.Insert(review)
}

func (s *ReviewService) GetAll() ([]*Entities.Review, error) {
	return s.reviewRepo.GetAll()
}
func (s *ReviewService) GetByID(id *string) (*Entities.Review, error) {
	return s.reviewRepo.GetByID(id)
}
func (s *ReviewService) Update(id *string, review *Entities.Review) error {
	return s.reviewRepo.Update(id, review)
}

func (s *ReviewService) Delete(id *string) error {
	return s.reviewRepo.Delete(id)
}
