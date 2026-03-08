package wallet

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetBalance(c *fiber.Ctx) error {

	uid := c.Locals("user_id").(float64)
	userID := uint(uid)

	balance, err := h.service.GetBalance(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "wallet not found",
		})
	}

	return c.JSON(fiber.Map{
		"balance": balance,
	})
}
