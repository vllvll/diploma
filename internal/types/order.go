package types

import "time"

type Order struct {
	Id         int       `json:"-"`
	Number     int64     `json:"number,string"`
	UserId     int       `json:"-"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}
