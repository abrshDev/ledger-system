package transaction

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(txn *Transaction) error {
	return r.db.Create(txn).Error
}
func (r *Repository) GetByUserID(userID uint) ([]Transaction, error) {
	var txns []Transaction

	err := r.db.
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&txns).Error

	return txns, err
}
