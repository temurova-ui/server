package models

import "errors"

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