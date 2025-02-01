package presetAdapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	presetUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/preset"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type PresetHandler struct {
	PresetUsecase presetUsecase.PresetUsecase
}

func NewPresetHandler(PresetUsecase presetUsecase.PresetUsecase) *PresetHandler {

	return &PresetHandler{PresetUsecase: PresetUsecase}

}

func (h *PresetHandler) CreatePreset(c *fiber.Ctx) error {
	preset := Entities.Preset{}
	if err := c.BodyParser(&preset); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the require fields", err, nil)
	}

	if err := h.PresetUsecase.CreatePreset(&preset); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "this role already exists", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "success", nil, preset)

}

func (h *PresetHandler) GetAllPreset(c *fiber.Ctx) error {
	preset, err := h.PresetUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, preset)
}

func (h *PresetHandler) GetByPresetID(c *fiber.Ctx) error {
	id := c.Params("id")
	preset, err := h.PresetUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, preset)
}

func (h *PresetHandler) GetByServiceID(c *fiber.Ctx) error {
	serviceID := c.Params("serviceID")
	preset, err := h.PresetUsecase.GetByServiceID(&serviceID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, preset)
}



func (h *PresetHandler) UpdatePreset(c *fiber.Ctx) error {
	id := c.Params("id")
	var preset Entities.Preset

	if err := c.BodyParser(&preset); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
	}

	// Convert the id to uuid.UUID
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format", err, nil)
	}

	preset.ID = uuidID // Ensure the ID is set to the one from the URL

	if err := h.PresetUsecase.UpdatePreset(&preset); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to update service", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Service updated successfully", nil, preset)
}

func (h *PresetHandler) DeletePreset(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.PresetUsecase.DeletePreset(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}
