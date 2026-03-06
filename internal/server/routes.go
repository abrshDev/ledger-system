package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	return app
}
