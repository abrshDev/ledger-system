package transaction

import (
	"github.com/abrshDev/ledger-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type DepositRequest struct {
	Amount int64 `json:"amount"`
}

type WithdrawRequest struct {
	Amount int64 `json:"amount"`
}

type TransferRequest struct {
	ToUserID uint  `json:"to_user_id"`
	Amount   int64 `json:"amount"`
}

func (h *Handler) Deposit(c *fiber.Ctx) error {

	var body DepositRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request"})
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	err = h.service.Deposit(userID, body.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "deposit successful",
	})
}

func (h *Handler) Withdraw(c *fiber.Ctx) error {

	var body WithdrawRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request"})
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	err = h.service.Withdraw(userID, body.Amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "withdraw successful",
	})
}
func (h *Handler) Transfer(c *fiber.Ctx) error {

	var body TransferRequest

	// Parse request body first
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	// Validate amount after parsing
	if body.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "amount must be greater than zero",
		})
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	err = h.service.Transfer(userID, body.ToUserID, body.Amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "transfer successful",
	})
}

func (h *Handler) GetTransactions(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	txns, err := h.service.GetUserTransactions(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch transactions",
		})
	}

	return c.JSON(fiber.Map{
		"transactions": txns,
	})
}
