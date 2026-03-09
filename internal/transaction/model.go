package transaction

import "time"

type TransactionType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
	Transfer TransactionType = "transfer"
)

type Transaction struct {
	ID         uint `gorm:"primaryKey"`
	FromUserID *uint
	ToUserID   *uint
	Amount     int64
	Type       TransactionType
	CreatedAt  time.Time
}
