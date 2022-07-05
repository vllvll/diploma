package types

import "time"

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type Withdraw struct {
	ID        int       `json:"-"`
	UserID    int       `json:"-"`
	Number    string    `json:"order"`
	Sum       float32   `json:"sum"`
	CreatedAt time.Time `json:"processed_at"`
}
