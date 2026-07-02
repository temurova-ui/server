package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	// "net/mail"
	"myservice/internal/models"
	"myservice/internal/repository"
	utils "myservice/internal/utils"
	"myservice/pkg/cache"
	"myservice/pkg/errs"
	"myservice/pkg/jwt"
	"myservice/pkg/password"
	smtp2 "myservice/pkg/smtp"
)

type UserService interface{
	Register(ctx context.Context, request models.RegisterRequest) error
	Login(ctx context.Context, req models.LoginRequest)(response models.LoginResponse, err error)
	GetProfile(ctx context.Context, userID int)(*models.User, error)
	ChangePassword(ctx context.Context, id int, req models.ChangePassword)error
	UpdateProfile(ctx context.Context, userID int, req models.UpdateProfileRequest)error
	DeleteAccount(ctx context.Context, userID int)error
	Verify (ctx context.Context, request models.VerifyRequest)error
	GetLoginHistory(ctx context.Context, userID int64)([]models.LoginHistory, error)
	Refresh(ctx context.Context, request models.RefreshRequest)(response models.LoginResponse, err error)
	Logout(ctx context.Context, request models.LogoutRequest) error

}

type userService struct{
	repo repository.UserRepo
	loginHistoryRepository *repository.LoginHistoryRepository
	cache cache.MemoryCache
	smtp  *smtp2.SMTP
}

func NewUserService(repo repository.UserRepo, cache cache.MemoryCache, smtp *smtp2.SMTP ) UserService{
	return &userService{
		repo: repo, cache: cache, smtp: smtp,
	}
}


func (s *userService) Register(ctx context.Context, request models.RegisterRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	exists, err := s.repo.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("user with this email already exists")
	}

	_, ok := s.cache.Get(request.Email)
	if ok {
		return errors.New("user with this email already exists")
	}

	passwordHash, err := password.Hash(request.Password)
	if err != nil {
		return err
	}

	otp := utils.GenerateOTP()
	err = s.smtp.SendOTP(ctx, request.Email, otp)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		PasswordHash: passwordHash,
		Role:     models.UserRole,
		OtpCode:  otp,
	}

	s.cache.Set(request.Email, user, time.Minute*5)
	return nil
}

func (s *userService)Login(ctx context.Context, req models.LoginRequest) (response models.LoginResponse, err error){
	err = req.Validate()
	if err != nil {
		return models.LoginResponse{}, err
	}

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return models.LoginResponse{}, err
	}
	fmt.Println(user.PasswordHash)
	fmt.Println(req.Password)


	err = password.Compare(user.PasswordHash, req.Password)
	if err != nil {
		return models.LoginResponse{}, err
	}

	accessToken, err := jwt.GenereteToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.LoginResponse{}, err
	}

	refreshToken := utils.GenerateRand()
	refreshTokenHash := utils.HashRefreshToken(refreshToken)
	fmt.Println(refreshToken)
	refreshTokenEntity := models.RefreshToken{
		UserID: int64(user.ID),
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	err = s.repo.AddRefreshToken(ctx, refreshTokenEntity)
	if err != nil {
		return models.LoginResponse{}, err
	}
	fmt.Println(accessToken)
	fmt.Println(refreshToken)
	response = models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
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
	err :=  s.repo.Delete(ctx, userID)
	if err != nil{
		return err
	}
	return nil
}

func (s *userService)Verify (ctx context.Context, request models.VerifyRequest)error{
		err := request.Validate()
		if err != nil{
			return err
		}

		cacheInfo, ok := s.cache.Get(request.Email)
		if !ok{
			return errors.New("user not founder")
		}

		user, ok := cacheInfo.(models.User)
		if !ok {
			return errors.New("user not found")
		}

		if user.AttemptedOTP >= 3 {
			return errors.New("user is too many attempts")
		}

		if request.Otp != user.OtpCode{
			user.AttemptedOTP++
			s.cache.Set(request.Email, user, time.Minute*5)
			return errors.New("invalid otp")
		}

		err = s.repo.Create(ctx, &user)
		if err != nil{
			return err
		}
		return nil
}

func (s *userService)GetLoginHistory(ctx context.Context, userID int64)([]models.LoginHistory, error) {
	if userID <= 0{
		return nil, errors.New("invalid user id")
	}
	return s.loginHistoryRepository.GetLoginHistory(ctx, userID)
}


func (s *userService) Refresh(ctx context.Context, request models.RefreshRequest) (response models.LoginResponse, err error) {
	if err = request.Validate(); err != nil {
		return models.LoginResponse{}, err
	}

	hash := utils.HashRefreshToken(request.RefreshToken)

	tokenEntity, err := s.repo.GetRefreshToken(ctx, hash)
	if err != nil {
		return models.LoginResponse{}, errors.New("unauthorized: invalid refresh token")
	}

	if time.Now().After(tokenEntity.ExpiresAt) {
		return models.LoginResponse{}, errors.New("unauthorized: expired refresh token")
	}

	user, err := s.repo.GetByID(ctx, int(tokenEntity.UserID))
	if err != nil {
		return models.LoginResponse{}, errors.New("unauthorized: user not found")
	}

	err = s.repo.DeleteRefreshToken(ctx, hash)
	if err != nil {
		return models.LoginResponse{}, err
	}

	accessToken, err := jwt.GenereteToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.LoginResponse{}, err
	}

	newRefreshToken := utils.GenerateRand()
	newRefreshTokenHash := utils.HashRefreshToken(newRefreshToken)

	newRefreshTokenEntity := models.RefreshToken{
		UserID:    int64(user.ID),
		TokenHash: newRefreshTokenHash,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	err = s.repo.AddRefreshToken(ctx, newRefreshTokenEntity)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *userService) Logout(ctx context.Context, request models.LogoutRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	hash := utils.HashRefreshToken(request.RefreshToken)

	err := s.repo.DeleteRefreshToken(ctx, hash)
	if err != nil {
		return err
	}
	return nil
}
