package types

import "time"

type Order struct {
	ID         int       `json:"-"`
	Number     string    `json:"number"`
	UserID     int       `json:"-"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type OrderLoyalty struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
