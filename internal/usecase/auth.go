package usecase

import (
	"api_shope/dto"
	"api_shope/internal/repository"
	"api_shope/utils/helper"
	"errors"
)

type AuthUsecase interface {
	Register(req *dto.RegisterReq) error
	Login(req *dto.LoginReq) (string, error)
}

type authUsecase struct {
	authRepo repository.AuthRepo
}

func NewAuthUsecase(authRepo repository.AuthRepo) AuthUsecase {
	return &authUsecase{authRepo}
}

func (u *authUsecase) Register(req *dto.RegisterReq) error {
	valid := helper.IsValidEmail(req.Email)
	if !valid {
		return helper.ErrInvalidEmail
	}
	hashsed, err := helper.HashPasswrd(req.Password)
	if err != nil {
		return err
	}

	req.Password = hashsed
	if err := u.authRepo.Register(req); err != nil {
		return err
	}
	return nil
}

func (u *authUsecase) Login(req *dto.LoginReq) (string, error) {
	valid := helper.IsValidEmail(req.Email)
	if !valid {
		return "", helper.ErrInvalidEmail
	}
	user, err := u.authRepo.LoginEmail(req.Email)
	if err != nil {
		return "", err
	}

	if valid := helper.ComparePassword(user.Password, req.Password); !valid {
		return "", errors.New("email dan password tidak cocok")
	}

	jwt, err := helper.GenerateJWT(user.Email, user.ID)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
