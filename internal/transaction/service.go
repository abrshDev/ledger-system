package transaction

import (
	"fmt"

	"github.com/abrshDev/ledger-system/internal/ledger"
	"github.com/abrshDev/ledger-system/internal/wallet"
	"gorm.io/gorm"
)

type Service struct {
	repo       *Repository
	walletSvc  *wallet.Service
	ledgerRepo *ledger.Repository
	db         *gorm.DB
}

func NewService(repo *Repository, walletSvc *wallet.Service, ledgerRepo *ledger.Repository, db *gorm.DB) *Service {
	return &Service{
		repo:       repo,
		walletSvc:  walletSvc,
		ledgerRepo: ledgerRepo,
		db:         db,
	}
}

// Deposit adds funds to a user's wallet
func (s *Service) Deposit(userID uint, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get wallet
		var w wallet.Wallet
		if err := tx.Where("user_id = ?", userID).First(&w).Error; err != nil {
			return err
		}

		// Update balance
		w.Balance += amount
		if err := tx.Save(&w).Error; err != nil {
			return err
		}

		// Create transaction
		txn := &Transaction{
			ToUserID: &userID,
			Amount:   amount,
			Type:     Deposit,
		}
		if err := tx.Create(txn).Error; err != nil {
			return err
		}

		// Ledger entry (credit)
		if err := tx.Create(&ledger.LedgerEntry{
			WalletID:      w.ID,
			TransactionID: txn.ID,
			Type:          "credit",
			Amount:        amount,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// Withdraw subtracts funds from a user's wallet
func (s *Service) Withdraw(userID uint, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get wallet
		var w wallet.Wallet
		if err := tx.Where("user_id = ?", userID).First(&w).Error; err != nil {
			return err
		}

		if w.Balance < amount {
			return fmt.Errorf("insufficient funds")
		}

		// Update balance
		w.Balance -= amount
		if err := tx.Save(&w).Error; err != nil {
			return err
		}

		// Create transaction
		txn := &Transaction{
			FromUserID: &userID,
			Amount:     amount,
			Type:       Withdraw,
		}
		if err := tx.Create(txn).Error; err != nil {
			return err
		}

		// Ledger entry (debit)
		if err := tx.Create(&ledger.LedgerEntry{
			WalletID:      w.ID,
			TransactionID: txn.ID,
			Type:          "debit",
			Amount:        amount,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// Transfer moves funds from one user to another
func (s *Service) Transfer(fromUserID, toUserID uint, amount int64) error {
	if fromUserID == toUserID {
		return fmt.Errorf("cannot transfer to yourself")
	}

	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get wallets
		var senderWallet, receiverWallet wallet.Wallet
		if err := tx.Where("user_id = ?", fromUserID).First(&senderWallet).Error; err != nil {
			fmt.Println("error senderwallet")
			return err

		}
		if err := tx.Where("user_id = ?", toUserID).First(&receiverWallet).Error; err != nil {
			fmt.Println("error in receiverwallet")
			return err
		}

		if senderWallet.Balance < amount {
			return fmt.Errorf("insufficient funds")
		}

		// Update balances
		senderWallet.Balance -= amount
		receiverWallet.Balance += amount

		if err := tx.Save(&senderWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&receiverWallet).Error; err != nil {
			return err
		}

		// Create transaction
		txn := &Transaction{
			FromUserID: &fromUserID,
			ToUserID:   &toUserID,
			Amount:     amount,
			Type:       Transfer,
		}
		if err := tx.Create(txn).Error; err != nil {
			return err
		}

		// Ledger entries
		if err := tx.Create(&ledger.LedgerEntry{
			WalletID:      senderWallet.ID,
			TransactionID: txn.ID,
			Type:          "debit",
			Amount:        amount,
		}).Error; err != nil {
			return err
		}
		if err := tx.Create(&ledger.LedgerEntry{
			WalletID:      receiverWallet.ID,
			TransactionID: txn.ID,
			Type:          "credit",
			Amount:        amount,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetUserTransactions returns all transactions for a given user
func (s *Service) GetUserTransactions(userID uint) ([]Transaction, error) {
	return s.repo.GetByUserID(userID)
}
