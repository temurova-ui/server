package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	// "net/mail"
	"myservice/internal/models"
	"myservice/internal/repository"
	"myservice/pkg/errs"
	"myservice/pkg/jwt"
	"myservice/pkg/password"
)

type UserService interface{
	Register(ctx context.Context, req models.RegisterRequest)(string, error)
	Login(ctx context.Context, req models.LoginRequest)(string, error)
	GetProfile(ctx context.Context, userID int)(*models.User, error)
	ChangePassword(ctx context.Context, id int, req models.ChangePassword)error
	UpdateProfile(ctx context.Context, userID int, req models.UpdateProfileRequest)error
	DeleteAccount(ctx context.Context, userID int)error

}

type userService struct{
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService{
	return &userService{
		repo: repo,
	}
}


func (s *userService)Register(ctx context.Context, req models.RegisterRequest)(string, error){
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

func (s *userService)Login(ctx context.Context, req models.LoginRequest)(string, error){
	email := strings.TrimSpace(strings.ToLower(req.Email))
	
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil{
		return "", errors.New("invalid email or password")
	}
	
	if err := password.Compare(user.PasswordHash, req.Password); err != nil{
		return  "", errors.New("invalid email or password")
	}
	fmt.Println(user.ID)
	return jwt.GenereteToken(user.ID, user.Email, user.Role)
}

func (s *userService)GetProfile(ctx context.Context, userID int)(*models.User, error){
	return s.repo.GetByID(ctx, userID)
}

func (s *userService)ChangePassword(ctx context.Context, userID int, req models.ChangePassword)error{

	err := req.Validate()
	if err != nil{
		return errs.ErrValidate
	}
	user, err := s.repo.GetByID(ctx, userID) 
	if err != nil{
		return err
	}

	err = password.Compare(user.PasswordHash, req.OldPassword)
	if err != nil{
		return errors.New("error from password.Compare")
	}

	newHash, err := password.Hash(req.NewPassword)
	if err != nil{
		return errors.New("error from newHash")
	}
	err = s.repo.UpdatePassword(ctx, userID, newHash)
	if err != nil{
		return errors.New("error from s.repo.UpdatePassword")
	}
	return nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID int, req models.UpdateProfileRequest)error{
	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(strings.ToLower(req.Email))
	return s.repo.Update(ctx, userID, name, email)
}

func (s *userService)DeleteAccount(ctx context.Context, userID int)error{
	return s.repo.Delete(ctx, userID)
}


