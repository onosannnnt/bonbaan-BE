package serviceAdapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/config"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	"gorm.io/gorm"

	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	AttachmentUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/attachment"
	RecommendationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/recommendation"
	ServiceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
	"google.golang.org/api/option"
)

type ServiceHandler struct {
    ServiceUsecase        ServiceUsecase.ServiceUsecase
    AttachmentUsecase     AttachmentUsecase.AttachmentUsecase
    RecommendationUsecase RecommendationUsecase.RecommendationUsecase
    DB                    *gorm.DB
}

// Modify NewServiceHandler to accept the recommendation use case and DB.
func NewServiceHandler(svcUsecase ServiceUsecase.ServiceUsecase, attachUsecase AttachmentUsecase.AttachmentUsecase, recUsecase RecommendationUsecase.RecommendationUsecase, db *gorm.DB) *ServiceHandler {
    return &ServiceHandler{
        ServiceUsecase:        svcUsecase,
        AttachmentUsecase:     attachUsecase,
        RecommendationUsecase: recUsecase,
        DB:                    db,
    }
}

func (h *ServiceHandler) CreateService(c *fiber.Ctx) error {
	// Parse the multipart form:
	form, err := c.MultipartForm()
	if err != nil {
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
	log.Println(input.Categories)
	for _, catID := range input.Categories {
		log.Println(catID)
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
	reviewUtils := Entities.Review_utils{
		ID:            uuid.New(),
		TotalReviewer: 0,
		TotalRete:     0,
		ServiceID:     service.ID,
	}
	service.Review_utils = reviewUtils
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
	// fmt.Println(services)
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

    // Parse JSON fields for basic info and associations.
    var input model.UpdateServiceInput
    form, err := c.MultipartForm()
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Error parsing form data", err, nil)
    }

    // If packages are passed via form-data as JSON, parse them:
    if v, exists := form.Value["packages"]; exists && len(v) > 0 {
        if err := json.Unmarshal([]byte(v[0]), &input.Packages); err != nil {
            return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid packages data", err, nil)
        }
    }

    // Use BodyParser to fill in pointer fields (Name, Description, Address, etc.)
    if err := c.BodyParser(&input); err != nil {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
    }

    // Retrieve the existing service from the database.
    existingService, err := h.ServiceUsecase.GetByID(&id)
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Service not found", err, nil)
    }
    if existingService == nil {
        return utils.ResponseJSON(c, fiber.StatusNotFound, "Service not found", nil, nil)
    }

    // Update basic fields only if new data was provided.
	if input.Name != "" {
		existingService.Name = input.Name
	}
	if input.Description != "" {
		existingService.Description = input.Description
	}
	if input.Address != "" {
		existingService.Address = input.Address
	}
    // ---------------------------------------------------------------------------------------
    // 1) Update Categories (replace associations) only if new categories are provided.
    // ---------------------------------------------------------------------------------------
    if len(input.Categories) > 0 {
        updatedCategories := []Entities.Category{}
        for _, catID := range input.Categories {
            uid, err := uuid.Parse(catID)
            if err != nil {
                return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid category id", err, nil)
            }
            updatedCategories = append(updatedCategories, Entities.Category{ID: uid})
        }
        existingService.Categories = updatedCategories
    }

    // ---------------------------------------------------------------------------------------
    // 2) Update Packages (replace associations) only if new packages data is provided.
    // ---------------------------------------------------------------------------------------
    if len(input.Packages) > 0 {
        updatedPackages := []Entities.Package{}
        uniqueOrderTypeIDs := make(map[uuid.UUID]bool)
        for _, pkgInput := range input.Packages {
            orderTypeID, err := uuid.Parse(pkgInput.OrderTypeID)
            if err != nil {
                return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid package order type id", err, nil)
            }
            pkg := Entities.Package{
                Name:        pkgInput.Name,
                Item:        pkgInput.Item,
                Price:       pkgInput.Price,
                Description: pkgInput.Description,
                OrderTypeID: orderTypeID,
            }
            updatedPackages = append(updatedPackages, pkg)
            uniqueOrderTypeIDs[orderTypeID] = true
        }

        // If the user wants a "Custom Package"
        if input.CustomPackage {
            for orderTypeID := range uniqueOrderTypeIDs {
                customPkg := Entities.Package{
                    Name:        "Custom Package",
                    Description: "Custom package to your needs",
                    Price:       0,
                    OrderTypeID: orderTypeID,
                }
                updatedPackages = append(updatedPackages, customPkg)
            }
        }
        existingService.Packages = updatedPackages
    }

    // ---------------------------------------------------------------------------------------
    // 3) Handle Attachment updates:
    //    - If no new attachment data or attachment_ids are provided, keep the existing ones.
    // ---------------------------------------------------------------------------------------

    // Parse the list of attachment IDs the user wants to keep, if provided.
    keepIDs := []string{}
    if v, exists := form.Value["attachments"]; exists && len(v) > 0 {
        keepIDs = v
    }

    // 3.1) Remove attachments that are NOT in keepIDs.
    // (Delete from Cloud Storage + DB)
    if len(keepIDs) > 0 {
        for _, oldAttach := range existingService.Attachments {
            oldAttachID := oldAttach.ID.String()
            if !stringInSlice(oldAttachID, keepIDs) {
                if err := deleteAttachmentFromStorage(oldAttach.URL); err != nil {
                    fmt.Printf("[WARN] Failed to remove from storage: %v\n", err)
                }
                if err := h.AttachmentUsecase.Delete(&oldAttachID); err != nil {
                    fmt.Printf("[WARN] Failed to remove from DB: %v\n", err)
                }
            }
        }
        // 3.2) Build the final slice of attachments that are kept.
        keptAttachments := []Entities.Attachment{}
        for _, oldAttach := range existingService.Attachments {
            if stringInSlice(oldAttach.ID.String(), keepIDs) {
                keptAttachments = append(keptAttachments, oldAttach)
            }
        }
        existingService.Attachments = keptAttachments
    }

    // 3.3) Check for any new attachments in "attachment_newfiles" and add them.
    newAttachmentHeaders := form.File["attachment_newfiles"]
    if len(newAttachmentHeaders) > 0 {
        ctx := context.Background()
        client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.BucketKey))
        if err != nil {
            return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create storage client", err, nil)
        }
        defer client.Close()
        bucketName := config.BucketName
        for _, fileHeader := range newAttachmentHeaders {
            srcFile, err := fileHeader.Open()
            if err != nil {
                return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Error opening file", err, nil)
            }
            objectName := fmt.Sprintf("images/%d_%s", time.Now().UnixNano(), fileHeader.Filename)
            token := uuid.New().String()
            wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
            wc.Metadata = map[string]string{
                "firebaseStorageDownloadTokens": token,
            }
            if _, err := io.Copy(wc, srcFile); err != nil {
                srcFile.Close()
                wc.Close()
                return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to write file to bucket", err, nil)
            }
            srcFile.Close()
            if err := wc.Close(); err != nil {
                return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to close writer", err, nil)
            }
            shareableURL := fmt.Sprintf(
                "https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
                bucketName,
                url.QueryEscape(objectName),
                token,
            )
            newAttach := Entities.Attachment{
                ID:        uuid.New(),
                URL:       shareableURL,
                ServiceID: existingService.ID,
            }
            existingService.Attachments = append(existingService.Attachments, newAttach)
        }
    }

    // ---------------------------------------------------------------------------------------
    // 4) Persist the updated service.
    // ---------------------------------------------------------------------------------------
    if err := h.ServiceUsecase.UpdateService(existingService); err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to update service", err, nil)
    }

    return utils.ResponseJSON(c, fiber.StatusOK, "Service updated successfully", nil, existingService)
}
// stringInSlice is a small helper to see if a string is in a slice
func stringInSlice(str string, list []string) bool {
    for _, item := range list {
        if item == str {
            return true
        }
    }
    return false
}

// deleteAttachmentFromStorage deletes an object from Cloud Storage given its download URL.
func deleteAttachmentFromStorage(fileURL string) error {
    ctx := context.Background()
    client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.BucketKey))
    if err != nil {
        return err
    }
    defer client.Close()

    // parseObjectNameFromURL is your helper function from your existing code
    objectName, err := parseObjectNameFromURL(fileURL)
    if err != nil {
        return err
    }
    if objectName == "" {
        return fmt.Errorf("object name is empty")
    }

    bucketName := config.BucketName
    return client.Bucket(bucketName).Object(objectName).Delete(ctx)
}

func parseObjectNameFromURL(fileURL string) (string, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	prefix := fmt.Sprintf("/v0/b/%s/o/", config.BucketName)
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



func (h *ServiceHandler) DeleteService(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.ServiceUsecase.DeleteService(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}
func (h *ServiceHandler) RecommendService(c *fiber.Ctx) error {
    var pagination model.Pagination
    if err := c.QueryParser(&pagination); err != nil {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid pagination parameters", err, nil)
    }
    if pagination.PageSize <= 0 {
        pagination.PageSize = 10
    }
    if pagination.CurrentPage <= 0 {
        pagination.CurrentPage = 1
    }

    userID, ok := c.Locals(constance.UserID_ctx).(string)
    if !ok || userID == "" {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Unable to retrieve userID from token", nil, nil)
    }

    // Check if the user has any transactions.
    var txCount int64
    if err := h.DB.Model(&Entities.Transaction{}).Where("user_id = ?", userID).Count(&txCount).Error; err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to count transactions", err, nil)
    }

    var outputs *[]model.ServiceOutput
    var pag *model.Pagination
    var err error

    if txCount > 0 {
        // Use SuggestNextServie if the user has transaction records.
        outputs, pag, err = h.RecommendationUsecase.SuggestNextServies(userID, &pagination)
    } else {
        // Otherwise use InterestRating.
        outputs, pag, err = h.RecommendationUsecase.InterestRatings(userID, &pagination)
    }
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get service recommendations", err, nil)
    }

    return utils.ResponseJSON(c, fiber.StatusOK, "Service recommendations retrieved successfully", nil, fiber.Map{
        "services":   outputs,
        "pagination": pag,
    })
}

// Bestseller returns a paginated list of Bestseller services.
func (h *ServiceHandler) Bestseller(c *fiber.Ctx) error {
    var pagination model.Pagination
    if err := c.QueryParser(&pagination); err != nil {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid pagination parameters", err, nil)
    }
    if pagination.PageSize <= 0 {
        pagination.PageSize = 10
    }
    if pagination.CurrentPage <= 0 {
        pagination.CurrentPage = 1
    }

    outputs, pag, err := h.RecommendationUsecase.Bestsellers(&pagination)
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get bestseller services", err, nil)
    }

    return utils.ResponseJSON(c, fiber.StatusOK, "Bestseller services retrieved successfully", nil, fiber.Map{
        "services":   outputs,
        "pagination": pag,
    })
}