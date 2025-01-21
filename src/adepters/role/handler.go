package roleAdapter

import (
	"github.com/gofiber/fiber/v2"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	roleUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/role"
)

type RoleHandler struct {
	roleUsecase roleUsecase.RoleUsecase
}

func NewRoleHandler(roleUsecase roleUsecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{
		roleUsecase: roleUsecase,
	}
}

func (h *RoleHandler) InsertRole() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var role Entities.Role
		if err := c.BodyParser(&role); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Please fill all the require fields",
				"error":   err.Error(),
			})
		}
		if err := h.roleUsecase.InsertRole(role); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Role created successfully",
			"role":    role,
		})
	}
}

func (h *RoleHandler) GetAll() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		roles, err := h.roleUsecase.GetAll()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
			"roles":   roles,
		})
	}
}
