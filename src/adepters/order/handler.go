package orderAdepter

import (
	"context"
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
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
	"google.golang.org/api/option"
)

type OrderHandler struct {
	OrderUsecase orderUsecase.OrderUsecase
}

func NewOrderHandler(usecase orderUsecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		OrderUsecase: usecase,
	}
}

func (h *OrderHandler) Insert(c *fiber.Ctx) error {
	order := model.OrderInsertRequest{}
	if err := c.BodyParser(&order); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	insertOrder := Entities.Order{}
	insertOrder.CancellationReason = order.CancellationReason
	insertOrder.OrderDetail = model.JSONB(order.OrderDetail)
	insertOrder.Note = order.Note
	parsedDate, err := time.Parse("2006-01-02", order.Deadline)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid date format", err, nil)
	}
	insertOrder.Deadline = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
	insertOrder.UserID, err = uuid.Parse(order.UserID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format for UserID", err, nil)
	}
	insertOrder.ServiceID, err = uuid.Parse(order.ServiceID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format for ServiceID", err, nil)
	}
	data, err := h.OrderUsecase.Insert(&insertOrder)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to insert order", err, err.Error())
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, data)
}

func (h *OrderHandler) GetAll(c *fiber.Ctx) error {
	status := c.Query("status")
	config := model.Pagination{}
	if err := c.QueryParser(&config); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request query", err, nil)
	}
	if status != "" {
		statusID, err := uuid.Parse(status)
		if err != nil {
			return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format for status", err, nil)
		}
		order, pagination, err := h.OrderUsecase.GetByStatus(&statusID, &config)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, fiber.Map{
			"orders":     order,
			"pagination": pagination,
		})
	}
	order, pagination, err := h.OrderUsecase.GetAll(&config)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, fiber.Map{
		"orders":     order,
		"pagination": pagination,
	})
}

func (h *OrderHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := h.OrderUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, order)
}

func (h *OrderHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	order := model.OrderInsertRequest{}
	if err := c.BodyParser(&order); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	updateOrder := Entities.Order{}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format", err, nil)
	}
	updateOrder.ID = parsedID
	updateOrder.StatusID = uuid.MustParse(order.StatusID)
	if err := h.OrderUsecase.Update(&id, &updateOrder); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to update order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.OrderUsecase.Delete(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to delete order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) Hook(c *fiber.Ctx) error {
	var event model.ChargeEvent

	if err := c.BodyParser(&event); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	if event.Data.Status == "failed" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Payment failed", nil, nil)
	}
	err := h.OrderUsecase.Hook(&event.Data.ID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to update order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) CancleOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var request struct {
		CancleReason string `json:"cancle_reason"`
	}
	if err := c.BodyParser(&request); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	if err := h.OrderUsecase.CancleOrder(&id, &request.CancleReason); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to cancel order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) AcceptOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.OrderUsecase.AcceptOrder(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to accept order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) SubmitOrder(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Error parsing form data", err, nil)
	}
	var input model.ConfirmOrderRequest
	if v, exists := form.Value["orderID"]; exists && len(v) > 0 {
		input.OrderID = v[0]
	} else {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Name is missing", nil, nil)
	}

	files := form.File["attachments"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No attachments provided")
	}

	if os.Getenv("TEST_MODE") == "true" {
		for range files {
			input.Attachments = append(input.Attachments, Entities.Attachment{URL: "http://dummy-url"})
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
			input.Attachments = append(input.Attachments, Entities.Attachment{URL: imgURL})
		}
	}
	if err := h.OrderUsecase.SubmitOrder(&input); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to confirm order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) CompleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.OrderUsecase.CompleteOrder(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to complete order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}
