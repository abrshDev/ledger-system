package server

import (
	"github.com/abrshDev/ledger-system/internal/auth"
	"github.com/abrshDev/ledger-system/internal/middleware"
	"github.com/abrshDev/ledger-system/internal/user"
	"github.com/abrshDev/ledger-system/internal/wallet"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	userRepo := user.NewRepository(db)
	walletRepo := wallet.NewRepository(db)

	walletService := wallet.NewService(walletRepo)

	authHandler := auth.NewAuthHandler(userRepo, walletService)
	walletHandler := wallet.NewHandler(walletService)

	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)
	app.Post("/auth/refresh", authHandler.Refresh)

	app.Get("/wallet/balance", middleware.Protected(), walletHandler.GetBalance)
}

func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	SetupRoutes(app, db)

	return app
}
