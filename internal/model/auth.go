package model

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Data  UserPresenter `json:"data"`
	Token string        `json:"token"`
}

type RegisterRequest struct {
	FirstName     string `json:"first_name" validate:"required,min=2,max=20"`
	LastName      string `json:"last_name" validate:"required,min=2,max=20"`
	BirthDate     string `json:"birth_date" validate:"required"`
	StreetAddress string `json:"street_address" validate:"required,min=5,max=40,alphanumspace"`
	City          string `json:"city" validate:"required,min=2,max=20,alpha"`
	Province      string `json:"province" validate:"required"`
	Phone         string `json:"phone" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Username      string `json:"username" validate:"required"`
	Password      string `json:"password" validate:"required"`
}

type UserPresenter struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BirthDate     string `json:"birth_date"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Username      string `json:"username"`
}
