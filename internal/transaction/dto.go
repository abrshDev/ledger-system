package transaction

type TransferRequest struct {
	ToUserID uint  `json:"to_user_id"`
	Amount   int64 `json:"amount"`
}
