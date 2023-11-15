package usecase

import (
	"context"
	"errors"

	"github.com/adityatresnobudi/job-portal/dto"
	"github.com/adityatresnobudi/job-portal/helper"
	"github.com/adityatresnobudi/job-portal/model"
	"github.com/adityatresnobudi/job-portal/repository"
	"github.com/adityatresnobudi/job-portal/shared"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo repository.UserRepository
}

type UserUsecase interface {
	CreateUsers(ctx context.Context, user dto.UserPayload) (dto.UserResponse, error)
	LoginUser(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uu *userUsecase) CreateUsers(ctx context.Context, user dto.UserPayload) (dto.UserResponse, error) {
	newUser := model.Users{}
	userRes := dto.UserResponse{}

	newUser.ID = user.ID
	newUser.Name = user.Name
	newUser.Email = user.Email
	newUser.Phone = user.Phone
	newUser.CurrentJob = user.CurrentJob
	newUser.Age = user.Age
	newUser.IsJobPoster = user.IsJobPoster

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return dto.UserResponse{}, shared.ErrCreateUsers
	}
	newUser.Password = string(hash)

	uc, err := uu.userRepo.Create(ctx, newUser)
	if err != nil {
		return dto.UserResponse{}, shared.ErrCreateUsers
	}

	userRes.ID = uc.ID
	userRes.Name = uc.Name
	userRes.Email = uc.Email
	userRes.Phone = uc.Phone
	userRes.CurrentJob = uc.CurrentJob
	userRes.Age = uc.Age
	userRes.IsJobPoster = uc.IsJobPoster

	return userRes, nil
}

func (uu *userUsecase) LoginUser(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	output := dto.LoginResponse{}

	user, err := uu.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user.ID == 0 {
		if errors.Is(err, shared.ErrRecordNotFound) {
			return output, shared.ErrUserDoesntExist
		}
		return output, shared.ErrFailedLogin
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return output, shared.ErrInvalidPassword
	}

	claims := helper.JWTClaims{
		UserId: user.ID,
	}

	userPayload := dto.UserPayload{
		ID:       user.ID,
		Name:     user.Email,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
	}

	token, err := helper.AuthorizedJWT(claims, userPayload)
	if err != nil {
		return output, shared.ErrFailedLogin
	}

	output = dto.LoginResponse{
		AccessToken: token,
	}

	return output, nil
}
