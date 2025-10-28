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
	GetAllJobPostsByUserId(ctx context.Context) ([]dto.GetJobPostsResponse, error)
	GetAllJobPosts(ctx context.Context) ([]dto.GetJobPostsResponse, error)

	GetJobPostByID(ctx context.Context, id string) (dto.GetJobPostsResponse, error)
	UpdateJobPost(ctx context.Context, req dto.UpdateJobPostRequest) (dto.GetJobPostsResponse, error)
	DeleteJobPost(ctx context.Context, id string) error
}

type Application interface {
	// GetApplicationsByApplicantID(ctx context.Context, applicantID string) ([]dto.GetApplicationsResponse, error)
	GetApplicationsByJobPostID(ctx context.Context, jobPostID string) ([]dto.GetApplicationsResponse, error)
}
