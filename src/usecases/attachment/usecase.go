package attachmentUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type AttachmentUsecase interface {
    CreateAttachment(attachment *Entities.Attachment) error
    GetAll() (*[]Entities.Attachment, error)
    GetByServiceID(serviceID *string) (*[]Entities.Attachment, error) 
    GetByID(id *string) (*Entities.Attachment, error) 
	Update(attachment *Entities.Attachment) error 
	Delete(id *string) error 
}

type AttachmentService struct {
	attachmentRepo AttachmentRepository
}

func NewAttachmentService(repo AttachmentRepository) AttachmentUsecase {
    return &AttachmentService{
        attachmentRepo: repo,
    }
}

func (s *AttachmentService) CreateAttachment(attachment *Entities.Attachment) error {
	return s.attachmentRepo.Insert(attachment)
}

func (s *AttachmentService) GetAll() (*[]Entities.Attachment, error) {
	return s.attachmentRepo.GetAll()
}

func (s *AttachmentService) GetByServiceID(serviceID *string) (*[]Entities.Attachment, error) {
	return s.attachmentRepo.GetByServiceID(serviceID)
}

func (s *AttachmentService) GetByID(id *string) (*Entities.Attachment, error) {
	return s.attachmentRepo.GetByID(id)
}

func (s *AttachmentService) Update(attachment *Entities.Attachment) error {
	return s.attachmentRepo.Update(attachment)
}	

func (s *AttachmentService) Delete(id *string) error {
	return s.attachmentRepo.Delete(id)
}

