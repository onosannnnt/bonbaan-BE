package reviewAdepter

import (
	"github.com/gofiber/fiber/v2"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	reviewUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/review"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type ReviewHandler struct {
	ReviewUsecase reviewUsecase.ReviewUsecase
}

func NewReviewHandler(usecase reviewUsecase.ReviewUsecase) *ReviewHandler {

	return &ReviewHandler{ReviewUsecase: usecase}

}

func (h *ReviewHandler) Insert(c *fiber.Ctx) error {
	var review Entities.Review

	if err := c.BodyParser(&review); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the require fields", err, nil)
	}

	if err := h.ReviewUsecase.Insert(&review); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "this role already exists", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "success", nil, review)

}