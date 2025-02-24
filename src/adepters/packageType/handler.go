package packageTypeAdapter

import (
	"github.com/gofiber/fiber/v2"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	packageTypeUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/packageType"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type PackageTypeHandler struct {
    packageTypeUsecase packageTypeUsecase.PackageTypeUsecase
}

func NewPackageTypeHandler(usecase packageTypeUsecase.PackageTypeUsecase) *PackageTypeHandler {
    return &PackageTypeHandler{
        packageTypeUsecase: usecase,
    }
}

func (h *PackageTypeHandler) Insert(c *fiber.Ctx) error {
    var packageType Entities.PackageType
    if err := c.BodyParser(&packageType); err != nil {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the required fields", err, nil)
    }
    if err := h.packageTypeUsecase.Insert(&packageType); err != nil {
        return utils.ResponseJSON(c, fiber.StatusConflict, "This package type already exists", err, nil)
    }
    return utils.ResponseJSON(c, fiber.StatusCreated, "Success", nil, nil)
}

func (h *PackageTypeHandler) GetAll(c *fiber.Ctx) error {
    packageTypes, err := h.packageTypeUsecase.GetAll()
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
    }
    return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, packageTypes)
}

func (h *PackageTypeHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	packageType, err := h.packageTypeUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, packageType)
}

func (h *PackageTypeHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var packageType Entities.PackageType
	if err := c.BodyParser(&packageType); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the required fields", err, nil)
	}
	if err := h.packageTypeUsecase.Update(&id, &packageType); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "This package type already exists", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}

func (h *PackageTypeHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.packageTypeUsecase.Delete(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, nil)
}
