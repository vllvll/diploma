package types

import "time"

type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type Withdraw struct {
	Id        int       `json:"-"`
	UserId    int       `json:"-"`
	Number    string    `json:"order"`
	Sum       int       `json:"sum"`
	CreatedAt time.Time `json:"processed_at"`
}
