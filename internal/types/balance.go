package types

import "time"

type Balance struct {
	ID        int
	UserID    int
	Sum       float32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ResponseBalance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
