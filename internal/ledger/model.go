package ledger

import "time"

type LedgerEntry struct {
	ID            uint `gorm:"primaryKey"`
	WalletID      uint
	TransactionID uint
	Type          string `gorm:"type:varchar(10)"`
	Amount        int64
	CreatedAt     time.Time
}
