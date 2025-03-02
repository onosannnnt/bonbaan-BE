package orderTypeAdapter

import (
	"github.com/gofiber/fiber/v2"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	orderTypeUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order_type"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type OrderTypeHandler struct {
	orderTypeUsecase orderTypeUsecase.OrderTypeUsecase
}

func NewOrderTypeHandler(usecase orderTypeUsecase.OrderTypeUsecase) *OrderTypeHandler {
	return &OrderTypeHandler{
		orderTypeUsecase: usecase,
	}
}

func (h *OrderTypeHandler) Insert(c *fiber.Ctx) error {
	var orderType Entities.OrderType
	if err := c.BodyParser(&orderType); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the required fields", err, nil)
	}
	if err := h.orderTypeUsecase.Insert(&orderType); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "This package type already exists", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusCreated, "Success", nil, nil)
}

func (h *OrderTypeHandler) GetAll(c *fiber.Ctx) error {
	orderTypes, err := h.orderTypeUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, orderTypes)
}

func (h *OrderTypeHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	orderType, err := h.orderTypeUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, orderType)
}

func (h *OrderTypeHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var orderType Entities.OrderType
	if err := c.BodyParser(&orderType); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the required fields", err, nil)
	}
	if err := h.orderTypeUsecase.Update(&id, &orderType); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "This package type already exists", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *OrderTypeHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.orderTypeUsecase.Delete(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}
