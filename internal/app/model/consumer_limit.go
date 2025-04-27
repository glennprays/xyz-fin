package model

import "time"

type ConsumerLimit struct {
	ConsumerNIK string    `json:"consumer_nik"`
	Tenor       int       `json:"tenor"`
	LimitAmount float64   `json:"limit_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ConsumerLimitResponse struct {
	LimitAmount float64 `json:"limit_amount"`
	Tenor       int     `json:"tenor"`
}
