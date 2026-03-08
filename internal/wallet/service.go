package wallet

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateWallet(UserId uint) error {
	wallet := &Wallet{
		UserID:  UserId,
		Balance: 0,
	}
	return s.repo.Create(wallet)
}

// get wallet balance
func (s *Service) GetBalance(userID uint) (int64, error) {

	wallet, err := s.repo.GetByUserID(userID)
	if err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}
