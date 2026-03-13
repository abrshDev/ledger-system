package ledger

import "time"

type LedgerEntry struct {
	ID        uint `gorm:"primaryKey"`
	WalletID  uint
	Type      string
	Amount    int64
	Reference string
	CreatedAt time.Time
}
