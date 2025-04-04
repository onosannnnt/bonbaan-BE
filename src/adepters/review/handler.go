package reviewAdepter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
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
	review := model.ReviewInsertRequest{}
	if err := c.BodyParser(&review); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the required fields", err, nil)
	}

	insertReview := Entities.Review{}
	insertReview.Rating = review.Rating
	insertReview.Detail = review.Detail
	userID, ok := c.Locals(constance.UserID_ctx).(string)
	if !ok {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Unable to get UserID from token", nil, nil)
	}
	parseUserID, err := uuid.Parse(userID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format for UserID", err, nil)
	}
	insertReview.UserID = parseUserID
	insertReview.OrderID, err = uuid.Parse(review.OrderID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format for OrderID", err, nil)
	}
	insertReview.ServiceID, err = uuid.Parse(review.ServiceID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format for ServiceID", err, nil)
	}

	if err := h.ReviewUsecase.Insert(&insertReview); err != nil {
		return utils.ResponseJSON(c, fiber.StatusConflict, "This review already exists", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Review created successfully", nil, review)
}

func (h *ReviewHandler) GetAll(c *fiber.Ctx) error {
	reviews, err := h.ReviewUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get reviews", nil, err)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Reviews retrieved successfully", nil, reviews)
}

func (h *ReviewHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	review, err := h.ReviewUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get review", nil, err)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Review retrieved successfully", nil, review)
}

func (h *ReviewHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var review Entities.Review
	if err := c.BodyParser(&review); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format", err, nil)
	}

	review.ID = uuidID

	if err := h.ReviewUsecase.Update(&id, &review); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to update review", err, nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Review updated successfully", nil, review)
}

func (h *ReviewHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.ReviewUsecase.Delete(&id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to delete review", nil, err)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Review deleted successfully", nil, nil)
}
