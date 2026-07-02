package models

import "time"


const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	Role string `json:"role"`
	OtpCode string `json:"otp_code"`
	AttemptedOTP int `json:"attempted_otp"`
	CreatedAt time.Time `json:"created_at"`
}

type UserAndOrder struct{
	User User `json:"user"`
	Orders []Order `json:"order"`
}