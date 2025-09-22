package user

import (
	"context"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	models "github.com/TNAHOM/ATS-system-main/internal/constants/model"
	"github.com/TNAHOM/ATS-system-main/internal/storage"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

	type user struct {
		log *zap.Logger
		db  *gorm.DB
	}

	func Init(log *zap.Logger, db *gorm.DB) storage.Users {
		return &user{log, db}
	}

	func (u *user) CreateUser(ctx context.Context, user dto.CreateUserRequest) (dto.CreateUserResponse, error) {
		params := models.User{
			ID:           user.ID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Phone:        user.Phone,
			Password:     user.Password,
			Email:        user.Email,
			UserType:     user.UserType,
			Token:        user.Token,
			RefreshToken: user.RefreshToken,
		}

		response := u.db.Create(&params)
		if response.Error != nil {
			u.log.Error(response.Error.Error(), zap.Any("request", user))
			return dto.CreateUserResponse{}, response.Error
		}
		return dto.CreateUserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			// Password:     user.Password,
			Email:        user.Email,
			UserType:     user.UserType,
			Token:        user.Token,
			RefreshToken: user.RefreshToken,
		}, nil
	}

func (u *user) UserExist(ctx context.Context, email string) (bool, error) {
	var count int64
	err := u.db.Model(&models.User{}).WithContext(ctx).Where("email = ?", email).Count(&count).Error
	if err != nil {
		u.log.Error(err.Error(), zap.Any("request", email))
		return false, err
	}

	return count > 0, nil

}

func (u *user) UpdateToken(ctx context.Context, updateTokenField dto.UpdateTokenResponse) (bool, error) {

	tx := u.db.Model(&models.User{}).Where("id = ?", updateTokenField.ID).Omit("ID").Updates(updateTokenField)
	if tx.Error != nil {
		return false, tx.Error
	}

	return true, nil
}

func (u *user) GetUserByEmail(ctx context.Context, req dto.LoginUserRequest) (dto.LoginUserResponse, error) {
	var userModel models.User
	tx := u.db.WithContext(ctx).Where("email = ?", req.Email).First(&userModel)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return dto.LoginUserResponse{}, tx.Error
		}
		return dto.LoginUserResponse{}, tx.Error
	}

	return dto.LoginUserResponse{
		ID:           userModel.ID,
		FirstName:    userModel.FirstName,
		LastName:     userModel.LastName,
		Email:        userModel.Email,
		UserType:     userModel.UserType,
		Password:     userModel.Password,
		Token:        userModel.Token,
		RefreshToken: userModel.RefreshToken,
	}, nil
}

func (u *user) GetAllUsers(ctx context.Context) (users []dto.GetAllUsers, error error) {
	var allUsers []models.User
	tx := u.db.WithContext(ctx).Find(&allUsers)
	if tx.Error != nil {
		return []dto.GetAllUsers{}, tx.Error
	}

	// Map models.User to dto.GetAllUsers
	users = make([]dto.GetAllUsers, len(allUsers))
	for i, user := range allUsers {
		users[i] = dto.GetAllUsers{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
			UserType:  user.UserType,
		}
	}
	return users, nil

}
