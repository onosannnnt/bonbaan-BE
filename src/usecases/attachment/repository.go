package attachmentUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type AttachmentRepository interface {
	Insert(attachment *Entities.Attachment) error
	GetAll() (*[]Entities.Attachment, error)
	GetByServiceID(serviceID *string) (*[]Entities.Attachment, error)
	GetByID(id *string) (*Entities.Attachment, error)
	Update(attachment *Entities.Attachment) error
	Delete(id *string) error
}
