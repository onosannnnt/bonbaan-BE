package recommendationAdepter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	recommendationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/recommendation"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

// RecommendationHandler handles recommendation-related HTTP requests.
type RecommendationHandler struct {
    RecommendationUsecase recommendationUsecase.RecommendationRepository
}

// NewRecommendationHandler creates a new RecommendationHandler.
func NewRecommendationHandler(usecase recommendationUsecase.RecommendationRepository) *RecommendationHandler {
    return &RecommendationHandler{
        RecommendationUsecase: usecase,
    }
}

// GetSuggestions retrieves service suggestions based on the user's latest vow record.
func (h *RecommendationHandler) SuggestNextServie(c *fiber.Ctx) error {
    // Parse pagination parameters.
    var pagination model.Pagination
    if err := c.QueryParser(&pagination); err != nil {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid pagination parameters", err, nil)
    }
    // Set default pagination if needed.
    if pagination.PageSize <= 0 {
        pagination.PageSize = 10
    }
    if pagination.CurrentPage <= 0 {
        pagination.CurrentPage = 1
    }

    // Get userID from the context (like in ReviewHandler).
    userID, ok := c.Locals(constance.UserID_ctx).(string)
    if !ok || userID == "" {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Unable to retrieve userID from token", nil, nil)
    }

    // Call the SuggestNextServie usecase method.
    outputs, pag, err := h.RecommendationUsecase.SuggestNextServie(userID, &pagination)
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get service suggestions", err, nil)
    }

    return utils.ResponseJSON(c, fiber.StatusOK, "Service suggestions retrieved successfully", nil, map[string]interface{}{
        "services":   outputs,
        "pagination": pag,
    })
}

// InterestRating retrieves services based on interest ratings.
func (h *RecommendationHandler) InterestRating(c *fiber.Ctx) error {
    // Parse pagination parameters.
    var pagination model.Pagination
    if err := c.QueryParser(&pagination); err != nil {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid pagination parameters", err, nil)
    }
    // Set default pagination if needed.
    if pagination.PageSize <= 0 {
        pagination.PageSize = 10
    }
    if pagination.CurrentPage <= 0 {
        pagination.CurrentPage = 1
    }

    // Get userID from the context (like in ReviewHandler).
    userID, ok := c.Locals(constance.UserID_ctx).(string)
    if !ok || userID == "" {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Unable to retrieve userID from token", nil, nil)
    }

    // Call the SuggestNextServie usecase method.
    outputs, pag, err := h.RecommendationUsecase.InterestRating(userID, &pagination)
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get service suggestions", err, nil)
    }

    return utils.ResponseJSON(c, fiber.StatusOK, "Service suggestions retrieved successfully", nil, map[string]interface{}{
        "services":   outputs,
        "pagination": pag,
    })
}	