package jobpost

import (
	"context"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/internal/module"
	"github.com/TNAHOM/ATS-system-main/internal/storage"
	"github.com/TNAHOM/ATS-system-main/platform/ai"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JobPost struct {
	log            *zap.Logger
	jobPostStorage storage.JobPosts
}

func Init(logger *zap.Logger, jobPostStorage storage.JobPosts) module.JobPost {
	return &JobPost{log: logger, jobPostStorage: jobPostStorage}
}

func (j *JobPost) CreateJobPost(ctx context.Context, jobPost dto.CreateJobPostRequest) (dto.CreateJobPostResponse, error) {
	client, err := ai.NewClient(ctx)
	if err != nil {
		j.log.Error("Failed to create AI client", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}

	descriptionEmbedding, err := ai.Embedding(ctx, client, jobPost.Description)
	if err != nil {
		j.log.Error("Failed to generate description embedding", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}
	responsibilityEmbedding, err := ai.Embedding(ctx, client, jobPost.Responsibilities[])
	if err != nil {
		j.log.Error("Failed to generate responsibility embedding", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}
	requirementEmbedding, err := ai.Embedding(ctx, client, jobPost.Requirements[])
	if err != nil {
		j.log.Error("Failed to generate requirement embedding", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}

	jobPost.DescriptionEmbedding = descriptionEmbedding
	jobPost.RequirementsEmbedding = requirementEmbedding
	jobPost.ResponsibilitiesEmbedding= responsibilityEmbedding

	jobPost.ID = uuid.New().String()

	// Create the job post
	res, err := j.jobPostStorage.CreateJobPost(ctx, jobPost)
	if err != nil {
		j.log.Error("Failed to create job post", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}

	return dto.CreateJobPostResponse{
		ID:          res.ID,
		Description: res.Description,
		UserID:      res.UserID,
		Deadline:    res.Deadline,
	}, nil
}

// func GetAllJobPosts() {}
