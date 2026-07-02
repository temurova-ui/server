package models

import (
	"errors"
)

type RegisterRequest struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest)Validate() error{
	if r.Name == "" || r.Email == "" || len(r.Password) < 8{
		return errors.New("validation error: name and email are required, password must be >= 8 chars")
	}
	return nil
}

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest)Validate()error{
	if r.Email == "" || len(r.Password) < 8{
		return errors.New("validation error: email and password are required")
	}
	return nil
}

type AuthResponse struct{
	Token string `json:"token"`
}

type UpdateProfileRequest struct{
	Name string `json:"name"`
	Email string `json:"email"`
}

func (r *UpdateProfileRequest)Validate()error{
	if r.Name == "" || r.Email == ""{
		return errors.New("validation error: name and email cannot be empty")
	}
	return nil
}

type ChangePassword struct{
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (c *ChangePassword)Validate()error{
	if len(c.OldPassword) == 0 || len(c.NewPassword) == 0{
		return errors.New("validation error: name and password cannot be empty")
	}
	return nil
}

type CreateOrderRequest struct {
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.ItemName == "" {
		return errors.New("validation error: item_name is required")
	}
	if r.Quantity <= 0 {
		return errors.New("validation error: quantity must be greater than 0")
	}
	return nil
}

type CancelOrder struct{
	Status string `json:"status"`
	UserID int `json:"userID"`
}

func (r *CancelOrder) Validate()error{
	if len(r.Status) == 0{
		return errors.New("validation error: status is empty")
	}
	
	return nil
}


type VerifyRequest struct{
	Email string `json:"email"`
	Otp string `json:"otp"`
}

func (r *VerifyRequest)Validate()error{
	if r.Email == "" || r.Otp == ""{
		return errors.New("validation error: emailand otp should not be empty")
	}
	return nil
}

type RefreshRequest struct{
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshRequest)Validate()error{
	if r.RefreshToken == "" {
		return errors.New("validation error: refresh_token is required")
	}
	return nil
}

type LogoutRequest struct{
	RefreshToken string `json:"refresh_token`
}

func (r *LogoutRequest)Validate()error{
	if r.RefreshToken == ""{
		return errors.New("validation error:refresh_token is required")
	}
	return nil
}