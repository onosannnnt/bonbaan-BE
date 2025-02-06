package categoryAdapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	categoryUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/category"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type Handler struct {
	usecase categoryUsecase.CategoryUsecase
}

func NewHandler(usecase categoryUsecase.CategoryUsecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	category := Entities.Category{}
	if err := c.BodyParser(&category); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
	}

	if err := h.usecase.CreateCategory(&category); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create category", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Category created successfully", nil, category)
}



// func (h *Handler) AddServiceToCategory(c *fiber.Ctx) error {
	
// 	var request model.AddServiceToCategoryRequest
// 	// Parse the JSON body into the request struct
// 	if err := c.BodyParser(&request); err != nil {
// 		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
// 	}

// 	// Extract categoryID and serviceID from the request struct
// 	categoryID := request.CategoryID
// 	serviceID := request.ServiceID

// 	// fmt.Println(categoryID, serviceID) // Debugging line to print categoryID and serviceID

// 	if err := h.usecase.AddServiceToCategory(&categoryID, &serviceID); err != nil {
// 		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to add service to category", err, nil)
// 	}
// 	return utils.ResponseJSON(c, fiber.StatusCreated, "Service added to category successfully", nil, request)
// }

// func (h *Handler) RemoveServiceFromCategory(c *fiber.Ctx) error {
// 	var request model.RemoveServiceFromCategoryRequest
// 	if err := c.BodyParser(&request); err != nil {
// 		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
// 	}
// 	categoryID := request.CategoryID
// 	serviceID := request.ServiceID
// 	fmt.Println(categoryID, serviceID)
	
// 	if err := h.usecase.RemoveServiceFromCategory(&categoryID, &serviceID); err != nil {
// 		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to remove service from category", err, nil)
// 	}
// 	return utils.ResponseJSON(c, fiber.StatusOK, "Service removed from category successfully", nil, request)
// }


func (h *Handler) GetAllCategory(c *fiber.Ctx) error {
	categories, err := h.usecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve categories", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Categories retrieved successfully", nil, categories)
}

func (h *Handler) GetByCategoryID(c *fiber.Ctx) error {
	id:= c.Params("id")
	category, err := h.usecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Category retrieved successfully", nil, category)
}

// func (h *Handler) GetServicesByCategoryID(c *fiber.Ctx) error {
// 	id:= c.Params("id")

	
// 	services, err := h.usecase.GetServicesByCategoryID(&id)
// 	if err != nil {
// 		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve services", err, nil)
// 	}
// 	return utils.ResponseJSON(c, fiber.StatusOK, "Services retrieved successfully", nil, services)
// }

// func (h *Handler) GetCategoriesByServiceID(c *fiber.Ctx) error {
// 	id:= c.Params("id")
// 	categories, err := h.usecase.GetCategoriesByServiceID(&id)
// 	if err != nil {
// 		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve categories", err, nil)
// 	}
// 	return utils.ResponseJSON(c, fiber.StatusOK, "Categories retrieved successfully", nil, categories)
// }


func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	id:= c.Params("id")
	var category Entities.Category
	if err := c.BodyParser(&category); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
	}
	
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format", err, nil)
	}

	

	category.ID = uuidID

	if err := h.usecase.Update(&category); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to update category", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Category updated successfully", nil, category)
}

func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	// id, err := uuid.Parse(c.Params("id"))
	id := c.Params("id")
	err := h.usecase.Delete(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Category deleted successfully", nil, nil)
}
