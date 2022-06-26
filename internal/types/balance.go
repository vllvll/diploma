package types

import "time"

type Balance struct {
	Id        int       `json:"-"`
	UserId    int       `json:"-"`
	Sum       int       `json:"current"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
