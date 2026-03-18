package config

import (
	"fmt"
	"os"

	"github.com/abrshDev/ledger-system/internal/ledger"

	"github.com/abrshDev/ledger-system/internal/transaction"
	"github.com/abrshDev/ledger-system/internal/user"
	"github.com/abrshDev/ledger-system/internal/wallet"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")

	var db *gorm.DB
	var err error

	if dbURL != "" {
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	}

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&user.User{},
		&wallet.Wallet{},

		&ledger.LedgerEntry{},
		&transaction.Transaction{},
	); err != nil {
		return nil, err
	}

	fmt.Println("Database connected successfully")

	return db, nil
}
