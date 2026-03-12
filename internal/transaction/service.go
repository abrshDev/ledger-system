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
func (s *Service) Transfer(fromUserID uint, toUserID uint, amount int64) error {

	if fromUserID == toUserID {
		return fmt.Errorf("cannot transfer to yourself")
	}

	// sender wallet
	senderWallet, err := s.walletSvc.GetWalletByUserID(fromUserID)
	if err != nil {
		return err
	}

	// receiver wallet
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

	err = s.walletSvc.Update(senderWallet)
	if err != nil {
		return err
	}

	err = s.walletSvc.Update(receiverWallet)
	if err != nil {
		return err
	}

	// transaction record
	txn := &Transaction{
		FromUserID: &fromUserID,
		ToUserID:   &toUserID,
		Amount:     amount,
		Type:       Transfer,
	}

	return s.repo.Create(txn)
}
