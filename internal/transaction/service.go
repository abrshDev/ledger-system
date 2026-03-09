package transaction

import "github.com/abrshDev/ledger-system/internal/wallet"

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
