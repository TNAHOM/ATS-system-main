package jobpost

import (
	"context"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	models "github.com/TNAHOM/ATS-system-main/internal/constants/model"
	"github.com/TNAHOM/ATS-system-main/internal/storage"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JobPost struct {
	log *zap.Logger
	db  *gorm.DB
}

func Init(log *zap.Logger, db *gorm.DB) storage.JobPosts {
	return &JobPost{log, db}
}

func (j *JobPost) CreateJobPost(ctx context.Context, jobPost dto.CreateJobPostRequest) (dto.CreateJobPostResponse, error) {
	params := models.JobPost{
		ID:               jobPost.ID,
		Title:            jobPost.Title,
		Description:      jobPost.Description,
		Responsibilities: jobPost.Responsibilities,
		Requirements:     jobPost.Requirements,

		DescriptionEmbedding:      jobPost.DescriptionEmbedding,
		RequirementsEmbedding:     jobPost.RequirementsEmbedding,
		ResponsibilitiesEmbedding: jobPost.ResponsibilitiesEmbedding,

		UserID:   jobPost.UserID,
		Deadline: jobPost.Deadline,
	}
	if err := j.db.Create(&params).Error; err != nil {
		j.log.Error("Failed to create job post", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}
	return dto.CreateJobPostResponse{
		ID:               jobPost.ID,
		Title:            jobPost.Title,
		Description:      jobPost.Description,
		Responsibilities: jobPost.Responsibilities,
		Requirements:     jobPost.Requirements,
		UserID:           jobPost.UserID,
		Deadline:         jobPost.Deadline,
	}, nil
}

func (j *JobPost) GetAllJobPosts(ctx context.Context) ([]dto.GetAllJobPostsResponse, error) {
	var jobPosts []models.JobPost
	if err := j.db.WithContext(ctx).Find(&jobPosts).Error; err != nil {
		j.log.Error("failed to fetch job posts", zap.Error(err))
		return nil, err
	}

	res := make([]dto.GetAllJobPostsResponse, len(jobPosts))
	for i, jp := range jobPosts {
		res[i] = dto.GetAllJobPostsResponse{
			ID:               jp.ID,
			Title:            jp.Title,
			Description:      jp.Description,
			Responsibilities: jp.Responsibilities,
			Requirements:     jp.Requirements,
			UserID:           jp.UserID,
			Deadline:         jp.Deadline,
		}
	}
	return res, nil
}
