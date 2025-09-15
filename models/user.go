package models

import "time"

type User struct {
	ID          int64     `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type RequestOTP struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyOTP struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}
