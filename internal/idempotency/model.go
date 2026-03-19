package idempotency

import "time"

type IdempotencyKey struct {
	ID          uint   `gorm:"primaryKey"`
	Key         string `gorm:"uniqueIndex:idx_key_user"`
	UserID      uint   `gorm:"uniqueIndex:idx_key_user"`
	RequestHash string
	Response    []byte
	StatusCode  int
	CreatedAt   time.Time
}
