package transaction

import (
	"fmt"

	"github.com/abrshDev/ledger-system/internal/wallet"
)

type Service struct {
	repo      *Repository
	walletSvc *wallet.Service
}

func NewService(repo *Repository, walletSvc *wallet.Service) *Service {
	return &Service{
		repo:      repo,
		walletSvc: walletSvc,
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
