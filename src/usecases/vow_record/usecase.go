package vowRecordUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type VowRecordService struct {
	vowRecordRepo VowRecordRepository
}

type VowrecordUsecase interface {
	Insert(vowRecord *Entities.VowRecord) error
	GetAll(config *model.Pagination) ([]*Entities.VowRecord, *model.Pagination, error)
	GetByID(id *string) (*Entities.VowRecord, error)
	GetByUserID(userID *string, config *model.Pagination) ([]*Entities.VowRecord, *model.Pagination, error)
	Update(id *string, vowRecord *Entities.VowRecord) error
	Delete(id *string) error
}

func NewVowRecordService(vowRecordRepo VowRecordRepository) *VowRecordService {
	return &VowRecordService{
		vowRecordRepo: vowRecordRepo,
	}
}

func (s *VowRecordService) Insert(vowRecord *Entities.VowRecord) error {
	return s.vowRecordRepo.Insert(vowRecord)
}

func (s *VowRecordService) GetAll(config *model.Pagination) ([]*Entities.VowRecord, *model.Pagination, error) {
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
	if config.OrderBy == "" {
		config.OrderBy = "created_at"
	}
	if config.OrderDirection == "" {
		config.OrderDirection = "asc"
	}
	vowRecords, totalRecords, err := s.vowRecordRepo.GetAll(config)
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
	return vowRecords, pagination, nil
}

func (s *VowRecordService) GetByID(id *string) (*Entities.VowRecord, error) {
	return s.vowRecordRepo.GetByID(id)
}

func (s *VowRecordService) GetByUserID(userID *string, config *model.Pagination) ([]*Entities.VowRecord, *model.Pagination, error) {
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
	if config.OrderBy == "" {
		config.OrderBy = "created_at"
	}
	if config.OrderDirection == "" {
		config.OrderDirection = "asc"
	}
	vowRecords, totalRecords, err := s.vowRecordRepo.GetByUserID(userID, config)
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
	return vowRecords, pagination, nil
}

func (s *VowRecordService) Update(id *string, vowRecord *Entities.VowRecord) error {
	return s.vowRecordRepo.Update(id, vowRecord)
}

func (s *VowRecordService) Delete(id *string) error {
	return s.vowRecordRepo.Delete(id)
}
