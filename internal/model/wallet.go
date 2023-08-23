package model

import "time"

type Wallet struct {
	ID        int64     `json:"id"`
	Address   string    `json:"address"`
	Balance   float64   `json:"balance"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CheckBalanceRequest struct {
	UserID  int64  `json:"user_id"`
	Address string `json:"address"`
}

type CheckBalanceResponse struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

type TopUpRequest struct {
	Nominal float64 `json:"nominal"`
	Address string  `json:"address"`
	UserID  int64   `json:"user_id"`
}

type PayRequest struct {
	NominalPayment float64 `json:"nominal_payment"`
	Address        string  `json:"address"`
	UserID         int64   `json:"user_id"`
}
