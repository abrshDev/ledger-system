package server

import (
	"github.com/abrshDev/ledger-system/internal/auth"
	"github.com/abrshDev/ledger-system/internal/middleware"
	"github.com/abrshDev/ledger-system/internal/transaction"
	"github.com/abrshDev/ledger-system/internal/user"
	"github.com/abrshDev/ledger-system/internal/wallet"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	// repositories
	userRepo := user.NewRepository(db)
	walletRepo := wallet.NewRepository(db)
	txnRepo := transaction.NewRepository(db)

	// services
	walletService := wallet.NewService(walletRepo)
	txnService := transaction.NewService(txnRepo, walletService)

	// handlers
	txnHandler := transaction.NewHandler(txnService)
	authHandler := auth.NewAuthHandler(userRepo, walletService)
	walletHandler := wallet.NewHandler(walletService)

	// auth routes
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)
	app.Post("/auth/refresh", authHandler.Refresh)

	// wallet routes
	app.Get("/wallet/balance", middleware.Protected(), walletHandler.GetBalance)

	// transaction routes
	app.Post("/wallet/deposit", middleware.Protected(), txnHandler.Deposit)
	app.Post("/wallet/withdraw", middleware.Protected(), txnHandler.Withdraw)
	app.Post("/wallet/transfer", middleware.Protected(), txnHandler.Transfer)
}
func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	SetupRoutes(app, db)

	return app
}
