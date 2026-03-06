package payment

import "time"

type Payment struct {
	ID           uint   `gorm:"primaryKey"`
	FromWalletID uint   `gorm:"not null"`
	ToWalletID   uint   `gorm:"not null"`
	Amount       int64  `gorm:"not null"`
	Status       string `gorm:"default:pending"`
	CreatedAt    time.Time
}
