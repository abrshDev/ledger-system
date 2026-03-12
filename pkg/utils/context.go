package utils

import "github.com/gofiber/fiber/v2"

func GetUserID(c *fiber.Ctx) (uint, error) {

	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case float64:
		return uint(v), nil
	case uint:
		return v, nil
	case int:
		return uint(v), nil
	default:
		return 0, fiber.NewError(fiber.StatusUnauthorized, "invalid user id")
	}
}
