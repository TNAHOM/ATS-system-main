package module

import (
	"context"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
)

type User interface {
	CreateUser(ctx context.Context, user dto.CreateUserRequest) (dto.CreateUserResponse, error)
	LoginUser(ctx context.Context, loginReq dto.LoginUserRequest) (dto.LoginUserResponse, error)

	GetAllUsers(ctx context.Context) ([]dto.GetAllUsers, error)
}

type JobPost interface {
	CreateJobPost(ctx context.Context, jobPost dto.CreateJobPostRequest) (dto.CreateJobPostResponse, error)
}
