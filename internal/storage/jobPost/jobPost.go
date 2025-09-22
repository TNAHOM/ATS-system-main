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
		ID:          jobPost.ID,
		Description: jobPost.Description,
		UserID:      jobPost.UserID,
		Deadline:    jobPost.Deadline,
	}, nil
}
