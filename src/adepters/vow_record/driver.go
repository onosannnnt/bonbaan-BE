package vowrecordAdapter

import (
	"fmt"

	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	vowrecordUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/vow_record"
	"gorm.io/gorm"
)

type VowRecordDriver struct {
	db *gorm.DB
}

func NewVowRecordDriver(db *gorm.DB) vowrecordUsecase.VowRecordRepository {
	return &VowRecordDriver{
		db: db,
	}
}

func (d *VowRecordDriver) Insert(vowRecord *Entities.VowRecord) error {
	if err := d.db.Create(vowRecord).Error; err != nil {
		return err
	}
	return nil
}

func (d *VowRecordDriver) GetAll(config *model.Pagination) ([]*Entities.VowRecord, int64, error) {
	var vowRecords []*Entities.VowRecord
	var count int64
	if err := d.db.Model(&Entities.VowRecord{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := d.db.Order(fmt.Sprintf("%s %s", config.OrderBy, config.OrderDirection)).Limit(config.PageSize).Offset((config.CurrentPage - 1) * config.PageSize).Find(&vowRecords).Error; err != nil {
		return nil, 0, err
	}
	return vowRecords, count, nil

}

func (d *VowRecordDriver) GetByID(id *string) (*Entities.VowRecord, error) {
	var vowRecord Entities.VowRecord
	if err := d.db.Where("id = ?", id).First(&vowRecord).Error; err != nil {
		return nil, err
	}
	return &vowRecord, nil
}

func (d *VowRecordDriver) GetByUserID(userID *string, config *model.Pagination) ([]*Entities.VowRecord, int64, error) {
	var vowRecords []*Entities.VowRecord
	var count int64
	if err := d.db.Model(&Entities.VowRecord{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := d.db.Where("user_id = ?", userID).Order(fmt.Sprintf("%s %s", config.OrderBy, config.OrderDirection)).Limit(config.PageSize).Offset((config.CurrentPage - 1) * config.PageSize).Find(&vowRecords).Error; err != nil {
		return nil, 0, err
	}
	return vowRecords, count, nil
}

func (d *VowRecordDriver) Update(id *string, vowRecord *Entities.VowRecord) error {
	if err := d.db.Model(&Entities.VowRecord{}).Where("id = ?", id).Updates(vowRecord).Error; err != nil {
		return err
	}
	return nil
}

func (d *VowRecordDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.VowRecord{}).Error; err != nil {
		return err
	}
	return nil
}
