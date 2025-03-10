package NotificationUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type NotificationUsecase interface {
	Insert(notification *Entities.Notification) error
	GetAll(config *model.Pagination) ([]*Entities.Notification, *model.Pagination, error)
	GetByID(id *string) (*Entities.Notification, error)
	Update(id *string, notification *Entities.Notification) error
	Delete(id *string) error
	Read(id *string) error
	GetByUserID(userID *string, config *model.Pagination) ([]*Entities.Notification, *model.Pagination, error)
}

type NotificationService struct {
	notificationRepo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) NotificationUsecase {
	return &NotificationService{
		notificationRepo: repo,
	}
}

func (s *NotificationService) Insert(notification *Entities.Notification) error {
	return s.notificationRepo.Insert(notification)
}

func (s *NotificationService) GetAll(config *model.Pagination) ([]*Entities.Notification, *model.Pagination, error) {
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
	notifications, totalRecords, err := s.notificationRepo.GetAll(config)
	if err != nil {
		return nil, nil, err
	}
	var totalPages int64
	if totalRecords%int64(config.PageSize) == 0 {
		totalPages = totalRecords / int64(config.PageSize)
	} else {
		totalPages = totalRecords/int64(config.PageSize) + 1
	}
	pagination := &model.Pagination{
		CurrentPage:  config.CurrentPage,
		PageSize:     config.PageSize,
		TotalPages:   int(totalPages),
		TotalRecords: int(totalRecords),
	}
	return notifications, pagination, nil
}

func (s *NotificationService) GetByID(id *string) (*Entities.Notification, error) {
	return s.notificationRepo.GetByID(id)
}

func (s *NotificationService) Update(id *string, notification *Entities.Notification) error {
	return s.notificationRepo.Update(id, notification)
}

func (s *NotificationService) Delete(id *string) error {
	return s.notificationRepo.Delete(id)
}

func (s *NotificationService) Read(id *string) error {
	return s.notificationRepo.Read(id)
}

func (s *NotificationService) GetByUserID(userID *string, config *model.Pagination) ([]*Entities.Notification, *model.Pagination, error) {
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
	notifications, totalRecords, err := s.notificationRepo.GetByUserID(userID, config)
	if err != nil {
		return nil, nil, err
	}
	var totalPages int64
	if totalRecords%int64(config.PageSize) == 0 {
		totalPages = totalRecords / int64(config.PageSize)
	} else {
		totalPages = totalRecords/int64(config.PageSize) + 1
	}
	pagination := &model.Pagination{
		CurrentPage:  config.CurrentPage,
		PageSize:     config.PageSize,
		TotalPages:   int(totalPages),
		TotalRecords: int(totalRecords),
	}
	return notifications, pagination, nil
}
