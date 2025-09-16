// models/user.go

package models

import "time"

type User struct {
	ID          int64     `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type RequestOTPPayload struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyOTPPayload struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

type LoginSuccessResponse struct {
	Token string `json:"token"`
}

type UserListResponse struct {
	Users []*User `json:"users"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
}
