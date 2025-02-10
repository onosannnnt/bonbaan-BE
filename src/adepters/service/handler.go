package serviceAdapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/Config"

	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	ServiceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
	"google.golang.org/api/option"
)

type ServiceHandler struct {
	ServiceUsecase ServiceUsecase.ServiceUsecase
}

func NewServiceHandler(ServiceUsecase ServiceUsecase.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{ServiceUsecase: ServiceUsecase}
}

func (h *ServiceHandler) CreateService(c *fiber.Ctx) error {
	// Parse the multipart form:
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Error parsing form data", err, nil)
	}

	// Retrieve the JSON part from the form data
	jsonData := form.Value["json"]
	if len(jsonData) == 0 {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "JSON data is missing", nil, nil)
	}

	// Unmarshal the JSON data into the input struct
	var input model.CreateServiceInput
	if err := json.Unmarshal([]byte(jsonData[0]), &input); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid JSON data", err, nil)
	}

	// Create the service entity
	service := Entities.Service{
		Name:        input.Name,
		Description: input.Description,
		Rate:        input.Rate,
	}

	// Map category IDs to category objects
	for _, catID := range input.Categories {
		uid, err := uuid.Parse(catID)
		if err != nil {
			return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid category id", err, nil)
		}
		service.Categories = append(service.Categories, Entities.Category{ID: uid})
	}

	// Map packages to the service
	for _, pkgInput := range input.Packages {
		pkg := Entities.Package{
			Name:        pkgInput.Name,
			Item:        pkgInput.Item,
			Price:       pkgInput.Price,
			Description: pkgInput.Description,
		}
		service.Packages = append(service.Packages, pkg)
	}

	// Retrieve files from the "attachments" form field
	files := form.File["attachments"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No attachments provided")
	}

	// Initialize the Cloud Storage client using the service account credentials.
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(Config.BucketKey))
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create storage client", err, nil)
	}
	defer client.Close()

	// Set your bucket name.
	bucketName := Config.BucketName

	// Slice to hold shareable URLs for all images.
	shareableURLs := make([]string, 0, len(files))

	for _, fileHeader := range files {
		// Open the file.
		file, err := fileHeader.Open()
		if err != nil {
			return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Error opening file", err, nil)
		}

		// Generate a unique object name.
		objectName := fmt.Sprintf("images/%d_%s", time.Now().UnixNano(), fileHeader.Filename)

		// Generate a random download token.
		token := uuid.New().String()

		// Create a writer to upload the file to the storage bucket.
		wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
		// Set the metadata with the download token.
		wc.Metadata = map[string]string{
			"firebaseStorageDownloadTokens": token,
		}

		// Copy the file's content to Cloud Storage.
		if _, err = io.Copy(wc, file); err != nil {
			file.Close()
			wc.Close()
			return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to write file to bucket", err, nil)
		}
		file.Close()
		if err := wc.Close(); err != nil {
			return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to close writer", err, nil)
		}

		// Construct the shareable URL.
		shareableURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
			bucketName, url.QueryEscape(objectName), token)
		shareableURLs = append(shareableURLs, shareableURL)
	}

	// Associate the shareable URLs with the service entity.
	// Assuming the Service entity has an Attachments field of type []string.
	for _, img_url := range shareableURLs {
		atg := Entities.Attachment{URL: img_url}
		service.Attachments = append(service.Attachments, atg)
	}
	

	// Create the service using the use case.
	if err := h.ServiceUsecase.CreateService(&service); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Success", nil, service)
}

func (h *ServiceHandler) GetAllServices(c *fiber.Ctx) error {
	services, err := h.ServiceUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, services)
}

func (h *ServiceHandler) GetByServiceID(c *fiber.Ctx) error {
	id := c.Params("id")
	service, err := h.ServiceUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, service)
}

func (h *ServiceHandler) GetPackagesbyServiceID(c *fiber.Ctx) error {
	serviceID := c.Params("id")
	presets, err := h.ServiceUsecase.GetPackageByServiceID(&serviceID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, presets)
}

func (h *ServiceHandler) UpdateService(c *fiber.Ctx) error {
	id := c.Params("id")
	var service Entities.Service

	if err := c.BodyParser(&service); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
	}

	// Convert the id to uuid.UUID
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format", err, nil)
	}

	service.ID = uuidID // Ensure the ID is set to the one from the URL

	if err := h.ServiceUsecase.UpdateService(&service); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to update service", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Service updated successfully", nil, service)
}

func (h *ServiceHandler) DeleteService(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.ServiceUsecase.DeleteService(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}
