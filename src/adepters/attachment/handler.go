package attachmentAdapter

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/Config"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	attachmentUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/attachment"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
	"google.golang.org/api/option"
)

type AttachmentHandler struct {
	attachmentUsecase attachmentUsecase.AttachmentUsecase
}

func NewAttachmentHandler(attachmentUsecase attachmentUsecase.AttachmentUsecase) *AttachmentHandler {
	return &AttachmentHandler{
		attachmentUsecase: attachmentUsecase,
	}
}

func (h *AttachmentHandler) CreateAttachment(c *fiber.Ctx) error {
	
	serviceID := c.FormValue("service_id")
	if serviceID == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "service_id is required", nil, nil)
	}

	uuidServiceID, err := uuid.Parse(serviceID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "invalid service_id", err, nil)
	}

	
	fileHeader, err := c.FormFile("attachments")
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "No file provided", err, nil)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Error opening file", err, nil)
	}
	defer file.Close()

	objectName := fmt.Sprintf("images/%s_%d_%s", serviceID, time.Now().UnixNano(), fileHeader.Filename)
	

	token := uuid.New().String()

	
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(Config.BucketKey))
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create storage client", err, nil)
	}
	defer client.Close()

	bucketName := Config.BucketName

	
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	
	wc.Metadata = map[string]string{
		"firebaseStorageDownloadTokens": token,
	}

	
	if _, err = io.Copy(wc, file); err != nil {
		wc.Close()
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to write file", err, nil)
	}
	if err := wc.Close(); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to close writer", err, nil)
	}

	
	shareableURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
		bucketName, url.QueryEscape(objectName), token)

	attachment := Entities.Attachment{
		ID:  uuid.New(),
		URL: shareableURL,
		ServiceID: uuidServiceID,
	}

	
	if err := h.attachmentUsecase.CreateAttachment(&attachment); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create attachment", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Attachment created successfully", nil, attachment)
}


func (h *AttachmentHandler) GetAllAttachment(c *fiber.Ctx) error {
	attachments, err := h.attachmentUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get attachments", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, attachments)
}

func (h *AttachmentHandler) GetAttachmentByServiceID(c *fiber.Ctx) error {
	serviceID := c.Params("service_id")
	if serviceID == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "service_id is required", nil, nil)
	}

	attachments, err := h.attachmentUsecase.GetByServiceID(&serviceID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get attachments", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, attachments)
}

func (h *AttachmentHandler) GetAttachmentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "id is required", nil, nil)
	}

	attachment, err := h.attachmentUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get attachment", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, attachment)
}


func (h *AttachmentHandler) UpdateAttachment(c *fiber.Ctx) error {
    
    id := c.Params("id")
    if id == "" {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "id is required", nil, nil)
    }

    
    existingAttachment, err := h.attachmentUsecase.GetByID(&id)
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get existing attachment", err, nil)
    }

  
    fileHeader, err := c.FormFile("attachments")
    if (err != nil) {
        
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "No new file provided", err, nil)
    }

    newFile, err := fileHeader.Open()
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Error opening new file", err, nil)
    }
    defer newFile.Close()

 
    ctx := context.Background()
    client, err := storage.NewClient(ctx, option.WithCredentialsFile(Config.BucketKey))
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create storage client", err, nil)
    }
    defer client.Close()

    bucketName := Config.BucketName

	oldURL := existingAttachment.URL
	objectName, err := parseObjectNameFromURL(oldURL)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "failed to parse object name from URL", err, nil)
	}
	if err := client.Bucket(bucketName).Object(objectName).Delete(ctx); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "failed to delete object", err, nil)
	}
	_ = client.Bucket(bucketName).Object(objectName).Delete(ctx)

	
    newObjectName := fmt.Sprintf("images/%s_%d_%s", 
        c.FormValue("service_id"), 
        time.Now().UnixNano(), 
        fileHeader.Filename,
    )
    newToken := uuid.New().String()

    wc := client.Bucket(bucketName).Object(newObjectName).NewWriter(ctx)
    wc.Metadata = map[string]string{
        "firebaseStorageDownloadTokens": newToken,
    }

    _, err = io.Copy(wc, newFile)
    if err != nil {
        wc.Close()
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to write file to storage", err, nil)
    }
    if err := wc.Close(); err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to close writer", err, nil)
    }


    newShareableURL := fmt.Sprintf(
        "https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
        bucketName, 
        url.QueryEscape(newObjectName), 
        newToken,
    )


    existingAttachment.URL = newShareableURL


    if err := h.attachmentUsecase.Update(existingAttachment); err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to update attachment in DB", err, nil)
    }


    return utils.ResponseJSON(c, fiber.StatusOK, "Attachment updated successfully", nil, existingAttachment)
}


func (h *AttachmentHandler) DeleteAttachment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "id is required", nil, nil)
	}

	attachment, err := h.attachmentUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get attachment", err, nil)
	}
	if attachment == nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Attachment not found", nil, nil)
	}

	
	objectName, err := parseObjectNameFromURL(attachment.URL)
	if err != nil || objectName == "" {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to parse object name from URL", err, nil)
	}

	
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(Config.BucketKey))
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create storage client", err, nil)
	}
	defer client.Close()

	bucketName := Config.BucketName

	
	if err := client.Bucket(bucketName).Object(objectName).Delete(ctx); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to delete file from storage", err, nil)
	}

	
	if err := h.attachmentUsecase.Delete(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to delete attachment record", err, nil)
	}

	
	return utils.ResponseJSON(c, fiber.StatusOK, "Attachment deleted successfully", nil, nil)
}


func parseObjectNameFromURL(fileURL string) (string, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	
	prefix := fmt.Sprintf("/v0/b/%s/o/", Config.BucketName)
	if !strings.HasPrefix(u.Path, prefix) {
		return "", fmt.Errorf("unexpected URL format")
	}

	
	encodedObjectName := strings.TrimPrefix(u.Path, prefix)

	
	objectName, err := url.QueryUnescape(encodedObjectName)
	if err != nil {
		return "", err
	}

	return objectName, nil
}
