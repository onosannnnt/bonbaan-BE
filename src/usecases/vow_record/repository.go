package vowRecordUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type VowRecordRepository interface {
	Insert(vowRecord *Entities.VowRecord) error
	GetAll(config *model.Pagination) ([]*Entities.VowRecord, int64, error)
	GetByID(id *string) (*Entities.VowRecord, error)
	GetByUserID(userID *string, config *model.Pagination) ([]*Entities.VowRecord, int64, error)
	Update(id *string, vowRecord *Entities.VowRecord) error
	Delete(id *string) error
}
