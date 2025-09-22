package user

import (
	"context"
	"errors"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/internal/module"
	"github.com/TNAHOM/ATS-system-main/internal/storage"
	"github.com/TNAHOM/ATS-system-main/platform"
	"github.com/TNAHOM/ATS-system-main/platform/encryption"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type user struct {
	logger      *zap.Logger
	userStorage storage.Users
	encryption  platform.Encryption
}

func Init(logger *zap.Logger, userStorage storage.Users, encryption platform.Encryption) module.User {
	return &user{logger: logger, userStorage: userStorage, encryption: encryption}
}

func (u *user) CreateUser(ctx context.Context, user dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	exists, err := u.userStorage.UserExist(ctx, user.Email)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}
	if exists {
		return dto.CreateUserResponse{}, errors.New("this user already exists with this email")
	}

	params := dto.GenerateUpdateToken{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UserType:  user.UserType,
	}

	signedToken, signedRefreshToken, err := u.encryption.GenerateToken(params)

	if err != nil {
		return dto.CreateUserResponse{}, err
	}
	hashedPassword, err := encryption.HashPassword(user.Password)
	if err != nil {
		u.logger.Error(err.Error(), zap.Any("request", user))
		return dto.CreateUserResponse{}, err
	}

	user.ID = uuid.New().String()
	user.Password = hashedPassword
	user.Token = signedToken
	user.RefreshToken = signedRefreshToken

	res, err := u.userStorage.CreateUser(ctx, user)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	return res, nil
}

func (u *user) LoginUser(ctx context.Context, loginReq dto.LoginUserRequest) (dto.LoginUserResponse, error) {
	res, err := u.userStorage.GetUserByEmail(ctx, loginReq)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	verified, err := encryption.VerifyPassword(res.Password, loginReq.Password)
	if err != nil {
		u.logger.Error(err.Error(), zap.Any("request", loginReq))
		return dto.LoginUserResponse{}, errors.New("internal server error")
	}
	if !verified {
		u.logger.Info("Password not correct", zap.Any("request", loginReq))
		return dto.LoginUserResponse{}, errors.New("password not correct")
	}
	params := dto.GenerateUpdateToken{
		ID:        res.ID,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		UserType:  res.UserType,
	}

	signedToken, signedRefreshToken, err := u.encryption.GenerateToken(params)
	if err != nil {
		u.logger.Error(err.Error(), zap.Any("request", res))
		return dto.LoginUserResponse{}, errors.New("internal server error")
	}

	res.Token = signedToken
	res.RefreshToken = signedRefreshToken

	updated, err := u.userStorage.UpdateToken(ctx, dto.UpdateTokenResponse{
		ID:           res.ID,
		Token:        res.Token,
		RefreshToken: res.RefreshToken,
	})
	if err != nil {
		u.logger.Error(err.Error(), zap.Any("request", res))
		return dto.LoginUserResponse{}, errors.New("internal server error")
	}
	if !updated {
		u.logger.Error("Failed to update the token", zap.Any("request", res))
		return dto.LoginUserResponse{}, errors.New("internal server error")
	}

	return res, nil
}

func (u *user) GetAllUsers(ctx context.Context) ([]dto.GetAllUsers, error) {
	res, err := u.userStorage.GetAllUsers(ctx)
	if err != nil {
		return []dto.GetAllUsers{}, err
	}

	return res, nil
}
