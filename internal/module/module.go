package module

import (
	"context"

	enum "github.com/TNAHOM/ATS-system-main/internal/constants/Enum"
	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
)

type User interface {
	CreateUser(ctx context.Context, user dto.CreateUserRequest) (dto.CreateUserResponse, error)
	LoginUser(ctx context.Context, loginReq dto.LoginUserRequest) (dto.LoginUserResponse, error)

	GetAllUsers(ctx context.Context) ([]dto.GetAllUsers, error)
}

type JobPost interface {
	CreateJobPost(ctx context.Context, jobPost dto.CreateJobPostRequest) (dto.CreateJobPostResponse, error)
	GetJobPostByID(ctx context.Context, id string) (dto.GetJobPostsResponse, error)
	GetAllJobPostsByUserId(ctx context.Context) ([]dto.GetJobPostsResponse, error)
	GetAllJobPosts(ctx context.Context) ([]dto.GetJobPostsResponse, error)
	UpdateJobPost(ctx context.Context, req dto.UpdateJobPostRequest) (dto.GetJobPostsResponse, error)
	DeleteJobPost(ctx context.Context, id string) error
}

type Application interface {
	// GetApplicationsByApplicantID(ctx context.Context, applicantID string) ([]dto.GetApplicationsResponse, error)
	GetApplicationsByJobPostID(ctx context.Context, jobPostID string, progressStatus enum.ProgressStatus) ([]dto.GetApplicationsResponse, error)
}
