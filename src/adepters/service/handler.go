package serviceAdapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	ServiceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type ServiceHandler struct {
	ServiceUsecase ServiceUsecase.ServiceUsecase
}

func NewServiceHandler(ServiceUsecase ServiceUsecase.ServiceUsecase) *ServiceHandler {

	return &ServiceHandler{ServiceUsecase: ServiceUsecase}

}

func (h *ServiceHandler) CreateService(c *fiber.Ctx) error {
	var service Entities.Service

	if err := c.BodyParser(&service); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
	}

	// Parse categories from the request body
	var categories []Entities.Category
	if err := c.BodyParser(&categories); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid categories format", err, nil)
	}

	// Associate categories with the service
	service.Categories = categories

	if err := h.ServiceUsecase.CreateService(&service); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create service", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Service created successfully", nil, service)

}

func (h *ServiceHandler) GetAllServices(c *fiber.Ctx) error {
	services, err := h.ServiceUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, services)
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
