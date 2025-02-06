package packageAdapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	packageUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/package"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type PackageHandler struct {
	PackageUsecase packageUsecase.PackageUsecase
}

func NewPackageHandler(PackageUsecase packageUsecase.PackageUsecase) *PackageHandler {

	return &PackageHandler{PackageUsecase: PackageUsecase}

}

func (h *PackageHandler) CreatePackage(c *fiber.Ctx) error {
	packages := Entities.Package{}
	if err := c.BodyParser(&packages); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the require fields", err, nil)
	}

	if err := h.PackageUsecase.CreatePackage(&packages); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "this role already exists", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "success", nil, packages)

}

func (h *PackageHandler) GetAllPackage(c *fiber.Ctx) error {
	packages, err := h.PackageUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, packages)
}

func (h *PackageHandler) GetByPackageID(c *fiber.Ctx) error {
	id := c.Params("id")
	packages, err := h.PackageUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, packages)
}

func (h *PackageHandler) GetByServiceID(c *fiber.Ctx) error {
	serviceID := c.Params("serviceID")
	packages, err := h.PackageUsecase.GetByServiceID(&serviceID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, packages)
}



func (h *PackageHandler) UpdatePackage(c *fiber.Ctx) error {
	id := c.Params("id")
	var packages Entities.Package

	if err := c.BodyParser(&packages); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
	}

	// Convert the id to uuid.UUID
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format", err, nil)
	}

	packages.ID = uuidID // Ensure the ID is set to the one from the URL

	if err := h.PackageUsecase.UpdatePackage(&packages); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to update service", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Service updated successfully", nil, packages)
}

func (h *PackageHandler) DeletePackage(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.PackageUsecase.DeletePackage(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}
