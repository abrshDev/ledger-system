package idempotency

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) GetByKey(tx *gorm.DB, key string, userID uint) (*IdempotencyKey, error) {
	var record IdempotencyKey
	err := tx.Where("key = ? AND user_id = ?", key, userID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil

}
func (r *Repository) Create(tx *gorm.DB, record *IdempotencyKey) error {
	return tx.Create(record).Error
}
