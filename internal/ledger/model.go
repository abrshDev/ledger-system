package ledger

import "time"

type LedgerEntry struct {
	ID        uint   `gorm:"primaryKey"`
	WalletID  uint   `gorm:"not null"`
	PaymentID uint   `gorm:"not null"`
	Type      string `gorm:"not null"` // debit | credit
	Amount    int64  `gorm:"not null"`
	CreatedAt time.Time
}
