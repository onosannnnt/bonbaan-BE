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
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	vowRecordUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/vow_record"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
	"google.golang.org/api/option"
)

type OrderHandler struct {
	OrderUsecase     orderUsecase.OrderUsecase
	VowRecordUsecase vowRecordUsecase.VowRecordService
}

func NewOrderHandler(usecase orderUsecase.OrderUsecase, vowRecordUsecase vowRecordUsecase.VowRecordService) *OrderHandler {
	return &OrderHandler{
		OrderUsecase:     usecase,
		VowRecordUsecase: vowRecordUsecase,
	}
}

func (h *OrderHandler) Insert(c *fiber.Ctx) error {
	order := model.OrderInputRequest{}
	if err := c.BodyParser(&order); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	order.UserID = c.Locals(constance.UserID_ctx).(string)
	h.OrderUsecase.Insert(&order)
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
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
	var order model.OrderInputRequest
	if err := c.BodyParser(&order); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	updateOrder := Entities.Order{}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format", err, nil)
	}
	updateOrder.ID = parsedID
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

func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var request struct {
		CancelReason string `json:"cancellation_reason"`
	}
	if err := c.BodyParser(&request); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	if err := h.OrderUsecase.CancelOrder(&id, &request.CancelReason); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to cancel order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) ApproveOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.OrderUsecase.ApproveOrder(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to accept order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) SubmitOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Error parsing form data", err, nil)
	}
	var input model.SubmitOrderRequest
	input.OrderID = id
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
		fmt.Println(err)
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

func (h *OrderHandler) InsertCustomOrder(c *fiber.Ctx) error {
	order := model.OrderInputRequest{}
	if err := c.BodyParser(&order); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	order.UserID = c.Locals(constance.UserID_ctx).(string)
	h.OrderUsecase.InsertCustomOrder(&order)
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) AcceptOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	updateOrder := model.ConfirmOrderRequest{
		OrderID: id,
	}
	if err := c.BodyParser(&updateOrder); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	if err := h.OrderUsecase.AcceptOrder(&updateOrder); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to confirm order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		userID = c.Locals(constance.UserID_ctx).(string)
	}
	config := model.Pagination{}
	status := c.Query("status")
	if status != "" {
		statusID, err := uuid.Parse(status)
		if err != nil {
			return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format for status", err, nil)
		}
		order, pagination, err := h.OrderUsecase.GetByUserIDAndStatusID(&userID, &statusID, &config)
		if err != nil {
			return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get order", err, nil)
		}
		return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, fiber.Map{
			"orders":     order,
			"pagination": pagination,
		})
	}
	if err := c.QueryParser(&config); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request query", err, nil)
	}
	order, pagination, err := h.OrderUsecase.GetByUserID(&userID, &config)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, fiber.Map{
		"orders":     order,
		"pagination": pagination,
	})
}

func (h *OrderHandler) GetByVowRecordByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		userID = c.Locals(constance.UserID_ctx).(string)
	}
	config := model.Pagination{}
	if err := c.QueryParser(&config); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request query", err, nil)
	}
	vowRecord, pagination, err := h.VowRecordUsecase.GetByUserID(&userID, &config)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, fiber.Map{
		"orders":     vowRecord,
		"pagination": pagination,
	})
}
