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
	wallet, err := s.walletSvc.GetWalletByUserID(userID)
	if err != nil {
		return err
	}
	wallet.Balance += amount
	err = s.walletSvc.Update(wallet)
	if err != nil {
		return err
	}
	txn := &Transaction{
		ToUserID: &userID,
		Amount:   amount,
		Type:     Deposit,
	}
	return s.repo.Create(txn)
}
func (s *Service) Withdraw(userID uint, amount int64) error {
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
