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

func (s *Service) Deposit(userID uint, amount int64) error {
	// 1. Validate input
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	// 2. Start DB transaction
	return s.db.Transaction(func(tx *gorm.DB) error {

		// 3. Get wallet INSIDE transaction
		var wallet wallet.Wallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		// 4. Update balance
		wallet.Balance += amount

		// 5. Save updated wallet
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		// 6. Create transaction record
		txn := &Transaction{
			ToUserID: &userID,
			Amount:   amount,
			Type:     Deposit,
		}

		if err := tx.Create(txn).Error; err != nil {
			return err
		}

		// 7. Create ledger entry (credit)
		if err := tx.Create(&ledger.LedgerEntry{
			WalletID:      wallet.ID,
			TransactionID: txn.ID,
			Type:          "credit",
			Amount:        amount,
		}).Error; err != nil {
			return err
		}

		// 8. Commit (by returning nil)
		return nil
	})
}
func (s *Service) Withdraw(userID uint, amount int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		return nil
	})
	wallet, err := s.walletSvc.GetWalletByUserID(userID)
	if err != nil {
		return err
	}
	if wallet.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// subtract balance
	wallet.Balance -= amount

	// update wallet
	err = s.walletSvc.Update(wallet)
	if err != nil {
		return err
	}
	// create transaction record
	txn := &Transaction{
		FromUserID: &userID,
		Amount:     amount,
		Type:       Withdraw,
	}

	return s.repo.Create(txn)

}
func (s *Service) Transfer(fromUserID uint, toUserID uint, amount int64) error {

	return s.db.Transaction(func(tx *gorm.DB) error {

		if fromUserID == toUserID {
			return fmt.Errorf("cannot transfer to yourself")
		}

		// get wallets
		senderWallet, err := s.walletSvc.GetWalletByUserID(fromUserID)
		if err != nil {
			return err
		}

		receiverWallet, err := s.walletSvc.GetWalletByUserID(toUserID)
		if err != nil {
			return err
		}

		if senderWallet.Balance < amount {
			return fmt.Errorf("insufficient funds")
		}

		// update balances
		senderWallet.Balance -= amount
		receiverWallet.Balance += amount

		if err := tx.Save(senderWallet).Error; err != nil {
			return err
		}

		if err := tx.Save(receiverWallet).Error; err != nil {
			return err
		}

		// create transaction
		txn := &Transaction{
			FromUserID: &fromUserID,
			ToUserID:   &toUserID,
			Amount:     amount,
			Type:       Transfer,
		}

		if err := tx.Create(txn).Error; err != nil {
			return err
		}

		// debit entry
		if err := tx.Create(&ledger.LedgerEntry{
			WalletID:      senderWallet.ID,
			TransactionID: txn.ID,
			Type:          "debit",
			Amount:        amount,
		}).Error; err != nil {
			return err
		}

		// credit entry
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
func (s *Service) GetUserTransactions(userId uint) ([]Transaction, error) {
	return s.repo.GetByUserID(userId)
}
