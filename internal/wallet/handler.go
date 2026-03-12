package wallet

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetBalance(c *fiber.Ctx) error {

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
