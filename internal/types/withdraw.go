package types

import "time"

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type Withdraw struct {
	Id        int       `json:"-"`
	UserId    int       `json:"-"`
	Number    int64     `json:"order,string"`
	Sum       float64   `json:"sum"`
	CreatedAt time.Time `json:"processed_at"`
}
