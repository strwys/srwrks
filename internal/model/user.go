package model

import "time"

type User struct {
	ID            int64     `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	BirthDate     string    `json:"birth_date" validate:"required"`
	StreetAddress string    `json:"street_address" validate:"required,min=5,max=40,alphanum"`
	City          string    `json:"city" validate:"required,min=2,max=20"`
	Province      string    `json:"province" validate:"required"`
	Phone         string    `json:"phone" validate:"required"`
	Email         string    `json:"email" validate:"required,email"`
	Username      string    `json:"username" validate:"required"`
	Password      string    `json:"password" validate:"required"`
	Token         string    `json:"token"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type LoginHistory struct {
	BrowserName  string    `json:"browser_name"`
	LoginSucceed int       `json:"login_succeed"`
	LoginFailed  int       `json:"login_failed"`
	UserID       int64     `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
