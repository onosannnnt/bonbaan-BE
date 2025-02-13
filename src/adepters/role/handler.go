package roleAdapter

import (
	"github.com/gofiber/fiber/v2"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	roleUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/role"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type RoleHandler struct {
	roleUsecase roleUsecase.RoleUsecase
}

func NewRoleHandler(roleUsecase roleUsecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{
		roleUsecase: roleUsecase,
	}
}

func (h *RoleHandler) InsertRole(c *fiber.Ctx) error {
	var role Entities.Role
	if err := c.BodyParser(&role); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the require fields", err, nil)
	}
	if err := h.roleUsecase.InsertRole(&role); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "this role already exists", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusCreated, "success", nil, nil)
}

func (h *RoleHandler) GetAllRole(c *fiber.Ctx) error {

	roles, err := h.roleUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, roles)
}
