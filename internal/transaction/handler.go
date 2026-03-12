package transaction

import (
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

	userIDValue := c.Locals("user_id")

	var userID uint

	switch v := userIDValue.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	default:
		return c.Status(500).JSON(fiber.Map{
			"error": "invalid user id type",
		})
	}

	err := h.service.Deposit(userID, body.Amount)
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

	userID := c.Locals("user_id").(uint)
	err := h.service.Withdraw(userID, body.Amount)
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

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request"})
	}

	userID := c.Locals("user_id").(uint)

	err := h.service.Transfer(userID, body.ToUserID, body.Amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "transfer successful",
	})
}
