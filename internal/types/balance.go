package types

import "time"

type Balance struct {
	Id        int       `json:"-"`
	UserId    int       `json:"-"`
	Sum       float64   `json:"current"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
