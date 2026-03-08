package server

import (
	"github.com/abrshDev/ledger-system/internal/middleware"
	"github.com/abrshDev/ledger-system/internal/wallet"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	walletRepo := wallet.NewRepository(db)
	walletService := wallet.NewService(walletRepo)
	walletHandler := wallet.NewHandler(walletService)

	app.Get("/wallet/balance", middleware.Protected(), walletHandler.GetBalance)
}

func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	SetupRoutes(app, db)

	return app
}
