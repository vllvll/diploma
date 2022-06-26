package types

import "time"

type Order struct {
	Id         int       `json:"-"`
	Number     int       `json:"number"`
	UserId     int       `json:"-"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}
