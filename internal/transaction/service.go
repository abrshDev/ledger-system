package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/abrshDev/ledger-system/internal/idempotency"
	"github.com/abrshDev/ledger-system/internal/ledger"
	"github.com/abrshDev/ledger-system/internal/wallet"
	"gorm.io/gorm"
)

type Service struct {
	repo       *Repository
	walletSvc  *wallet.Service
	ledgerRepo *ledger.Repository
	db         *gorm.DB
	idemRepo   *idempotency.Repository
}

func NewService(repo *Repository, walletSvc *wallet.Service, ledgerRepo *ledger.Repository, db *gorm.DB, idemRepo *idempotency.Repository) *Service {
	return &Service{
		repo:       repo,
		walletSvc:  walletSvc,
		ledgerRepo: ledgerRepo,
		db:         db,
		idemRepo:   idemRepo,
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

func (s *Service) Transfer(fromUserID, toUserID uint, amount int64, key string) ([]byte, int, error) {

	if key == "" {
		return nil, 400, fmt.Errorf("missing idempotency key")
	}

	if fromUserID == toUserID {
		return nil, 400, fmt.Errorf("cannot transfer to yourself")
	}

	if amount <= 0 {
		return nil, 400, fmt.Errorf("amount must be greater than zero")
	}

	reqHash := hashTransfer(fromUserID, toUserID, amount)

	var response []byte
	var status int

	err := s.db.Transaction(func(tx *gorm.DB) error {

		existing, err := s.idemRepo.GetByKey(tx, key, fromUserID)
		if err == nil {
			if existing.RequestHash != reqHash {
				return fmt.Errorf("idempotency key reused with different payload")
			}

			response = existing.Response
			status = existing.StatusCode
			return nil
		}

		var senderWallet, receiverWallet wallet.Wallet

		if err := tx.Where("user_id = ?", fromUserID).First(&senderWallet).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", toUserID).First(&receiverWallet).Error; err != nil {
			return err
		}

		if senderWallet.Balance < amount {
			return fmt.Errorf("insufficient funds")
		}

		senderWallet.Balance -= amount
		receiverWallet.Balance += amount

		if err := tx.Save(&senderWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&receiverWallet).Error; err != nil {
			return err
		}

		txn := &Transaction{
			FromUserID: &fromUserID,
			ToUserID:   &toUserID,
			Amount:     amount,
			Type:       Transfer,
		}

		if err := tx.Create(txn).Error; err != nil {
			return err
		}

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

		// build response
		result := map[string]interface{}{
			"message": "transfer successful",
			"amount":  amount,
		}

		resBytes, _ := json.Marshal(result)

		// save idempotency record
		err = s.idemRepo.Create(tx, &idempotency.IdempotencyKey{
			Key:         key,
			UserID:      fromUserID,
			RequestHash: reqHash,
			Response:    resBytes,
			StatusCode:  200,
		})
		if err != nil {
			return err
		}

		response = resBytes
		status = 200

		return nil
	})

	return response, status, err
}

// GetUserTransactions returns all transactions for a given user
func (s *Service) GetUserTransactions(userID uint) ([]Transaction, error) {
	return s.repo.GetByUserID(userID)
}
func hashTransfer(fromUserID, toUserID uint, amount int64) string {
	data := fmt.Sprintf("%d:%d:%d", fromUserID, toUserID, amount)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
