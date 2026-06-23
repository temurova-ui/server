package service

import(
	"context"
	"net/mail"
	"myservice/internal/models"
	"myservice/pkg/errs"
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

func (s *UserService)Register(ctx context.Context, name, email, password string)(*models.User, error){
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" || len(password) < 6{
		return nil, errs.NewInvalidInputError("name cannot be empty and password must be at least 6 chars")
	}
	if _, err := mail.ParseAddress(email); err != nil{
		return nil, errs.ErrInvalidInput 
	}

	hashedPassword := "hashed_" + password

	user := &models.User{
		Name: name,
		Email: email,
		PasswordHash: hashedPassword, 
		Role: "user", 
	}

	err := s.repo.Create(ctx, user)
	if err != nil{
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint"){
			return nil, errs.ErrEmailConflict
		}
		return nil, err
	}
	return user, nil
}
