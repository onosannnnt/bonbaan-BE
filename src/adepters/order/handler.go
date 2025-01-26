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
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", nil, err)
	}
	insertOrder := Entities.Order{}
	insertOrder.CancellationReason = order.CancellationReason
	insertOrder.OrderDetail = model.JSONB(order.OrderDetail)
	insertOrder.Note = order.Note
	parsedDate, err := time.Parse("2006-01-02", order.Deadline)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid date format", nil, err)
	}
	insertOrder.Deadline = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
	insertOrder.UserID = uuid.MustParse(order.UserID)
	insertOrder.ServiceID = uuid.MustParse(order.ServiceID)
	if err := h.OrderUsecase.Insert(&insertOrder); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to insert order", err, err.Error())
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) GetAll(c *fiber.Ctx) error {
	config := model.OrderGetAll{}
	if err := c.QueryParser(&config); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse query", nil, err)
	}
	order, err := h.OrderUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get order", nil, err)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, order)
}

func (h *OrderHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := h.OrderUsecase.GetOne(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get order", nil, err)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, order)
}

func (h *OrderHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	order := model.OrderInsertRequest{}
	if err := c.BodyParser(&order); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to parse request body", nil, err)
	}
	updateOrder := Entities.Order{}
	updateOrder.ID = id
	updateOrder.StatusID = uuid.MustParse(order.StatusID)
	if err := h.OrderUsecase.Update(&id, &updateOrder); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to update order", nil, err)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.OrderUsecase.Delete(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to delete order", nil, err)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}
