package transactionAdepter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	transactionUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/transaction"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type TransactioHandler struct {
	transactionUsecase transactionUsecase.TransectionUsecase
}

func NewStatusHandler(transactionUsecase transactionUsecase.TransectionUsecase) *TransactioHandler {
	return &TransactioHandler{
		transactionUsecase: transactionUsecase,
	}
}

func (h *TransactioHandler) InsertTransaction(c *fiber.Ctx) error {
	var transaction Entities.Transaction

	if err := c.BodyParser(&transaction); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the require fields", err, nil)
	}

	if err := h.transactionUsecase.Insert(&transaction); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to insert transaction", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusCreated, "success", nil, transaction)
}

func (h *TransactioHandler) GetAllTransaction(c *fiber.Ctx) error {
	transactions, err := h.transactionUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to get all transaction", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, transactions)
}

func (h *TransactioHandler) GetTransactionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	transaction, err := h.transactionUsecase.GetByID(id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to get transaction", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, transaction)
}

func (h *TransactioHandler) UpdateTransaction(c *fiber.Ctx) error {
	var transaction Entities.Transaction
	id := c.Params("id")
	transaction.ID = uuid.MustParse(id)
	if err := c.BodyParser(&transaction); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the require fields", err, nil)
	}

	if err := h.transactionUsecase.Update(&transaction); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to update transaction", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, transaction)
}

func (h *TransactioHandler) DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.transactionUsecase.Delete(id); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Failed to delete transaction", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, nil)
}
