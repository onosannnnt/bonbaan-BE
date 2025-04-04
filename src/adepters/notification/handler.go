package notificationAdapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	notificationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/notification"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type NotificationHandler struct {
	NotificationUsecase notificationUsecase.NotificationUsecase
}

func NewNotificationHandler(usecase notificationUsecase.NotificationUsecase) NotificationHandler {
	return NotificationHandler{
		NotificationUsecase: usecase,
	}
}

func (h *NotificationHandler) Insert(c *fiber.Ctx) error {
	inputNotification := model.NotificationInsertRequest{}
	if err := c.BodyParser(&inputNotification); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	if inputNotification.UserID == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "UserID is required", nil, nil)
	}
	uuidUserID := uuid.MustParse(inputNotification.UserID)
	var uuidOrderID uuid.UUID
	if inputNotification.OrderID != "" {
		uuidOrderID = uuid.MustParse(inputNotification.OrderID)
	}
	notification := Entities.Notification{
		UserID:  uuidUserID,
		Header:  inputNotification.Header,
		Body:    inputNotification.Body,
		OrderID: uuidOrderID,
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", h.NotificationUsecase.Insert(&notification), nil)
}

func (h *NotificationHandler) GetAll(c *fiber.Ctx) error {
	config := model.Pagination{}
	if err := c.QueryParser(&config); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse query", err, nil)
	}
	notifications, pagination, err := h.NotificationUsecase.GetAll(&config)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get notifications", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, fiber.Map{
		"notifications": notifications,
		"pagination":    pagination,
	})
}

func (h *NotificationHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	notification, err := h.NotificationUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get notification", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, notification)
}

func (h *NotificationHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	notification := Entities.Notification{}
	if err := c.BodyParser(&notification); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", h.NotificationUsecase.Update(&id, &notification), nil)
}

func (h *NotificationHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "ID is required", nil, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", h.NotificationUsecase.Delete(&id), nil)
}

func (h *NotificationHandler) Read(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "ID is required", nil, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", h.NotificationUsecase.Read(&id), nil)
}

func (h *NotificationHandler) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")
	isReadProvided := c.Query("is-read") != ""
	isRead := c.QueryBool("is-read")
	if userID == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "UserID is required", nil, nil)
	}
	config := model.Pagination{}
	if err := c.QueryParser(&config); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse query", err, nil)
	}
	if isReadProvided {
		notifications, pagination, err := h.NotificationUsecase.GetUnreadByUserID(&userID, &isRead, &config)
		if err != nil {
			return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get notifications", err, nil)
		}
		return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, fiber.Map{
			"notifications": notifications,
			"pagination":    pagination,
		})
	}
	notifications, pagination, err := h.NotificationUsecase.GetByUserID(&userID, &config)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get notifications", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, fiber.Map{
		"notifications": notifications,
		"pagination":    pagination,
	})
}
