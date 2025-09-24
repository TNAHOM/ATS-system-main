package storage

import (
	"context"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
)

type Users interface {
	CreateUser(ctx context.Context, user dto.CreateUserRequest) (dto.CreateUserResponse, error)
	UserExist(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, req dto.LoginUserRequest) (dto.LoginUserResponse, error)
	GetAllUsers(ctx context.Context) (users []dto.GetAllUsers, error error)

	UpdateToken(ctx context.Context, updateFieldToken dto.UpdateTokenResponse) (bool, error)
}

type JobPosts interface {
	CreateJobPost(ctx context.Context, jobPost dto.CreateJobPostRequest) (dto.CreateJobPostResponse, error)
	GetAllJobPosts(ctx context.Context) ([]dto.GetAllJobPostsResponse, error)
	GetJobPostByID(ctx context.Context, id string) (dto.UpdateJobPostResponse, error)
	UpdateJobPost(ctx context.Context, req dto.UpdateJobPostRequest) (dto.UpdateJobPostResponse, error)
	DeleteJobPost(ctx context.Context, id string) error
}
