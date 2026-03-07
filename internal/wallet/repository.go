package wallet

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(wallet *Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *Repository) GetByUserID(userID uint) (*Wallet, error) {
	var wallet Wallet

	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
