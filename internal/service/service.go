package service

import (
	"context"
	"errors"

	// "net/mail"
	"myservice/internal/models"
	"myservice/pkg/errs"
	"myservice/pkg/jwt"
	"myservice/pkg/password"
	"strings"
)

type UserRepo interface{
	Create(ctx context.Context, u *models.User)error
	GetByEmail (ctx context.Context, email string)(*models.User, error)
}

type UserService struct{
	repo UserRepo
}

func NewUserService(repo UserRepo)*UserService{
	return &UserService{repo: repo}
}

func (s *UserService)Register(ctx context.Context, req models.RegisterRequest)(string, error){
	hashedPassword, err := password.Hash(req.Password)
	if err != nil{
		return "", err
	}

	user := &models.User{
		Name: strings.TrimSpace(req.Name),
		Email: strings.TrimSpace(strings.ToLower(req.Email)),
		PasswordHash: hashedPassword, 
		Role: "user", 
	}

	err = s.repo.Create(ctx, user)
	if err != nil{
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint"){
			return "", errs.ErrEmailConflict
		}
		return "", err
	}

	return jwt.GenereteToken(user.ID, user.Email, user.Role)
}

func (s *UserService)Login(ctx context.Context, req models.LoginRequest)(string, error){
	email := strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil{
		return "", errors.New("invalid email or password")
	}

	if err := password.Compare(user.PasswordHash, req.Password); err != nil{
		return  "", errors.New("invalid email or password")
	}

	return jwt.GenereteToken(user.ID, user.Email, user.Role)
}