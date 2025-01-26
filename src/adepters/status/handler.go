package statusAdapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type StatusHandler struct {
	statusUsecase statusUsecase.StatusUsecase
}

func NewStatusHandler(statusUsecase statusUsecase.StatusUsecase) *StatusHandler {
	return &StatusHandler{
		statusUsecase: statusUsecase,
	}
}

func (h *StatusHandler) InsertStatus(c *fiber.Ctx) error {
	var status Entities.Status

	if err := c.BodyParser(&status); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the require fields", err, nil)
	}

	if err := h.statusUsecase.Insert(&status); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "this role already exists", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "success", nil, status)
}

func (h *StatusHandler) GetStatusByID(c *fiber.Ctx) error {
	id := c.Params("id")

	status, err := h.statusUsecase.GetStatusByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Status not found", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, status)
}

func (h *StatusHandler) GetStatusByName(c *fiber.Ctx) error {
	name := c.Params("name")

	status, err := h.statusUsecase.GetStatusByName(&name)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Status not found", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, status)
}

func (h *StatusHandler) GetAllStatus(c *fiber.Ctx) error {
	status, err := h.statusUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Status not found", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, status)
}

func (h *StatusHandler) UpdateStatus(c *fiber.Ctx) error {
	var id = c.Params("id")
	var status Entities.Status
	status.ID = uuid.MustParse(id)
	if err := c.BodyParser(&status); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the require fields", err, nil)
	}

	if err := h.statusUsecase.Update(&status); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "this role already exists", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, status)
}

func (h *StatusHandler) DeleteStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.statusUsecase.Delete(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Status not found", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, nil)
}
