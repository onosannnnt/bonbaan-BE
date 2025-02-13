package orderAdepter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
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
	config := model.OrderGetAll{}
	if err := c.QueryParser(&config); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse query", err, nil)
	}
	order, err := h.OrderUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get order", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, order)
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
