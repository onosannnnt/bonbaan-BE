package serviceAdapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/config"

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
	if (err != nil) {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Error parsing form data", err, nil)
	}
	// fmt.Println(form.Value)

	var input model.CreateServiceInput

	// Extract individual fields from form data.
	if v, exists := form.Value["name"]; exists && len(v) > 0 {
		input.Name = v[0]
	} else {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Name is missing", nil, nil)
	}
	if v, exists := form.Value["description"]; exists && len(v) > 0 {
		input.Description = v[0]
	} else {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Description is missing", nil, nil)
	}
	// Optional: parse rate if provided (assumed to be an integer)
	// Categories should be provided as multiple values.
	if v, exists := form.Value["categories"]; exists && len(v) > 0 {
		input.Categories = v
	} else {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Categories are missing", nil, nil)
	}
	if v, exists := form.Value["address"]; exists && len(v) > 0 {
		input.Address = v[0]
	} else {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Address is missing", nil, nil)
	}

	// Parse packages as a JSON string.
	if v, exists := form.Value["packages"]; exists && len(v) > 0 {
		// fmt.Println(v[0])
		// type of v[0]
		// fmt.Printf("%T\n", v[0])
		if err := json.Unmarshal([]byte(v[0]), &input.Packages); err != nil {
			return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid packages data", err, nil)
		}
	} else {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Packages data is missing", nil, nil)
	}

	// Create the service entity.
	service := Entities.Service{
		Name:        input.Name,
		Description: input.Description,
		Address:     input.Address,
	}

	// Map category IDs to category objects.
	for _, catID := range input.Categories {
		// fmt.Println(catID)
		uid, err := uuid.Parse(catID)
		if err != nil {
			return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid category id", err, nil)
		}
		service.Categories = append(service.Categories, Entities.Category{ID: uid})
	}

	// Map packages to the service.

	uniqueOrderTypeIDs := make(map[uuid.UUID]bool)

	for _, pkgInput := range input.Packages {
		orderTypeID, err := uuid.Parse(pkgInput.OrderTypeID)
		if err != nil {
			return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid type id", err, nil)
		}
		fmt.Println(orderTypeID)
		pkg := Entities.Package{
			Name:        pkgInput.Name,
			Item:        pkgInput.Item,
			Price:       pkgInput.Price,
			Description: pkgInput.Description,
			OrderTypeID: orderTypeID,
		}
		service.Packages = append(service.Packages, pkg)
		uniqueOrderTypeIDs[orderTypeID] = true
	}
	if v, exists := form.Value["custom_package"]; exists && len(v) > 0 {
		for orderTypeID := range uniqueOrderTypeIDs {
			customPackage := Entities.Package{
				Name:        "Custom Package",
				Description: "Custom package to your needs",
				Price:       0,
				OrderTypeID: orderTypeID,
			}
			service.Packages = append(service.Packages, customPackage)
		}
	}

	// Retrieve files from the "attachments" form field.
	files := form.File["attachments"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No attachments provided")
	}

	// Check if we are in test mode.
	if os.Getenv("TEST_MODE") == "true" {
		for range files {
			service.Attachments = append(service.Attachments, Entities.Attachment{URL: "http://dummy-url"})
		}
	} else {
		// Initialize the Cloud Storage client using the service account credentials.
		ctx := context.Background()
		client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.BucketKey))
		if err != nil {
			return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create storage client", err, nil)
		}
		defer client.Close()

		bucketName := config.BucketName
		shareableURLs := make([]string, 0, len(files))

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Error opening file", err, nil)
			}

			objectName := fmt.Sprintf("images/%d_%s", time.Now().UnixNano(), fileHeader.Filename)
			token := uuid.New().String()

			wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
			wc.Metadata = map[string]string{
				"firebaseStorageDownloadTokens": token,
			}

			if _, err = io.Copy(wc, file); err != nil {
				file.Close()
				wc.Close()
				return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to write file to bucket", err, nil)
			}
			file.Close()
			if err := wc.Close(); err != nil {
				return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to close writer", err, nil)
			}

			shareableURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
				bucketName, url.QueryEscape(objectName), token)
			shareableURLs = append(shareableURLs, shareableURL)
		}

		for _, imgURL := range shareableURLs {
			service.Attachments = append(service.Attachments, Entities.Attachment{URL: imgURL})
		}
	}

	// Create the service using the use case.
	if err := h.ServiceUsecase.CreateService(&service); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Success", nil, service)
}

func (h *ServiceHandler) GetAllServices(c *fiber.Ctx) error {
	config := model.Pagination{}
	if err := c.QueryParser(&config); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request query", err, nil)
	}

	services, pagination, err := h.ServiceUsecase.GetAll(&config)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	fmt.Println(services)
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, fiber.Map{
		"services":   services,
		"pagination": pagination,
	})
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
