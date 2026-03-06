package wallet

import "time"

type Wallet struct {
	ID        uint  `gorm:"primaryKey"`
	UserID    uint  `gorm:"not null;unique"`
	Balance   int64 `gorm:"not null;default:0"`
	CreatedAt time.Time
}
