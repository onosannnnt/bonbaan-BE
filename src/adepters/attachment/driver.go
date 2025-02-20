package attachmentAdapter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	attachmentUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/attachment"
	"gorm.io/gorm"
)

type AttachmentDriver struct {
	db *gorm.DB
}

func NewAttachmentDriver(db *gorm.DB) attachmentUsecase.AttachmentRepository {
	return &AttachmentDriver{
		db: db,
	}
}

func (d *AttachmentDriver) Insert(attachment *Entities.Attachment) error {
	if err := d.db.Create(attachment).Error; err != nil {
		return err
	}
	return nil
}

func (d *AttachmentDriver) GetAll() (*[]Entities.Attachment, error) {
	var attachments *[]Entities.Attachment
	if err := d.db.Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (d *AttachmentDriver) GetByServiceID(serviceID *string) (*[]Entities.Attachment, error) {
	var attachments *[]Entities.Attachment
	if err := d.db.Where("service_id = ?", serviceID).Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (d *AttachmentDriver) GetByID(id *string) (*Entities.Attachment, error) {
	var attachment Entities.Attachment
	if err := d.db.Where("id = ?", id).First(&attachment).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (d *AttachmentDriver) Update(attachment *Entities.Attachment) error {
	if err := d.db.Model(&Entities.Attachment{}).Where("id = ?", attachment.ID).Updates(attachment).Error; err != nil {
		return err
	}
	return nil
}

func (d *AttachmentDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.Attachment{}).Error; err != nil {
		return err
	}
	return nil
}

